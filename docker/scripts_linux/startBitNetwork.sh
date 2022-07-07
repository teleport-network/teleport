docker rm -f $(docker ps -aq)

# remove existed network
docker network rm  bitnetwork-br7
# create new network
docker network create --subnet=172.172.0.0/24 bitnetwork-br7

# remove pre data
sudo rm -rf ~/bitnetwork_dev
sudo cp -r ../bitnetwork_dev ~/bitnetwork_dev

# start validators
docker run -itd  --net bitnetwork-br7  --ip 172.172.0.2 --name=peer0  -p 26656:26656 -p 26657:26657 -v ~/bitnetwork_dev/validators/validator0/bitnetwork/:/root/bitnetwork -v ~/go/bin/bitnetworkd:/usr/bin/bitnetwork ubuntu:20.04  bitnetwork start --home /root/bitnetwork --log_level info
docker run -itd  --net bitnetwork-br7  --ip 172.172.0.3 --name=peer1  -p 26666:26656 -p 26667:26657 -v ~/bitnetwork_dev/validators/validator1/bitnetwork/:/root/bitnetwork  -v ~/go/bin/bitnetworkd:/usr/bin/bitnetwork ubuntu:20.04  bitnetwork start --home /root/bitnetwork --log_level info
docker run -itd  --net bitnetwork-br7  --ip 172.172.0.4 --name=peer2  -p 26676:26656 -p 26677:26657 -v ~/bitnetwork_dev/validators/validator2/bitnetwork/:/root/bitnetwork  -v ~/go/bin/bitnetworkd:/usr/bin/bitnetwork ubuntu:20.04  bitnetwork start --home /root/bitnetwork --log_level info
docker run -itd  --net bitnetwork-br7  --ip 172.172.0.5 --name=peer3  -p 26686:26656 -p 26687:26657 -v ~/bitnetwork_dev/validators/validator3/bitnetwork/:/root/bitnetwork  -v ~/go/bin/bitnetworkd:/usr/bin/bitnetwork ubuntu:20.04  bitnetwork start --home /root/bitnetwork --log_level info
docker run -itd  --net bitnetwork-br7  --ip 172.172.0.6 --name=peer4  -p 26696:26656 -p 26697:26657 -v ~/bitnetwork_dev/validators/validator4/bitnetwork/:/root/bitnetwork  -v ~/go/bin/bitnetworkd:/usr/bin/bitnetwork ubuntu:20.04  bitnetwork start --home /root/bitnetwork --log_level info
docker run -itd  --net bitnetwork-br7  --ip 172.172.0.7 --name=peer5  -p 26706:26656 -p 26707:26657 -v ~/bitnetwork_dev/validators/validator5/bitnetwork/:/root/bitnetwork  -v ~/go/bin/bitnetworkd:/usr/bin/bitnetwork ubuntu:20.04  bitnetwork start --home /root/bitnetwork --log_level info
docker run -itd  --net bitnetwork-br7  --ip 172.172.0.8 --name=peer6  -p 26716:26656 -p 26717:26657 -v ~/bitnetwork_dev/validators/validator6/bitnetwork/:/root/bitnetwork  -v ~/go/bin/bitnetworkd:/usr/bin/bitnetwork ubuntu:20.04  bitnetwork start --home /root/bitnetwork --log_level info


docker logs -f peer0 --tail=500
