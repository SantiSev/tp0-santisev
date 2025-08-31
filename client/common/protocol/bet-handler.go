package protocol

import (
	"fmt"
	"strings"

	"github.com/7574-sistemas-distribuidos/docker-compose-init/client/common/network" // odio los imports de golang :D
	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("log")

const HEADER = "\x02"
const EOF = "\xFF"
const SUCCESS_HEADER = "\x01"
const SUCCESS_MESSAGE_SIZE = 64
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

	err := connSock.SendData([]byte(HEADER))
	if err != nil {
		return err
	}

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
	data := make([]byte, len(SUCCESS_HEADER))

	err := connSock.ReceiveData(data)
	if err != nil {
		return err
	}

	success_header := string(data)
	if success_header == SUCCESS_HEADER {
		success_message := make([]byte, SUCCESS_MESSAGE_SIZE)
		err = connSock.ReceiveData(success_message)
		if err != nil {
			return err
		}
		successMsgStr := string(success_message)
		parts := strings.Split(successMsgStr, ",")
		if len(parts) != 2 {
			return fmt.Errorf("invalid bet confirmation format: %s", successMsgStr)
		}
		log.Infof("action: apuesta_enviada | result: success | dni: %s | numero: %s", parts[0], parts[1])
	} else {
		log.Info("Bet confirmation: FAIL")
	}

	return nil
}
