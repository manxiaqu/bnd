package main

import (
	"crypto/ecdsa"
	"io/ioutil"
	"math/big"
	"os"
	"path/filepath"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/params"
	log "github.com/inconshreveable/log15"
	yamlv2 "gopkg.in/yaml.v2"
)

var (
	chainID  = big.NewInt(1000)
	epoch    = uint64(3000)
	period   = uint64(15)
	build    = "build"
	ethash   = "ethash"
	clique   = "clique"
	ks       = "ks"
	dethbase = common.BigToAddress(big.NewInt(1))
)

// Config contains is the top level to deploy a ethash/clique network.
type Config struct {
	BaseConfig    *BaseConfig    `yaml:"base"`
	GenesisConfig *GenesisConfig `yaml:"genesis"`
	MinerConfig   *MinerConfig   `yaml:"miner"`
	RPCConfig     *RPCConfig     `yaml:"rpc,omitempty"`
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

// AssignDefault validates config and fill in default value if any error occurs.
func (c *Config) AssignDefault() {
	c.assignBase()
	// assign miner must before assign genesis if consensus is clique
	// because system will generate new key for every miner for mining purpose.
	c.assignMiner()
	c.assignGenesis()
	if c.RPCConfig == nil {
		c.RPCConfig = DRPCConfig()
	}
}

func (c *Config) assignMiner() {
	if c.MinerConfig == nil {
		c.MinerConfig = DMinerConfig()
		return
	}
	if c.MinerConfig.Amount < 0 {
		c.MinerConfig.Amount = 1
	}
	if len(c.MinerConfig.Ethbases) < c.MinerConfig.Amount {
		if c.MinerConfig.Ethbases == nil || len(c.MinerConfig.Ethbases) == 0 {
			c.MinerConfig.Ethbases = make([]common.Address, 0)
		}

		for i := 0; i < c.MinerConfig.Amount-len(c.MinerConfig.Ethbases); i++ {
			c.MinerConfig.Ethbases = append(c.MinerConfig.Ethbases, dethbase)
		}
	}
}

func (c *Config) assignGenesis() {
	if c.GenesisConfig == nil {
		c.GenesisConfig = DGenesisConfig()
		return
	}
	switch c.GenesisConfig.Consensus {
	case ethash:
	case clique:
		// always generate new key for signers.
		c.GenesisConfig.keys = make([]*ecdsa.PrivateKey, 0)
		c.GenesisConfig.Signers = make([]common.Address, 0)
		c.MinerConfig.Ethbases = make([]common.Address, 0)
		for i := 0; i < c.MinerConfig.Amount; i++ {
			key, _ := crypto.GenerateKey()
			c.GenesisConfig.keys = append(c.GenesisConfig.keys, key)
			c.GenesisConfig.Signers = append(c.GenesisConfig.Signers, crypto.PubkeyToAddress(key.PublicKey))

			// reset miner ethbase to signer address.
			c.MinerConfig.Ethbases = append(c.MinerConfig.Ethbases, crypto.PubkeyToAddress(key.PublicKey))
		}
		if c.GenesisConfig.Period == 0 {
			c.GenesisConfig.Period = period
		}
		if c.GenesisConfig.Epoch == 0 {
			c.GenesisConfig.Epoch = epoch
		}

	default:
		log.Warn("consensus not support, using ethash instead")
		c.GenesisConfig = DGenesisConfig()
	}

	if c.GenesisConfig.ChainID == nil || c.GenesisConfig.ChainID.Sign() != 1 {
		c.GenesisConfig.ChainID = chainID
	}
}

func (c *Config) assignBase() {
	if c.BaseConfig == nil {
		c.BaseConfig = DBaseConfig()
		return
	}

	// check dir is exist or not.
	if f, err := os.Stat(c.BaseConfig.BaseDir); err != nil {
		if os.IsNotExist(err) {
			log.Debug("base dir not exist, create it", "path", c.BaseConfig.BaseDir)
			os.MkdirAll(c.BaseConfig.BaseDir, 0755)
		} else {
			log.Debug("stat base dir failed, using defalut instead", "err", err)
			c.BaseConfig.BaseDir = getDBaseDir()
		}
	} else if !f.IsDir() {
		log.Debug("basedir isn't a dir, using default instead", "path", c.BaseConfig.BaseDir)
		c.BaseConfig.BaseDir = getDBaseDir()
	}

	c.BaseConfig.BaseDir, _ = filepath.Abs(c.BaseConfig.BaseDir)
}

// GetEthbase returns ethbase for miner i; returns default address if out of range.
func (c *Config) GetEthbase(i int) common.Address {
	if i > len(c.MinerConfig.Ethbases)+1 {
		return dethbase
	}

	return c.MinerConfig.Ethbases[i]
}

// BaseConfig contains normal config to work with deployer, like workdir.
type BaseConfig struct {
	BaseDir string `yaml:"basedir"`
}

// DBaseConfig returns default config to set base info.
func DBaseConfig() *BaseConfig {
	return &BaseConfig{
		BaseDir: getDBaseDir(),
	}
}

func getDBaseDir() string {
	dir, err := os.Getwd()
	if err != nil {
		log.Warn("get current directory failed", "err", err)
		// try using current relative path in linux
		dir = "."
	}

	return filepath.Join(dir, build)
}

// MinerConfig contains all config for miners to mine block on a eth/clique network.
type MinerConfig struct {
	// Amount == len(Ethbases)
	Amount int `yaml:"amount"`
	// will be covered in clique by default.
	Ethbases []common.Address `yaml:"ethbases"`
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
	Amount int `yaml:"amount"`
}

// DRPCConfig returns default config to set up rpc nodes.
func DRPCConfig() *RPCConfig {
	return &RPCConfig{
		Amount: 1,
	}
}

// GenesisConfig contains all config for a network.
type GenesisConfig struct {
	Consensus    string `yaml:"consensus"`
	CommonConfig `yaml:"common"`
	CliqueConfig `yaml:"clique,omitempty"`
}

// CliqueConfig contains clique network special config.
type CliqueConfig struct {
	// Signers assign automatically by deployer.
	// always generate new key for miner in clique.
	Signers             []common.Address `yaml:"-"`
	keys                []*ecdsa.PrivateKey
	params.CliqueConfig `yaml:"clique"`
}

// CommonConfig contains common config for genesis block.
type CommonConfig struct {
	ChainID *big.Int         `yaml:"chaindid"`
	Alloc   []common.Address `yaml:"alloc"`
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
		Consensus: ethash,
		CommonConfig: CommonConfig{
			ChainID: chainID,
			Alloc:   testAlloc,
		},
	}
}

func testEthashGenesisConfig() *GenesisConfig {
	return DGenesisConfig()
}

func testCliqueGenesisConfig() *GenesisConfig {
	return &GenesisConfig{
		Consensus: clique,
		CliqueConfig: CliqueConfig{
			Signers: testSigners,
			CliqueConfig: params.CliqueConfig{
				Period: period,
				Epoch:  epoch,
			},
		},
		CommonConfig: CommonConfig{
			ChainID: chainID,
			Alloc:   testAlloc,
		},
	}
}
