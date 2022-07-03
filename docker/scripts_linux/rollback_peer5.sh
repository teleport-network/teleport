#!/bin/bash

docker stop peer5
sudo rm -rf /usr/bin/bitchain
sudo cp -r ~/go/bin/bitchain /usr/bin
sudo bitchain rollback-any --home ~/bitchain_dev/validators/validator5/bitchain/ --height 50

sleep 20
docker start peer5
docker logs -f peer5 --tail=100
