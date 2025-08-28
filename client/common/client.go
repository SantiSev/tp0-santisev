package common

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/7574-sistemas-distribuidos/docker-compose-init/client/common/network"
	"github.com/7574-sistemas-distribuidos/docker-compose-init/client/common/protocol"
	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("log")
var err error

type Client struct {
	config      ClientConfig
	connManager network.ConnectionManager
	connSocket  *network.ConnectionInterface
	betHandler  protocol.BetHandler
}

func NewClient(config ClientConfig) *Client {
	client := &Client{
		config:      config,
		connManager: *network.NewConnectionManager(),
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

	c.connSocket, err = c.connManager.Connect(c.config.ServerAddress, c.config.ID)

	if err != nil {
		return
	}

	for msgID := 1; msgID <= c.config.LoopAmount; msgID++ {

		select {
		case <-done:
			log.Infof("action: exit | result: success | client_id: %v", c.config.ID)
			return
		default:
		}

		if err != nil {
			log.Errorf("action: create_bet | result: fail | client_id: %v | error: %v", c.config.ID, err)
			c.Shutdown()
			return
		}

		err = c.betHandler.SendBet(*c.config.Bet, c.connSocket)

		if err != nil {
			log.Errorf("action: send_message | result: fail | client_id: %v | error: %v",
				c.config.ID,
				err,
			)
			c.Shutdown()
			return
		}

		err = c.betHandler.RecvBetConfirmation(c.connSocket)

		if err != nil {
			log.Errorf("action: recv_message | result: fail | client_id: %v | error: %v",
				c.config.ID,
				err,
			)
			c.Shutdown()
			return
		}

		time.Sleep(c.config.LoopPeriod)

	}
	log.Infof("action: transmission finished | result: success | client_id: %v", c.config.ID)
	c.betHandler.SendDone(c.connSocket)
	c.Shutdown()
}

func (c *Client) Shutdown() {
	c.connSocket.Close()

}
