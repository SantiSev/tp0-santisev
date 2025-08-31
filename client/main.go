package main

import (
	"github.com/op/go-logging"

	"github.com/7574-sistemas-distribuidos/docker-compose-init/client/common"
	"github.com/7574-sistemas-distribuidos/docker-compose-init/client/common/config"
)

var log = logging.MustGetLogger("log")

func main() {
	clientConfig, err := config.InitConfig()
	if err != nil {
		log.Criticalf("%s", err)
	}

	if err := config.InitLogger(clientConfig.LogLevel); err != nil {
		log.Criticalf("%s", err)
	}

	// Print program config with debugging purposes
	clientConfig.PrintConfig()

	client := common.NewClient(*clientConfig)
	client.StartClientLoop()
}
