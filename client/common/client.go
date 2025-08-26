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

// ClientConfig Configuration used by the client
type ClientConfig struct {
	ID            string
	ServerAddress string
	LoopAmount    int
	LoopPeriod    time.Duration
}

type Client struct {
	config      ClientConfig
	connManager network.ConnectionManager
	connSocket  *network.ConnectionSocket
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

	// This is how i handle SIGTERM signals

	sigChannel := make(chan os.Signal, 1)
	signal.Notify(sigChannel, syscall.SIGTERM)
	done := make(chan bool, 1)

	go func() {
		<-sigChannel
		c.HandleShutdown()
		done <- true
	}()

	for msgID := 1; msgID <= c.config.LoopAmount; msgID++ {

		select {
		case <-done:
			log.Infof("action: exit | result: success | client_id: %v", c.config.ID)
			return
		default:
		}

		// Create the connection the server in every loop iteration. Send an
		c.connSocket = c.connManager.Connect(c.config.ServerAddress, c.config.ID)

		bet := protocol.NewBet(
			1,
			"santiago",
			"sev",
			42951041,
			"2000-10-08",
			42069,
		)

		err := c.betHandler.SendBet(*bet, c.connSocket)

		if err != nil {
			log.Errorf("action: send_message | result: fail | client_id: %v | error: %v",
				c.config.ID,
				err,
			)
			return
		}

		err = c.betHandler.RecvBetConfirmation(c.connSocket)

		if err != nil {
			log.Errorf("action: recv_message | result: fail | client_id: %v | error: %v",
				c.config.ID,
				err,
			)
			return
		}

		// Wait a time between sending one message and the next one
		time.Sleep(c.config.LoopPeriod)

	}
	log.Infof("action: loop_finished | result: success | client_id: %v", c.config.ID)
}

func (c *Client) HandleShutdown() {
	c.connSocket.Close()

}
