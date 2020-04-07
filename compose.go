package main

import (
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
)

// GetDMiner returns default config for docker-compose to start a ethereum miner.
func GetDMiner(name string, i int, ethbase string, miners, bootnodes []string) *config.ServiceConfig {
	service := &config.ServiceConfig{}
	service.ContainerName = name
	service.Image = ETHImage

	volumes := yaml.Volumes{}
	volumes.Volumes = make([]*yaml.Volume, 0)
	volumes.Volumes = append(volumes.Volumes, &yaml.Volume{"./" + name, "/root", ""})
	service.Volumes = &volumes

	service.Ports = []string{strconv.Itoa(lisPort+i) + ":30303"}

	// entrypoint
	command := yaml.Command{}
	command = append(command, "geth")
	command = append(command, "--nodekey")
	command = append(command, "/root/boot.key")
	command = append(command, "--miner.etherbase="+ethbase)
	command = append(command, "--mine")
	command = append(command, "--miner.threads=2")
	if i > 0 && len(bootnodes) > 0 {
		command = append(command, "--bootnodes")
		command = append(command, Contact(bootnodes))
	}

	if i > 0 && len(miners) > 0 {
		service.DependsOn = miners
	}

	service.Entrypoint = command

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
