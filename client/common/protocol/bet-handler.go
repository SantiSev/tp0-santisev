package protocol

import (
	"encoding/binary"
	"fmt"
	"strings"

	"github.com/7574-sistemas-distribuidos/docker-compose-init/client/common/network"
	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("log")

type BetHandler struct {
	MaxBatchAmount int
}

func NewBetHandler() *BetHandler {
	return &BetHandler{}
}

func (b *BetHandler) SendBets(bet string, connSock *network.ConnectionInterface) error {

	err := connSock.SendData([]byte(HEADER))
	if err != nil {
		return err
	}

	lenBytes := make([]byte, 2)
	binary.BigEndian.PutUint16(lenBytes, uint16(len(bet)))
	err = connSock.SendData(lenBytes)
	if err != nil {
		return err
	}
	err = connSock.SendData([]byte(bet))
	if err != nil {
		return err
	}
	log.Debug("action: send_bet_batch | result: success | batch_length: %d Bytes", len(bet))

	return err
}

func (b *BetHandler) RecvConfirmation(connSock *network.ConnectionInterface) error {
	headerData := make([]byte, SUCCESS_HEADER_SIZE)
	err := connSock.ReceiveData(headerData)
	if err != nil {
		return fmt.Errorf("failed to receive header: %v", err)
	}
	success_header := string(headerData)

	if success_header == SUCCESS_HEADER {

		lengthData := make([]byte, 1)
		err = connSock.ReceiveData(lengthData)

		if err != nil {
			return fmt.Errorf("failed to receive message length: %v", err)
		}

		messageLength := int(lengthData[0])
		messageData := make([]byte, messageLength)
		err = connSock.ReceiveData(messageData)

		if err != nil {
			return fmt.Errorf("failed to receive message data: %v", err)
		}

		messageParts := strings.Split(string(messageData), ",")

		log.Debug("action: apuesta_enviada | result: success | dni: %s | numero: %s", messageParts[0], messageParts[1])
	} else {
		log.Errorf("action: batch_confirmation | result: fail")
		return fmt.Errorf("batch processing failed")
	}

	return nil
}
