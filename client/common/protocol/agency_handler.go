package protocol

import (
	"encoding/binary"
	"fmt"

	"github.com/7574-sistemas-distribuidos/docker-compose-init/client/common/network"
	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("log")

type AgencyHandler struct {
	MaxBatchAmount int
}

func NewAgencyHandler() *AgencyHandler {
	return &AgencyHandler{}
}

func (a *AgencyHandler) SendBets(bets string, connSock *network.ConnectionInterface) error {

	if len(bets) > MAX_BATCH_SIZE {
		return fmt.Errorf("bets size too big to send: %d bytes (max %d)", len(bets), MAX_BATCH_SIZE)
	}

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

func (a *AgencyHandler) RecvConfirmation(connSock *network.ConnectionInterface) error {
	headerData := make([]byte, SUCCESS_HEADER_SIZE)
	err := connSock.ReceiveData(headerData)
	if err != nil {
		return fmt.Errorf("failed to receive header: %v", err)
	}
	success_header := string(headerData)

	if success_header == SUCCESS_HEADER {
		log.Debug("action: batch_confirmation | result: success")
	} else {
		log.Errorf("action: batch_confirmation | result: fail")
		return fmt.Errorf("batch processing failed")
	}

	return nil
}

func (a *AgencyHandler) SendDone(connSock *network.ConnectionInterface) error {
	err := connSock.SendData([]byte(EOF))
	return err
}

func (a *AgencyHandler) GetResults(connSock *network.ConnectionInterface) (string, error) {
	log.Debug("action: waiting_lottery_results | result: success")
	headerData := make([]byte, WINNER_HEADER_SIZE)
	err := connSock.ReceiveData(headerData)
	if err != nil {
		return "", fmt.Errorf("failed to receive header: %v", err)
	}
	success_header := string(headerData)

	log.Debug("action: lottery_results_received | result: success")

	if success_header == WINNERS_HEADER {
		winners_data_count := make([]byte, WINNER_COUNT_SIZE)
		err := connSock.ReceiveData(winners_data_count)
		if err != nil {
			return "", fmt.Errorf("failed to receive length data: %v", err)
		}
		winnerCount := binary.BigEndian.Uint16(winners_data_count)

		winnerData := make([]byte, winnerCount)
		err = connSock.ReceiveData(winnerData)
		if err != nil {
			return "", fmt.Errorf("failed to receive winner data: %v", err)
		}
		winnerStr := string(winnerData)

		return winnerStr, nil

	} else {
		log.Errorf("action: winner_confirmation | result: fail")
		return "", fmt.Errorf("an error occurred tallying up the winners")
	}
}
