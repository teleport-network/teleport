#!/bin/bash

KEY="mykey"
CHAINID="bitnetwork_9000-1"
MONIKER="mymoniker"
DATA_DIR=$(mktemp -d -t bitnetwork-datadir.XXXXX)

echo "create and add new keys"
./bitnetwork keys add $KEY --home $DATA_DIR --no-backup --chain-id $CHAINID --algo "eth_secp256k1" --keyring-backend test
echo "init BitNetwork with moniker=$MONIKER and chain-id=$CHAINID"
./bitnetwork init $MONIKER --chain-id $CHAINID --home $DATA_DIR
echo "prepare genesis: Allocate genesis accounts"
./bitnetwork add-genesis-account \
"$(./bitnetwork keys show $KEY -a --home $DATA_DIR --keyring-backend test)" 1000000000000000000abit,1000000000000000000stake \
--home $DATA_DIR --keyring-backend test
echo "prepare genesis: Sign genesis transaction"
./bitnetwork gentx $KEY 1000000000000000000stake --keyring-backend test --home $DATA_DIR --keyring-backend test --chain-id $CHAINID
echo "prepare genesis: Collect genesis tx"
./bitnetwork collect-gentxs --home $DATA_DIR
echo "prepare genesis: Run validate-genesis to ensure everything worked and that the genesis file is setup correctly"
./bitnetwork validate-genesis --home $DATA_DIR

echo "starting bitnetwork node $i in background ..."
./bitnetwork start --pruning=nothing --rpc.unsafe \
--keyring-backend test --home $DATA_DIR \
>$DATA_DIR/node.log 2>&1 & disown

echo "started bitnetwork node"
tail -f /dev/null
