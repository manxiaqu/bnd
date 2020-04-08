# bnd
blockchain network deployer-bash script to deploy ethereum poa/pow network

# rquirement

`docker`, `docker-compose` is needed to run `bnd`(only test it on ubuntu).

# run 

```bash

cd $bnd
go build .
./bnd #./bnd -c config.yaml

```

config example

```yaml

# data location
base:
  basedir: build

genesis:
  # support ethash, clique
  consensus: ethash
  common:
    chainid: "1000"
    alloc:
    - "0x00000000000000000000000000000000000003e8"
    - "0x00000000000000000000000000000000000003e9"
    - "0x00000000000000000000000000000000000003ea"
  # signer will be generated automatically
  clique：
    clique:
      period: 15
      epoch: 3000

miner:
  amount: 1
  # will be replaced by signer address if consensus is clique
  ethbases:
  - "0x0000000000000000000000000000000000000001"

rpc:
  amount: 1
```


