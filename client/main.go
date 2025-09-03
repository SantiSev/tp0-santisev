package main

import (
	"os"

	"github.com/op/go-logging"

	"github.com/7574-sistemas-distribuidos/docker-compose-init/client/common/client"
	"github.com/7574-sistemas-distribuidos/docker-compose-init/client/common/config"
)

var log = logging.MustGetLogger("log")

func main() {
	clientConfig, err := config.InitConfig()
	if err != nil {
		log.Fatalf("%s", err)
		os.Exit(1)
	}

	if err := config.InitLogger(clientConfig.LogLevel); err != nil {
		log.Fatalf("%s", err)
		os.Exit(1)
	}

	clientConfig.PrintConfig()

	client := client.NewClient(*clientConfig)

	if client == nil {
		log.Fatalf("Failed to create client")
		os.Exit(1)
	}

	err = client.Run()
	if err != nil {
		log.Fatalf("%s", err)
		os.Exit(1)
	}
}
