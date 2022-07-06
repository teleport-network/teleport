#!/bin/bash

docker stop peer5
sudo rm -rf /usr/bin/bitchaind
sudo cp -r ~/go/bin/bitchaind /usr/bin
sudo bitchaind rollback-any --home ~/bitchain_dev/validators/validator5/bitchain/ --height 20

sleep 20
docker start peer5
docker logs -f peer5 --tail=100
