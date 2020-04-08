package main

import (
	"html/template"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"

	log "github.com/inconshreveable/log15"
	"github.com/manxiaqu/libcompose/config"
	yamlv2 "gopkg.in/yaml.v2"
)

// Deployer represents ethash/clique network deployer.
type Deployer struct {
	config *Config

	rawServices   config.RawServiceMap
	composeConfig *config.Config

	// tmp info
	miners, rpcs []string
	bootnodes    []string
}

// NewDeployer creates a new instance of deployer.
func NewDeployer(conf *Config) *Deployer {
	return &Deployer{
		config:        conf,
		composeConfig: &config.Config{Version: "2.1"},
		rawServices:   make(config.RawServiceMap),
		miners:        make([]string, 0),
		rpcs:          make([]string, 0),
		bootnodes:     make([]string, 0),
	}
}

// Deploy deploys ethereum network base on config.
func (dp *Deployer) Deploy() {
	path := filepath.Join(dp.config.BaseConfig.BaseDir, genesisPath)
	MakeGenesis(dp.config.GenesisConfig, path)
	log.Info("generate genesis json success", "path", path)

	dp.setDMiners()
	dp.setDRPC()
	dp.flush()
}

// setDMiners sets servial configs for docker-compose to up servial miners.
func (dp *Deployer) setDMiners() {
	for i := 0; i < dp.config.MinerConfig.Amount; i++ {
		minerName := miner + strconv.Itoa(i)
		dp.rawServices[minerName] = convertServiceToRawService(dp.SetDMiner(minerName, i))
	}
}

// setDRPC sets config for docker-compose to start ethereum rpc node.
func (dp *Deployer) setDRPC() {
	for i := 0; i < dp.config.RPCConfig.Amount; i++ {
		rpcName := rpc + strconv.Itoa(i)
		// mkdir for every rpc.
		os.Mkdir(filepath.Join(dp.config.BaseConfig.BaseDir, rpcName), 0777)

		dp.rawServices[rpcName] = convertServiceToRawService(GetDRPC(rpcName, i, dp.miners, dp.bootnodes))
		dp.rpcs = append(dp.rpcs, rpcName)
	}
}

// flush writes shell/compose files to disk which is necessary for setting up ethereum network.
func (dp *Deployer) flush() {
	dp.composeConfig.Services = dp.rawServices
	out, err := yamlv2.Marshal(dp.composeConfig)
	if err != nil {
		panic(err)
	}

	Flush(out, filepath.Join(dp.config.BaseConfig.BaseDir, composeFile))
	tmpl, err := template.New("test").Parse(shellTemplate)
	if err != nil {
		panic(err)
	}

	data := Shell{}
	data.BaseDir = dp.config.BaseConfig.BaseDir
	data.Genesis = "/root/" + genesisPath
	data.Nodes = dp.miners
	data.Nodes = append(data.Nodes, dp.rpcs...)
	data.Compose = composeFile

	f, err := os.OpenFile("shell.sh", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		log.Error("create file failed", "err", err)
		return
	}
	if err := tmpl.Execute(f, data); err != nil {
		log.Error("generate shell by template failed", "err", err)
	}

	err = exec.Command("/bin/sh", "shell.sh").Run()
	if err != nil {
		log.Error("run shell.sh failed", "err", err)
	} else {
		log.Info("network success started using docker-compose")
	}
}
