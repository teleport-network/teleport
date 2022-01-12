#!/bin/bash

KEY="mykey"
CHAINID="teleport_9000-1"
MONIKER="mymoniker"
DATA_DIR=$(mktemp -d -t teleport-datadir.XXXXX)

echo "create and add new keys"
./teleport keys add $KEY --home $DATA_DIR --no-backup --chain-id $CHAINID --algo "eth_secp256k1" --keyring-backend test
echo "init Teleport with moniker=$MONIKER and chain-id=$CHAINID"
./teleport init $MONIKER --chain-id $CHAINID --home $DATA_DIR
echo "prepare genesis: Allocate genesis accounts"
./teleport add-genesis-account \
"$(./teleport keys show $KEY -a --home $DATA_DIR --keyring-backend test)" 1000000000000000000atele,1000000000000000000stake \
--home $DATA_DIR --keyring-backend test
echo "prepare genesis: Sign genesis transaction"
./teleport gentx $KEY 1000000000000000000stake --keyring-backend test --home $DATA_DIR --keyring-backend test --chain-id $CHAINID
echo "prepare genesis: Collect genesis tx"
./teleport collect-gentxs --home $DATA_DIR
echo "prepare genesis: Run validate-genesis to ensure everything worked and that the genesis file is setup correctly"
./teleport validate-genesis --home $DATA_DIR

echo "starting teleport node $i in background ..."
./teleport start --pruning=nothing --rpc.unsafe \
--keyring-backend test --home $DATA_DIR \
>$DATA_DIR/node.log 2>&1 & disown

echo "started teleport node"
tail -f /dev/null