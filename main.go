package main

import (
	"flag"

	log "github.com/inconshreveable/log15"
)

func main() {
	var (
		confPath = flag.String("c", "config.yaml", "config file path")
	)
	flag.Parse()

	conf, err := LoadConfig(*confPath)
	if err != nil {
		log.Warn("load config file failed, using default config instead", "err", err)
		conf = DConfig()
	}

	d := NewDeployer(conf)
	d.Deploy()
}
