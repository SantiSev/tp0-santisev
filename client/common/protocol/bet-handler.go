package protocol

import (
	"encoding/binary"
	"fmt"

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

func (b *BetHandler) SendBets(bets string, connSock *network.ConnectionInterface) error {

	err := connSock.SendData([]byte(HEADER))
	if err != nil {
		return err
	}

	lenBytes := make([]byte, 2)
	binary.BigEndian.PutUint16(lenBytes, uint16(len(bets)))
	err = connSock.SendData(lenBytes)
	if err != nil {
		return err
	}
	err = connSock.SendData([]byte(bets))
	if err != nil {
		return err
	}
	log.Debug("action: send_bet_batch | result: success | batch_length: %d Bytes", len(bets))

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
		log.Infof("action: batch_confirmation | result: success")
	} else {
		log.Errorf("action: batch_confirmation | result: fail")
		return fmt.Errorf("batch processing failed")
	}

	return nil
}

func (b *BetHandler) SendDone(connSock *network.ConnectionInterface) error {
	err := connSock.SendData([]byte(EOF))
	return err
}

func (b *BetHandler) GetResults(connSock *network.ConnectionInterface) (string, error) {
	log.Debug("action: waiting_lottery_results | result: success")
	headerData := make([]byte, WINNER_HEADER_SIZE)
	err := connSock.ReceiveData(headerData)
	if err != nil {
		return "", fmt.Errorf("failed to receive header: %v", err)
	}
	success_header := string(headerData)

	if success_header == WINNERS_HEADER {
		log.Infof("action: winner_confirmation | result: success")
		winnerCountBytes := make([]byte, WINNER_COUNT_SIZE)
		err := connSock.ReceiveData(winnerCountBytes)
		if err != nil {
			return "", fmt.Errorf("failed to receive length data: %v", err)
		}
		winnerCount := binary.BigEndian.Uint16(winnerCountBytes)

		log.Infof("action: consulta_ganadores | result: success | cant_ganadores: %d", winnerCount)

	} else {
		log.Errorf("action: winner_confirmation | result: fail")
		return "", fmt.Errorf("an error occurred tallying up the winners")
	}

	return "", nil
}
