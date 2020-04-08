# bnd

一键部署以太坊

# 环境依赖

目前仅支持linux，运行本程序前需要安装docker, docker-compose

# 运行

```bash

cd $bnd
go build .
./bnd #./bnd -c config.yaml

```

配置示例

```yaml

# 数据存储路径
base:
  basedir: build

# 创世块配置
genesis:
  # 可选ethash, clique
  consensus: ethash
  common:
    chainid: "1000"
    alloc:
    - "0x00000000000000000000000000000000000003e8"
    - "0x00000000000000000000000000000000000003e9"
    - "0x00000000000000000000000000000000000003ea"
  # clique相关配置，目前signer都是自动生成得
  clique：
    clique:
      period: 15
      epoch: 3000

# 矿工配置
miner:
  amount: 1
  # 当共识模式为clique时，该部分自动会替换为clique中signer的地址
  ethbases:
  - "0x0000000000000000000000000000000000000001"

# rpc节点数量配置
rpc:
  amount: 1
```
