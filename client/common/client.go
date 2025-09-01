package common

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/7574-sistemas-distribuidos/docker-compose-init/client/common/config"
	"github.com/7574-sistemas-distribuidos/docker-compose-init/client/common/network"
	"github.com/7574-sistemas-distribuidos/docker-compose-init/client/common/protocol"
	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("log")
var err error

// const BET_DATA_FILE = "data/agency_bets.csv"
const BET_DATA_FILE = "../.data/agency-1.csv"

type Client struct {
	config      config.ClientConfig
	connManager network.ConnectionManager
	connSocket  *network.ConnectionInterface
	betHandler  protocol.BetHandler
}

func NewClient(config config.ClientConfig) *Client {
	client := &Client{
		config:      config,
		connManager: *network.NewConnectionManager(),
		betHandler:  *protocol.NewBetHandler(config.MaxBatchAmount),
	}
	return client
}

func (c *Client) StartClientLoop() {

	sigChannel := make(chan os.Signal, 1)
	signal.Notify(sigChannel, syscall.SIGTERM)
	done := make(chan bool, 1)

	go func() {
		<-sigChannel
		c.Shutdown()
		done <- true
	}()

	select {
	case <-done:
		log.Infof("action: exit | result: success | client_id: %v", c.config.Id)
		return
	default:
	}

	c.connSocket, err = c.connManager.Connect(c.config.ServerAddress)

	if err != nil {
		log.Infof("action: connect | result: fail | client_id: %v", c.config.Id)
		return
	}

	err = c.betHandler.SendAllBetData(c.config.Id, BET_DATA_FILE, c.connSocket)

	if err != nil {
		log.Errorf("action: send_message | result: fail | client_id: %v | error: %v",
			c.config.Id,
			err,
		)
		c.Shutdown()
		return
	}

	log.Infof("action: transmission finished | result: success | client_id: %v", c.config.Id)
	c.Shutdown()
}

func (c *Client) Shutdown() {
	c.connSocket.Close()

}
