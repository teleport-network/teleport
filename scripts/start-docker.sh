#!/bin/bash

KEY="mykey"
CHAINID="bitchain_9000-1"
MONIKER="mymoniker"
DATA_DIR=$(mktemp -d -t bitchain-datadir.XXXXX)

echo "create and add new keys"
./bitchain keys add $KEY --home $DATA_DIR --no-backup --chain-id $CHAINID --algo "eth_secp256k1" --keyring-backend test
echo "init Bitchain with moniker=$MONIKER and chain-id=$CHAINID"
./bitchain init $MONIKER --chain-id $CHAINID --home $DATA_DIR
echo "prepare genesis: Allocate genesis accounts"
./bitchain add-genesis-account \
"$(./bitchain keys show $KEY -a --home $DATA_DIR --keyring-backend test)" 1000000000000000000atele,1000000000000000000stake \
--home $DATA_DIR --keyring-backend test
echo "prepare genesis: Sign genesis transaction"
./bitchain gentx $KEY 1000000000000000000stake --keyring-backend test --home $DATA_DIR --keyring-backend test --chain-id $CHAINID
echo "prepare genesis: Collect genesis tx"
./bitchain collect-gentxs --home $DATA_DIR
echo "prepare genesis: Run validate-genesis to ensure everything worked and that the genesis file is setup correctly"
./bitchain validate-genesis --home $DATA_DIR

echo "starting bitchain node $i in background ..."
./bitchain start --pruning=nothing --rpc.unsafe \
--keyring-backend test --home $DATA_DIR \
>$DATA_DIR/node.log 2>&1 & disown

echo "started bitchain node"
tail -f /dev/null
