#!/bin/sh

set -o errexit -o nounset
sudo rm -rf ~/.bitchain/

KEYRING="test"
KEYALGO="secp256k1"
KEY="mykey"

CHAINID="bitchain_7001-1"
GENACCT="mykey"

# if $KEY exists it should be deleted
bitchaind keys add $KEY --keyring-backend $KEYRING --algo $KEYALGO

# Build genesis file incl account for passed address
coins="10000000000stake,100000000000samoleans"
bitchaind init --chain-id $CHAINID $CHAINID
bitchaind keys add validator --keyring-backend $KEYRING
bitchaind add-genesis-account $(bitchaind keys show validator -a --keyring-backend $KEYRING) $coins
bitchaind add-genesis-account $GENACCT $coins --keyring-backend $KEYRING
bitchaind gentx validator 5000000000stake --keyring-backend $KEYRING --chain-id $CHAINID
bitchaind collect-gentxs

# Set proper defaults and change ports
sed -i 's#"tcp://127.0.0.1:26657"#"tcp://0.0.0.0:26657"#g' ~/.bitchain/config/config.toml
sed -i 's/timeout_commit = "5s"/timeout_commit = "1s"/g' ~/.bitchain/config/config.toml
sed -i 's/timeout_propose = "3s"/timeout_propose = "1s"/g' ~/.bitchain/config/config.toml
sed -i 's/index_all_keys = false/index_all_keys = true/g' ~/.bitchain/config/config.toml

# Start the bitchaind
bitchaind start --pruning=nothing
