package main

// DONTCOVER

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
	"path/filepath"
	"strconv"

	"github.com/spf13/cobra"

	tmconfig "github.com/tendermint/tendermint/config"
	tmrand "github.com/tendermint/tendermint/libs/rand"
	tmtypes "github.com/tendermint/tendermint/types"
	tmtime "github.com/tendermint/tendermint/types/time"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdkserver "github.com/cosmos/cosmos-sdk/server"
	srvconfig "github.com/cosmos/cosmos-sdk/server/config"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	mintypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/tharsis/ethermint/crypto/hd"
	"github.com/tharsis/ethermint/server/config"
	srvflags "github.com/tharsis/ethermint/server/flags"
	ethermint "github.com/tharsis/ethermint/types"
	evmtypes "github.com/tharsis/ethermint/x/evm/types"

	"github.com/teleport-network/teleport/testutil/network"
	"github.com/teleport-network/teleport/types"
)

var (
	flagNodeDirPrefix  = "node-dir-prefix"
	flagNumValidators  = "v"
	flagOutputDir      = "output-dir"
	flagNodeDaemonHome = "node-daemon-home"
	flagTotalBalance   = "total-balance"
	flagIPAddrs        = "ip-addresses"
	flagEnableLogging  = "enable-logging"
	flagRPCAddress     = "rpc.address"
	flagAPIAddress     = "api.address"
	flagPrintMnemonic  = "print-mnemonic"
	mnemonicPath       = "/key_seed.json"
)

type initArgs struct {
	algo           string
	chainID        string
	keyringBackend string
	minGasPrices   string
	totalBalance   string
	nodeDaemonHome string
	nodeDirPrefix  string
	outputDir      string
	ipAddresses    []string
}

type startArgs struct {
	algo           string
	apiAddress     string
	chainID        string
	grpcAddress    string
	minGasPrices   string
	outputDir      string
	rpcAddress     string
	jsonrpcAddress string
	numValidators  int
	enableLogging  bool
	printMnemonic  bool
}

type KeySeed struct {
	Secret string `json:"secret"`
}

func addTestnetFlagsToCmd(cmd *cobra.Command) {
	cmd.Flags().Int(flagNumValidators, 4, "Number of validators to initialize the testnet with")
	cmd.Flags().StringP(flagOutputDir, "o", "./.testnets", "Directory to store initialization data for the testnet")
	cmd.Flags().StringP(flagTotalBalance, "", "55000", "total balance of chain")
	cmd.Flags().String(flags.FlagChainID, "", "genesis file chain-id, if left blank will be randomly created")
	cmd.Flags().String(sdkserver.FlagMinGasPrices, fmt.Sprintf("0.000006%s", types.AttoTele), "Minimum gas prices to accept for transactions; All fees in a tx must meet this minimum (e.g. 0.01tele,0.001stake)")
	cmd.Flags().String(flags.FlagKeyAlgorithm, string(hd.EthSecp256k1Type), "Key signing algorithm to generate keys for")
}

// NewTestnetCmd creates a root testnet command with subcommands to run an in-process testnet or initialize
// validator configuration files for running a multi-validator testnet in a separate process
func NewTestnetCmd(mbm module.BasicManager, genBalIterator banktypes.GenesisBalancesIterator) *cobra.Command {
	testnetCmd := &cobra.Command{
		Use:                        "testnet",
		Short:                      "subcommands for starting or configuring local testnets",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	testnetCmd.AddCommand(testnetStartCmd())
	testnetCmd.AddCommand(testnetInitFilesCmd(mbm, genBalIterator))

	return testnetCmd
}

// get cmd to initialize all files for tendermint testnet and application
func testnetInitFilesCmd(mbm module.BasicManager, genBalIterator banktypes.GenesisBalancesIterator) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init-files",
		Short: "Initialize config directories & files for a multi-validator testnet running locally via separate processes (e.g. Docker Compose or similar)",
		Long: `init-files will setup ip-addresses number of directories and populate each with
necessary files (private validator, genesis, config, etc.) for running ip-addresses number of validator nodes.

Booting up a network with these validator folders is intended to be used with Docker Compose,
or a similar setup where each node has a manually configurable IP address.

Note, strict routability for addresses is turned off in the config file.

Example:
	teleport testnet init-files --output-dir ./.testnets --ip-addresses 192.168.10.2
	`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			serverCtx := sdkserver.GetServerContextFromCmd(cmd)

			args := initArgs{}
			args.outputDir, _ = cmd.Flags().GetString(flagOutputDir)
			args.keyringBackend, _ = cmd.Flags().GetString(flags.FlagKeyringBackend)
			args.chainID, _ = cmd.Flags().GetString(flags.FlagChainID)
			args.minGasPrices, _ = cmd.Flags().GetString(sdkserver.FlagMinGasPrices)
			args.nodeDirPrefix, _ = cmd.Flags().GetString(flagNodeDirPrefix)
			args.nodeDaemonHome, _ = cmd.Flags().GetString(flagNodeDaemonHome)
			args.ipAddresses, _ = cmd.Flags().GetStringSlice(flagIPAddrs)
			args.totalBalance, _ = cmd.Flags().GetString(flagTotalBalance)
			args.algo, _ = cmd.Flags().GetString(flags.FlagKeyAlgorithm)
			if len(args.ipAddresses) == 0 {
				return errors.New("IP address list cannot be empty")
			}

			return initTestnetFiles(clientCtx, cmd, serverCtx.Config, mbm, genBalIterator, args)
		},
	}

	addTestnetFlagsToCmd(cmd)
	cmd.Flags().String(flagNodeDirPrefix, "node", "Prefix the directory name for each node with (node results in node0, node1, ...)")
	cmd.Flags().String(flagNodeDaemonHome, "teleport", "Home directory of the node's daemon configuration")
	cmd.Flags().StringSlice(flagIPAddrs, []string{"192.168.0.1"}, "List of IP addresses to use (i.e. `192.168.0.1,172.168.0.1` results in persistent peers list ID0@192.168.0.1:46656, ID1@172.168.0.1)")
	cmd.Flags().String(flags.FlagKeyringBackend, flags.DefaultKeyringBackend, "Select keyring's backend (os|file|test)")

	return cmd
}

// get cmd to start multi validator in-process testnet
func testnetStartCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start",
		Short: "Launch an in-process multi-validator testnet",
		Long: `testnet will launch an in-process multi-validator testnet,
and generate "v" directories, populated with necessary validator configuration files
(private validator, genesis, config, etc.).

Example:
	teleport testnet --v 4 --output-dir ./.testnets
	`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			args := startArgs{}
			args.outputDir, _ = cmd.Flags().GetString(flagOutputDir)
			args.chainID, _ = cmd.Flags().GetString(flags.FlagChainID)
			args.minGasPrices, _ = cmd.Flags().GetString(sdkserver.FlagMinGasPrices)
			args.numValidators, _ = cmd.Flags().GetInt(flagNumValidators)
			args.algo, _ = cmd.Flags().GetString(flags.FlagKeyAlgorithm)
			args.enableLogging, _ = cmd.Flags().GetBool(flagEnableLogging)
			args.rpcAddress, _ = cmd.Flags().GetString(flagRPCAddress)
			args.apiAddress, _ = cmd.Flags().GetString(flagAPIAddress)
			args.grpcAddress, _ = cmd.Flags().GetString(srvflags.GRPCAddress)
			args.jsonrpcAddress, _ = cmd.Flags().GetString(srvflags.JSONRPCAddress)
			args.printMnemonic, _ = cmd.Flags().GetBool(flagPrintMnemonic)

			return startTestnet(cmd, args)
		},
	}

	addTestnetFlagsToCmd(cmd)
	cmd.Flags().Bool(flagEnableLogging, false, "Enable INFO logging of tendermint validator nodes")
	cmd.Flags().String(flagRPCAddress, "tcp://0.0.0.0:26657", "the RPC address to listen on")
	cmd.Flags().String(flagAPIAddress, "tcp://0.0.0.0:1317", "the address to listen on for REST API")
	cmd.Flags().String(srvflags.GRPCAddress, config.DefaultGRPCAddress, "the gRPC server address to listen on")
	cmd.Flags().String(srvflags.JSONRPCAddress, config.DefaultJSONRPCAddress, "the JSON-RPC server address to listen on")
	cmd.Flags().Bool(flagPrintMnemonic, true, "print mnemonic of first validator to stdout for manual testing")
	return cmd
}

const nodeDirPerm = 0o755

// initTestnetFiles initializes testnet files for a testnet to be run in a separate process
func initTestnetFiles(
	clientCtx client.Context,
	cmd *cobra.Command,
	nodeConfig *tmconfig.Config,
	mbm module.BasicManager,
	genBalIterator banktypes.GenesisBalancesIterator,
	args initArgs,
) error {

	if args.chainID == "" {
		args.chainID = fmt.Sprintf("teleport_%d-1", tmrand.Int63n(9999999999999)+1)
	}
	numValidators := len(args.ipAddresses)
	nodeIDs := make([]string, numValidators)
	valPubKeys := make([]cryptotypes.PubKey, numValidators)

	appConfig := config.DefaultConfig()
	appConfig.MinGasPrices = args.minGasPrices
	appConfig.API.Enable = true
	appConfig.Telemetry.Enabled = true
	appConfig.Telemetry.PrometheusRetentionTime = 60
	appConfig.Telemetry.EnableHostnameLabel = false
	appConfig.Telemetry.GlobalLabels = [][]string{{"chain_id", args.chainID}}

	// If the total-balance is 100,000,000, the numvalidators is 11,
	// 0-9 validators' balance is 9,000,000, named lowBalance
	// the 10 validator's balance is 10,000,000, named highBlance
	totalBalance, _ := strconv.ParseInt(args.totalBalance, 10, 64)
	var lowBalance, highBalance int64
	if totalBalance/int64(numValidators) != 0 {
		avgBalance := strconv.FormatInt(totalBalance/int64(numValidators), 10)
		topDigit, _ := strconv.Atoi(avgBalance[0:1])
		lowBalance = int64(topDigit) * (math.Exp(big.NewInt(10), big.NewInt(int64(len(avgBalance)-1))).Int64())
		highBalance = totalBalance - lowBalance*int64(numValidators-1)
	} else {
		highBalance = totalBalance / int64(numValidators)
		lowBalance = highBalance
	}

	var (
		genAccounts []authtypes.GenesisAccount
		genBalances []banktypes.Balance
		genFiles    []string
	)

	inBuf := bufio.NewReader(cmd.InOrStdin())
	// generate private keys, node IDs, and initial transactions
	for i := 0; i < numValidators; i++ {
		nodeDirName := fmt.Sprintf("%s%d", args.nodeDirPrefix, i)
		nodeDir := filepath.Join(args.outputDir, nodeDirName, args.nodeDaemonHome)
		gentxsDir := filepath.Join(args.outputDir, "gentxs")

		nodeConfig.SetRoot(nodeDir)
		nodeConfig.RPC.ListenAddress = "tcp://0.0.0.0:26657"

		if err := os.MkdirAll(filepath.Join(nodeDir, "config"), nodeDirPerm); err != nil {
			_ = os.RemoveAll(args.outputDir)
			return err
		}

		nodeConfig.Moniker = nodeDirName

		ip := args.ipAddresses[i]
		var err error
		nodeIDs[i], valPubKeys[i], err = genutil.InitializeNodeValidatorFiles(nodeConfig)
		if err != nil {
			_ = os.RemoveAll(args.outputDir)
			return err
		}

		memo := fmt.Sprintf("%s@%s:26656", nodeIDs[i], ip)
		genFiles = append(genFiles, nodeConfig.GenesisFile())

		kb, err := keyring.New(sdk.KeyringServiceName(), args.keyringBackend, nodeDir, inBuf, hd.EthSecp256k1Option())
		if err != nil {
			return err
		}

		keyringAlgos, _ := kb.SupportedAlgorithms()
		algo, err := keyring.NewSigningAlgoFromString(args.algo, keyringAlgos)
		if err != nil {
			return err
		}

		var secret string
		var keyseed KeySeed
		var addr sdk.AccAddress

		//whether the pre-mnemonic existed
		_, err = os.Stat(nodeDir + mnemonicPath)
		if err != nil {
			addr, secret, err = sdkserver.GenerateSaveCoinKey(kb, nodeDirName, true, algo)
			if err != nil {
				_ = os.RemoveAll(args.outputDir)
				return err
			}
		} else {
			//remove the pre-keyring which may be still existed in the os
			_ = kb.Delete(nodeDirName)

			//restore the pre-mnemonic that had stored in key_seed.json
			mnemonicBytes, err := ioutil.ReadFile(nodeDir + mnemonicPath)
			if err != nil {
				return nil
			}

			err = json.Unmarshal(mnemonicBytes, &keyseed)
			if err != nil {
				return nil
			}

			info, err := kb.NewAccount(nodeDirName, keyseed.Secret, keyring.DefaultBIP39Passphrase, sdk.GetConfig().GetFullBIP44Path(), algo)
			if err != nil {
				return nil
			}
			addr = sdk.AccAddress(info.GetPubKey().Address())
			secret = keyseed.Secret
		}

		info := map[string]string{"secret": secret}

		cliPrint, err := json.Marshal(info)
		if err != nil {
			return err
		}

		// save private key seed words
		if err := network.WriteFile(fmt.Sprintf("%v.json", "key_seed"), nodeDir, cliPrint); err != nil {
			return err
		}

		var accStakingTokens sdk.Int
		if i != numValidators-1 {
			accStakingTokens = sdk.TokensFromConsensusPower(lowBalance, types.PowerReduction)
		} else {
			accStakingTokens = sdk.TokensFromConsensusPower(highBalance, types.PowerReduction)
		}
		coins := sdk.Coins{
			sdk.NewCoin(types.AttoTele, accStakingTokens),
		}

		genBalances = append(genBalances, banktypes.Balance{Address: addr.String(), Coins: coins.Sort()})
		genAccounts = append(genAccounts, &ethermint.EthAccount{
			BaseAccount: authtypes.NewBaseAccount(addr, nil, 0, 0),
			CodeHash:    common.BytesToHash(evmtypes.EmptyCodeHash).Hex(),
		})

		valTokens := sdk.TokensFromConsensusPower(100, types.PowerReduction)
		createValMsg, err := stakingtypes.NewMsgCreateValidator(
			sdk.ValAddress(addr),
			valPubKeys[i],
			sdk.NewCoin(types.AttoTele, valTokens),
			stakingtypes.NewDescription(nodeDirName, "", "", "", ""),
			stakingtypes.NewCommissionRates(sdk.OneDec(), sdk.OneDec(), sdk.OneDec()),
			sdk.OneInt(),
		)
		if err != nil {
			return err
		}

		txBuilder := clientCtx.TxConfig.NewTxBuilder()
		if err := txBuilder.SetMsgs(createValMsg); err != nil {
			return err
		}

		txBuilder.SetMemo(memo)

		txFactory := tx.Factory{}
		txFactory = txFactory.
			WithChainID(args.chainID).
			WithMemo(memo).
			WithKeybase(kb).
			WithTxConfig(clientCtx.TxConfig)

		if err := tx.Sign(txFactory, nodeDirName, txBuilder, true); err != nil {
			return err
		}

		txBz, err := clientCtx.TxConfig.TxJSONEncoder()(txBuilder.GetTx())
		if err != nil {
			return err
		}

		if err := network.WriteFile(fmt.Sprintf("%v.json", nodeDirName), gentxsDir, txBz); err != nil {
			return err
		}

		customAppTemplate, customAppConfig := config.AppConfig(types.AttoTele)
		srvconfig.SetConfigTemplate(customAppTemplate)
		if err := sdkserver.InterceptConfigsPreRunHandler(cmd, customAppTemplate, customAppConfig); err != nil {
			return err
		}

		srvconfig.WriteConfigFile(filepath.Join(nodeDir, "config/app.toml"), appConfig)
	}

	if err := initGenFiles(clientCtx, mbm, args.chainID, types.AttoTele, genAccounts, genBalances, genFiles, numValidators); err != nil {
		return err
	}

	err := collectGenFiles(
		clientCtx, nodeConfig, args.chainID, nodeIDs, valPubKeys, numValidators,
		args.outputDir, args.nodeDirPrefix, args.nodeDaemonHome, genBalIterator,
	)
	if err != nil {
		return err
	}

	cmd.PrintErrf("Successfully initialized %d node directories\n", numValidators)
	return nil
}

func initGenFiles(
	clientCtx client.Context,
	mbm module.BasicManager,
	chainID,
	coinDenom string,
	genAccounts []authtypes.GenesisAccount,
	genBalances []banktypes.Balance,
	genFiles []string,
	numValidators int,
) error {

	appGenState := mbm.DefaultGenesis(clientCtx.Codec)
	// set the accounts in the genesis state
	var authGenState authtypes.GenesisState
	clientCtx.Codec.MustUnmarshalJSON(appGenState[authtypes.ModuleName], &authGenState)

	accounts, err := authtypes.PackAccounts(genAccounts)
	if err != nil {
		return err
	}

	authGenState.Accounts = accounts
	appGenState[authtypes.ModuleName] = clientCtx.Codec.MustMarshalJSON(&authGenState)

	// set the balances in the genesis state
	var bankGenState banktypes.GenesisState
	clientCtx.Codec.MustUnmarshalJSON(appGenState[banktypes.ModuleName], &bankGenState)

	bankGenState.Balances = genBalances
	appGenState[banktypes.ModuleName] = clientCtx.Codec.MustMarshalJSON(&bankGenState)

	var stakingGenState stakingtypes.GenesisState
	clientCtx.Codec.MustUnmarshalJSON(appGenState[stakingtypes.ModuleName], &stakingGenState)

	stakingGenState.Params.BondDenom = coinDenom
	appGenState[stakingtypes.ModuleName] = clientCtx.Codec.MustMarshalJSON(&stakingGenState)

	var govGenState govtypes.GenesisState
	clientCtx.Codec.MustUnmarshalJSON(appGenState[govtypes.ModuleName], &govGenState)

	govGenState.DepositParams.MinDeposit[0].Denom = coinDenom
	appGenState[govtypes.ModuleName] = clientCtx.Codec.MustMarshalJSON(&govGenState)

	var mintGenState mintypes.GenesisState
	clientCtx.Codec.MustUnmarshalJSON(appGenState[mintypes.ModuleName], &mintGenState)

	mintGenState.Params.MintDenom = coinDenom
	appGenState[mintypes.ModuleName] = clientCtx.Codec.MustMarshalJSON(&mintGenState)

	var crisisGenState crisistypes.GenesisState
	clientCtx.Codec.MustUnmarshalJSON(appGenState[crisistypes.ModuleName], &crisisGenState)

	crisisGenState.ConstantFee.Denom = coinDenom
	appGenState[crisistypes.ModuleName] = clientCtx.Codec.MustMarshalJSON(&crisisGenState)

	var evmGenState evmtypes.GenesisState
	clientCtx.Codec.MustUnmarshalJSON(appGenState[evmtypes.ModuleName], &evmGenState)

	evmGenState.Params.EvmDenom = coinDenom
	appGenState[evmtypes.ModuleName] = clientCtx.Codec.MustMarshalJSON(&evmGenState)

	appGenStateJSON, err := json.MarshalIndent(appGenState, "", "  ")
	if err != nil {
		return err
	}

	genDoc := tmtypes.GenesisDoc{
		ChainID:    chainID,
		AppState:   appGenStateJSON,
		Validators: nil,
	}

	// generate empty genesis files for each validator and save
	for i := 0; i < numValidators; i++ {
		if err := genDoc.SaveAs(genFiles[i]); err != nil {
			return err
		}
	}
	return nil
}

func collectGenFiles(
	clientCtx client.Context, nodeConfig *tmconfig.Config, chainID string,
	nodeIDs []string, valPubKeys []cryptotypes.PubKey, numValidators int,
	outputDir, nodeDirPrefix, nodeDaemonHome string, genBalIterator banktypes.GenesisBalancesIterator,
) error {

	var appState json.RawMessage
	genTime := tmtime.Now()

	for i := 0; i < numValidators; i++ {
		nodeDirName := fmt.Sprintf("%s%d", nodeDirPrefix, i)
		nodeDir := filepath.Join(outputDir, nodeDirName, nodeDaemonHome)
		gentxsDir := filepath.Join(outputDir, "gentxs")
		nodeConfig.Moniker = nodeDirName

		nodeConfig.SetRoot(nodeDir)

		nodeID, valPubKey := nodeIDs[i], valPubKeys[i]
		initCfg := genutiltypes.NewInitConfig(chainID, gentxsDir, nodeID, valPubKey)

		genDoc, err := tmtypes.GenesisDocFromFile(nodeConfig.GenesisFile())
		if err != nil {
			return err
		}

		nodeAppState, err := genutil.GenAppStateFromConfig(clientCtx.Codec, clientCtx.TxConfig, nodeConfig, initCfg, *genDoc, genBalIterator)
		if err != nil {
			return err
		}

		if appState == nil {
			// set the canonical application state (they should not differ)
			appState = nodeAppState
		}

		genFile := nodeConfig.GenesisFile()

		// overwrite each validator's genesis file to have a canonical genesis time
		if err := genutil.ExportGenesisFileWithTime(genFile, chainID, nil, appState, genTime); err != nil {
			return err
		}
	}

	return nil
}

// startTestnet starts an in-process testnet
func startTestnet(cmd *cobra.Command, args startArgs) error {
	networkConfig := network.DefaultConfig()

	// Default networkConfig.ChainID is random, and we should only override it if chainID provided
	// is non-empty
	if args.chainID != "" {
		networkConfig.ChainID = args.chainID
	}
	networkConfig.SigningAlgo = args.algo
	networkConfig.MinGasPrices = args.minGasPrices
	networkConfig.NumValidators = args.numValidators
	networkConfig.EnableTMLogging = args.enableLogging
	networkConfig.RPCAddress = args.rpcAddress
	networkConfig.APIAddress = args.apiAddress
	networkConfig.GRPCAddress = args.grpcAddress
	networkConfig.JSONRPCAddress = args.jsonrpcAddress
	networkConfig.PrintMnemonic = args.printMnemonic
	networkLogger := network.NewCLILogger(cmd)

	baseDir := fmt.Sprintf("%s/%s", args.outputDir, networkConfig.ChainID)
	if _, err := os.Stat(baseDir); !os.IsNotExist(err) {
		return fmt.Errorf(
			"testnests directory already exists for chain-id '%s': %s, please remove or select a new --chain-id",
			networkConfig.ChainID, baseDir)
	}

	testnet, err := network.New(networkLogger, baseDir, networkConfig)
	if err != nil {
		return err
	}

	_, err = testnet.WaitForHeight(1)
	if err != nil {
		return err
	}

	cmd.Println("press the Enter Key to terminate")
	fmt.Scanln() // wait for Enter Key
	testnet.Cleanup()

	return nil
}
