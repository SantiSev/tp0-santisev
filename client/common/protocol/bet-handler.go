package protocol

import (
	"fmt"

	"github.com/7574-sistemas-distribuidos/docker-compose-init/client/common/network" // odio los imports de golang :D
	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("log")

const HEADER = "\x02"
const EOF = "\xFF"
const SUCCESS = "\x00"
const BET_DATA_SIZE = 256

type BetHandler struct {
}

func NewBetHandler() *BetHandler {
	return &BetHandler{}
}
func (b *BetHandler) SendBet(bet Bet, connSock *network.ConnectionInterface) error {
	betString := bet.To_string()
	betBytes := []byte(betString)
	if len(betBytes) > BET_DATA_SIZE {
		return fmt.Errorf("bet data too large: %d bytes", len(betBytes))
	}

	// Send header byte first
	err := connSock.SendData([]byte(HEADER))
	if err != nil {
		return err
	}

	// Send 1024-byte data payload
	data := make([]byte, BET_DATA_SIZE)
	copy(data, betBytes)
	err = connSock.SendData(data)
	return err
}

func (b *BetHandler) SendDone(connSock *network.ConnectionInterface) error {
	err := connSock.SendData([]byte(EOF))
	return err
}

func (b *BetHandler) RecvBetConfirmation(connSock *network.ConnectionInterface) error {
	data := make([]byte, len(SUCCESS))

	err := connSock.ReceiveData(data)
	if err != nil {
		return err
	}

	response := string(data)
	if response == SUCCESS {
		log.Info("Bet confirmation: SUCCESS")
	} else {
		log.Info("Bet confirmation: FAIL")
	}

	return nil
}
