package main

import (
	"os"
	"path/filepath"
	"strconv"

	"github.com/manxiaqu/libcompose/config"
	"github.com/manxiaqu/libcompose/yaml"
	yamlv2 "gopkg.in/yaml.v2"
)

var (
	// ETHImage is the offical ethereum docker image.
	ETHImage    = "ethereum/client-go"
	rpcPort     = 8545
	lisPort     = 30303
	miner       = "miner"
	rpc         = "rpc"
	bootkey     = "boot.key"
	dbase       = "build"
	composeFile = "compose.yaml"
	pw          = "123"
	pwfile      = "pwfile"
	signerkey   = "key"
)

// SetDMiner sets default config for docker-compose to start a ethereum miner.
func (dp *Deployer) SetDMiner(name string, i int) *config.ServiceConfig {
	service := &config.ServiceConfig{}
	service.ContainerName = name
	service.Image = ETHImage

	// mkdir for every miner.
	minerDir := filepath.Join(dp.config.BaseConfig.BaseDir, name)
	os.Mkdir(minerDir, 0777)
	// generate boot key for every miner.
	key, _ := GenSaveKey(filepath.Join(minerDir, bootkey))

	volumes := yaml.Volumes{}
	volumes.Volumes = make([]*yaml.Volume, 0)
	volumes.Volumes = append(volumes.Volumes, &yaml.Volume{minerDir, "/root", ""})
	service.Volumes = &volumes

	service.Ports = []string{strconv.Itoa(lisPort+i) + ":30303"}

	// entrypoint
	command := yaml.Command{}
	command = append(command, "geth")
	command = append(command, "--nodekey")
	command = append(command, "/root/boot.key")
	command = append(command, "--miner.etherbase="+dp.config.GetEthbase(i).String())
	command = append(command, "--mine")
	command = append(command, "--miner.threads=2")
	if i > 0 && len(dp.bootnodes) > 0 {
		command = append(command, "--bootnodes")
		command = append(command, Contact(dp.bootnodes))
	}
	// unlock the account for mining if signer not empty.
	if dp.config.GenesisConfig.Consensus == clique {
		// mkdir keystore.
		keystorepath := filepath.Join(minerDir, ks)
		os.MkdirAll(keystorepath, 0777)

		// save signer key to keystore and password file to miner dir.
		StoreKey(dp.config.GenesisConfig.keys[i], keystorepath, filepath.Join(keystorepath, signerkey))
		Flush([]byte(pw), filepath.Join(minerDir, pwfile))

		command = append(command, "--unlock="+dp.config.GenesisConfig.Signers[i].Hex())
		command = append(command, "--password")
		command = append(command, filepath.Join("/root", pwfile))
		command = append(command, "--keystore")
		command = append(command, filepath.Join("/root", ks))
	}

	if i > 0 && len(dp.miners) > 0 {
		service.DependsOn = dp.miners
	}

	service.Entrypoint = command

	// add miner info for later useage.
	dp.bootnodes = append(dp.bootnodes, GetBootNodeString(key, name, lisPort))
	dp.miners = append(dp.miners, name)

	return service
}

// GetDRPC .
func GetDRPC(name string, i int, miners, bootnodes []string) *config.ServiceConfig {
	service := &config.ServiceConfig{}
	service.ContainerName = name
	service.Image = ETHImage

	volumes := yaml.Volumes{}
	volumes.Volumes = make([]*yaml.Volume, 0)
	volumes.Volumes = append(volumes.Volumes, &yaml.Volume{"./" + name, "/root", ""})
	service.Volumes = &volumes

	service.Ports = []string{strconv.Itoa(rpcPort+i) + ":8545"}

	// entrypoint
	// --rpc --rpcaddr 0.0.0.0 --rpcapi eth,net
	command := yaml.Command{}
	command = append(command, "geth")
	command = append(command, "--rpc")
	command = append(command, "--rpcaddr")
	command = append(command, "0.0.0.0")
	command = append(command, "--rpcapi")
	command = append(command, "eth,net")
	if len(bootnodes) > 0 {
		command = append(command, "--bootnodes")
		command = append(command, Contact(bootnodes))
	}
	service.Entrypoint = command

	// add depends_on
	if len(miners) > 0 {
		service.DependsOn = miners
	}

	return service
}

func convertServiceToRawService(service *config.ServiceConfig) config.RawService {
	rawService := make(config.RawService)
	out, _ := yamlv2.Marshal(service)
	yamlv2.Unmarshal(out, &rawService)

	return rawService
}
