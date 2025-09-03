package client

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/7574-sistemas-distribuidos/docker-compose-init/client/common/business"
	"github.com/7574-sistemas-distribuidos/docker-compose-init/client/common/network"
	"github.com/7574-sistemas-distribuidos/docker-compose-init/client/common/protocol"
	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("log")
var err error

type Client struct {
	config        ClientConfig
	connManager   network.ConnectionManager
	connInterface *network.ConnectionInterface
	betHandler    protocol.BetHandler
	agencyService business.AgencyService
}

func NewClient(config ClientConfig) *Client {
	agencyService, err := business.NewAgencyService(config.Bet, config.Id)
	if err != nil {
		log.Errorf("action: init_agency_service | result: fail | client_id: %v | error: %v", config.Id, err)
		return nil
	}
	client := &Client{
		config:        config,
		connManager:   *network.NewConnectionManager(),
		betHandler:    *protocol.NewBetHandler(),
		agencyService: *agencyService,
	}
	return client
}

func (c *Client) Run() error {

	c.setupGracefulShutdown()

	c.connInterface, err = c.connManager.Connect(c.config.ServerAddress)

	if err != nil {
		log.Infof("action: connect | result: fail | client_id: %v", c.config.Id)
		return err
	}

	bet, err := c.agencyService.ReadBets()

	if err != nil {
		log.Errorf("action: send_bets | result: fail | client_id: %v | error: %v", c.config.Id, err)
		c.Shutdown()
		return err
	}

	err = c.betHandler.SendBets(bet, c.connInterface)
	if err != nil {
		log.Errorf("action: send_bets | result: fail | client_id: %v | error: %v", c.config.Id, err)
		c.Shutdown()
		return err
	}

	err = c.betHandler.RecvConfirmation(c.connInterface)
	if err != nil {
		log.Errorf("action: recv_confirmation | result: fail | client_id: %v | error: %v", c.config.Id, err)
		c.Shutdown()
		return err
	}

	log.Infof("action: transmission finished | result: success | client_id: %v", c.config.Id)
	c.Shutdown()
	return nil
}

func (c *Client) setupGracefulShutdown() {
	/// This is a graceful non-blocking setup to shut down the process in case
	sigChannel := make(chan os.Signal, 1)
	signal.Notify(sigChannel, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		<-sigChannel
		log.Infof("action: shutdown_signal | result: received")
		c.Shutdown()
		os.Exit(0)
	}()
}

func (c *Client) Shutdown() {
	time.Sleep(100 * time.Millisecond)
	c.connInterface.Close()
	log.Infof("action: exit | result: success")
}
