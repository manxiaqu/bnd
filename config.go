package main

import (
	"io/ioutil"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/params"
	yamlv2 "gopkg.in/yaml.v2"
)

// Config contains is the top level to deploy a ethash/clique network.
type Config struct {
	BaseConfig    *BaseConfig    `yaml:"baseConfig"`
	GenesisConfig *GenesisConfig `yaml:"genesisConfig"`
	MinerConfig   *MinerConfig   `yaml:"minerConfig"`
	RPCConfig     *RPCConfig     `yaml:"rpcConfig,omitempty"`
}

// WriteConfig flush config to disk.
func WriteConfig(path string, conf *Config) error {
	data, err := yamlv2.Marshal(conf)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(path, data, 0755)
}

// LoadConfig loads config by yaml file.
func LoadConfig(path string) (*Config, error) {
	raw, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var c Config
	if err := yamlv2.Unmarshal(raw, &c); err != nil {
		return nil, err
	}

	return &c, nil
}

// DConfig returns default config to set up a network.
func DConfig() *Config {
	return &Config{
		DBaseConfig(),
		DGenesisConfig(),
		DMinerConfig(),
		DRPCConfig(),
	}
}

// GetEthbase returns ethbase for miner i; returns default address if out of range.
func (c *Config) GetEthbase(i int) common.Address {
	if i > len(c.MinerConfig.Ethbases)+1 {
		return common.BigToAddress(big.NewInt(1))
	}

	return c.MinerConfig.Ethbases[i]
}

// BaseConfig contains normal config to work with deployer, like workdir.
type BaseConfig struct {
	BaseDir string
}

// DBaseConfig returns default config to set base info.
func DBaseConfig() *BaseConfig {
	return &BaseConfig{
		BaseDir: "./build",
	}
}

// MinerConfig contains all config for miners to mine block on a eth/clique network.
type MinerConfig struct {
	// Amount == len(Ethbases)
	Amount   int
	Ethbases []common.Address
}

// DMinerConfig returns default miner config to set up miner.
func DMinerConfig() *MinerConfig {
	return &MinerConfig{
		Amount:   1,
		Ethbases: []common.Address{common.BigToAddress(big.NewInt(1))},
	}
}

// RPCConfig contains all config to run servial rpc nodes.
type RPCConfig struct {
	Amount int
}

// DRPCConfig returns default config to set up rpc nodes.
func DRPCConfig() *RPCConfig {
	return &RPCConfig{
		Amount: 1,
	}
}

// GenesisConfig contains all config for a network.
type GenesisConfig struct {
	Consensus string
	CommonConfig
	CliqueConfig `yaml:"cliqueConfig,omitempty"`
}

// CliqueConfig contains clique network special config.
type CliqueConfig struct {
	Signers []common.Address
	params.CliqueConfig
}

// CommonConfig contains common config for genesis block.
type CommonConfig struct {
	ChainID *big.Int
	Alloc   []common.Address
}

var (
	testAlloc = []common.Address{
		common.BigToAddress(big.NewInt(1000)),
		common.BigToAddress(big.NewInt(1001)),
		common.BigToAddress(big.NewInt(1002)),
	}
	testSigners = []common.Address{
		common.BigToAddress(big.NewInt(2000)),
		common.BigToAddress(big.NewInt(2001)),
		common.BigToAddress(big.NewInt(2002)),
	}
)

// DGenesisConfig returns default ethash config to generate genesis block.
func DGenesisConfig() *GenesisConfig {
	return &GenesisConfig{
		Consensus: "ethash",
		CommonConfig: CommonConfig{
			ChainID: new(big.Int).SetUint64(1011),
			Alloc:   testAlloc,
		},
	}
}

func testEthashGenesisConfig() *GenesisConfig {
	return DGenesisConfig()
}

func testCliqueGenesisConfig() *GenesisConfig {
	return &GenesisConfig{
		Consensus: "clique",
		CliqueConfig: CliqueConfig{
			Signers: testSigners,
			CliqueConfig: params.CliqueConfig{
				Period: 15,
				Epoch:  3000,
			},
		},
		CommonConfig: CommonConfig{
			ChainID: new(big.Int).SetUint64(1011),
			Alloc:   testAlloc,
		},
	}
}
