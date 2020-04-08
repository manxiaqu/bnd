cd /home/ubuntu/git/bnd/build
docker run -d --name ethereum-node -v /home/ubuntu/git/bnd/build:/root ethereum/client-go init /root/genesis.json
sleep 1
docker stop ethereum-node
docker rm ethereum-node

sudo cp -r .ethereum  miner0 
 sudo cp -r .ethereum  rpc0 
 
sudo rm -fr .ethereum
docker-compose -f compose.yaml up -d
