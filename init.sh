#!/bin/bash

set -o errexit -o nounset
sudo rm -rf ~/.bitchain/

KEYRING="test"
KEYALGO="eth_secp256k1"
KEY="mykey"

CHAINID="bitchain_7001-1"
GENACCT="mykey"

# if $KEY exists it should be deleted
bitchaind keys add $KEY --keyring-backend $KEYRING --algo $KEYALGO

# Build genesis file incl account for passed address
coins="10000000000abit"
bitchaind init --chain-id $CHAINID $CHAINID
bitchaind keys add validator --keyring-backend $KEYRING
bitchaind add-genesis-account $(bitchaind keys show validator -a --keyring-backend $KEYRING) $coins
bitchaind add-genesis-account $GENACCT $coins --keyring-backend $KEYRING
bitchaind gentx validator 5000000000abit --keyring-backend $KEYRING --chain-id $CHAINID
bitchaind collect-gentxs


# Change parameter token denominations to atele
cat $HOME/.bitchain/config/genesis.json | jq '.app_state["staking"]["params"]["bond_denom"]="abit"' > $HOME/.bitchain/config/tmp_genesis.json && mv $HOME/.bitchain/config/tmp_genesis.json $HOME/.bitchain/config/genesis.json
cat $HOME/.bitchain/config/genesis.json | jq '.app_state["crisis"]["constant_fee"]["denom"]="abit"' > $HOME/.bitchain/config/tmp_genesis.json && mv $HOME/.bitchain/config/tmp_genesis.json $HOME/.bitchain/config/genesis.json
cat $HOME/.bitchain/config/genesis.json | jq '.app_state["gov"]["deposit_params"]["min_deposit"][0]["denom"]="abit"' > $HOME/.bitchain/config/tmp_genesis.json && mv $HOME/.bitchain/config/tmp_genesis.json $HOME/.bitchain/config/genesis.json
cat $HOME/.bitchain/config/genesis.json | jq '.app_state["mint"]["params"]["mint_denom"]="abit"' > $HOME/.bitchain/config/tmp_genesis.json && mv $HOME/.bitchain/config/tmp_genesis.json $HOME/.bitchain/config/genesis.json
cat $HOME/.bitchain/config/genesis.json | jq '.app_state["liquidity"]["params"]["pool_creation_fee"][0]["denom"]="abit"' > $HOME/.bitchain/config/tmp_genesis.json && mv $HOME/.bitchain/config/tmp_genesis.json $HOME/.bitchain/config/genesis.json


# Set proper defaults and change ports
if [[ "$OSTYPE" == "darwin"* ]]; then
      sed -i '' 's#"tcp://127.0.0.1:26657"#"tcp://0.0.0.0:26657"#g' ~/.bitchain/config/config.toml
      sed -i '' 's/timeout_commit = "5s"/timeout_commit = "1s"/g' ~/.bitchain/config/config.toml
      sed -i '' 's/timeout_propose = "3s"/timeout_propose = "1s"/g' ~/.bitchain/config/config.toml
      sed -i '' 's/index_all_keys = false/index_all_keys = true/g' ~/.bitchain/config/config.toml
  else
      sed -i 's#"tcp://127.0.0.1:26657"#"tcp://0.0.0.0:26657"#g' ~/.bitchain/config/config.toml
      sed -i 's/timeout_commit = "5s"/timeout_commit = "1s"/g' ~/.bitchain/config/config.toml
      sed -i 's/timeout_propose = "3s"/timeout_propose = "1s"/g' ~/.bitchain/config/config.toml
      sed -i 's/index_all_keys = false/index_all_keys = true/g' ~/.bitchain/config/config.toml
fi

# Start the bitchaind
bitchaind start --pruning=nothing
