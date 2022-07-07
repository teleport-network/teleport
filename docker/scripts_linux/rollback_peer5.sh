#!/bin/bash

docker stop peer5
sudo rm -rf /usr/bin/bitnetworkd
sudo cp -r ~/go/bin/bitnetworkd /usr/bin
sudo bitnetworkd rollback-any --home ~/bitnetwork_dev/validators/validator5/bitnetwork/ --height 20

sleep 20
docker start peer5
docker logs -f peer5 --tail=100
