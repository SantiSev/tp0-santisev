package protocol

import (
	"fmt"
	"strings"

	"github.com/7574-sistemas-distribuidos/docker-compose-init/client/common/network"
	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("log")

type BetHandler struct {
	MaxBatchAmount int
}

func NewAgencyHandler() *BetHandler {
	return &BetHandler{}
}

func (b *BetHandler) SendBets(bet string, connSock *network.ConnectionInterface) error {

	// Send header
	if err := connSock.SendData([]byte(HEADER)); err != nil {
		return err
	}

	// Send length as single byte
	betBytes := []byte(bet)
	if err := connSock.SendData([]byte{byte(len(betBytes))}); err != nil {
		return err
	}

	// Send bet data
	if err := connSock.SendData(betBytes); err != nil {
		return err
	}

	log.Debugf("action: send_bet | result: success | batch_length: %d Bytes", len(betBytes))
	return nil
}

func (b *BetHandler) RecvConfirmation(connSock *network.ConnectionInterface) error {
	headerData := make([]byte, SUCCESS_HEADER_SIZE)
	err := connSock.ReceiveData(headerData)
	if err != nil {
		return fmt.Errorf("failed to receive header: %v", err)
	}
	success_header := string(headerData)

	if success_header == SUCCESS_HEADER {

		lengthData := make([]byte, RESPONSE_DATA_SIZE)
		if err := connSock.ReceiveData(lengthData); err != nil {
			return fmt.Errorf("failed to receive message length: %v", err)
		}

		messageData := make([]byte, lengthData[0])
		if err := connSock.ReceiveData(messageData); err != nil {
			return fmt.Errorf("failed to receive message data: %v", err)
		}

		parts := strings.Split(string(messageData), ",")
		if len(parts) >= 2 {
			log.Infof("action: apuesta_enviada | result: success | dni: %s | numero: %s", parts[0], parts[1])
		}
	} else {
		log.Errorf("action: batch_confirmation | result: fail")
		return fmt.Errorf("batch processing failed")
	}

	return nil
}
