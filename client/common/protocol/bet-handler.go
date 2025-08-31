package protocol

import (
	"bufio"
	"fmt"
	"os"
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
const MAX_FILE_SIZE = 8192

type BetHandler struct {
	MaxBatchAmount int
}

func NewBetHandler(maxBatchAmount int) *BetHandler {
	return &BetHandler{
		MaxBatchAmount: maxBatchAmount,
	}
}
func (b *BetHandler) SendAllBetData(agency_id int64, agency_data_file string, connSock *network.ConnectionInterface) error {

	fileInfo, err := os.Stat(agency_data_file)
	if err != nil {
		return fmt.Errorf("failed to get file info: %v", err)
	}

	if fileInfo.Size() > MAX_FILE_SIZE {
		return fmt.Errorf("file too large: %d bytes | the max amount to send is 8 Kb", fileInfo.Size())
	}

	err = connSock.SendData([]byte(HEADER))
	if err != nil {
		return err
	}

	b._sendBatch(agency_id, agency_data_file, connSock)
	return err
}

func (b *BetHandler) _sendBatch(agency_id int64, agency_data_file string, connSock *network.ConnectionInterface) error {
	file, err := os.Open(agency_data_file)
	if err != nil {
		return fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for i := 0; i < b.MaxBatchAmount && scanner.Scan(); i++ {
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines
		if line == "" {
			i--
			continue
		}

		bet, err := _processBet(line, agency_id)
		if err != nil {
			log.Errorf("action: parse_bet | result: fail | line: %s | error: %v", line, err)
			return err
		}

		err = b._sendBet(bet, connSock)
		if err != nil {
			log.Errorf("action: send_bet | result: fail | batch_index: %d | error: %v", i, err)
			return err
		}

		err = b._recvBetConfirmation(connSock)
		if err != nil {
			log.Errorf("action: recv_confirmation | result: fail | batch_index: %d | error: %v", i, err)
			return err
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading file: %v", err)
	}

	log.Infof("action: send_batch | result: success")
	return nil
}

func (b *BetHandler) SendDone(connSock *network.ConnectionInterface) error {
	err := connSock.SendData([]byte(EOF))
	return err
}

func _processBet(data string, agencyId int64) (*Bet, error) {
	var firstName, lastName, birthdate string
	var document, number int64

	_, err := fmt.Sscanf(data, "%d,%s,%s,%d,%s,%d", &agencyId, &firstName, &lastName, &document, &birthdate, &number)
	if err != nil {
		return nil, err
	}

	return NewBet(agencyId, firstName, lastName, document, birthdate, number), nil
}

func (b *BetHandler) _sendBet(bet *Bet, connSock *network.ConnectionInterface) error {
	betString := bet.To_string()

	betBytes := []byte(betString)
	if len(betBytes) > BET_DATA_SIZE {
		return fmt.Errorf("bet data too large: %d bytes", len(betBytes))
	}

	data := make([]byte, BET_DATA_SIZE)
	copy(data, betBytes)

	err := connSock.SendData(data)
	if err != nil {
		log.Errorf("action: send_bet | result: fail | error: %v", err)
		return err
	}

	return nil
}

func (b *BetHandler) _recvBetConfirmation(connSock *network.ConnectionInterface) error {
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
