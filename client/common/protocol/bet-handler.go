package protocol

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/7574-sistemas-distribuidos/docker-compose-init/client/common/network" // odio los imports de golang :D
	"github.com/op/go-logging"
)

// TODO: SEND THE LENGTH AND THEN THE DATA ! ! ! DONT USE FIXED AMOUNTS
// this is due to using strings
// e.g:
// main header: 1 B
// [bet length] [bet data] ;
// [bet length] [bet data] ;
// . . .
// EOF header

var log = logging.MustGetLogger("log")

const HEADER = "\x02"
const EOF = "\xFF"
const SUCCESS_HEADER = "\x01"
const MAX_DATA_LEN = 255
const SUCCESS_MESSAGE_SIZE = 64
const MAX_BATCH_SIZE = 8192 // 8 kB

type BetHandler struct {
	MaxBatchAmount int
}

func NewBetHandler(maxBatchAmount int) *BetHandler {
	return &BetHandler{
		MaxBatchAmount: maxBatchAmount,
	}
}
func (b *BetHandler) SendAllBetData(agency_id int64, agency_data_file string, connSock *network.ConnectionInterface) error {

	file, err := os.Open(agency_data_file)
	if err != nil {
		return fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	// First, we send the header
	err = connSock.SendData([]byte(HEADER))
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {

		if err := scanner.Err(); err != nil {
			return fmt.Errorf("error reading file: %v", err)
		}

		betBatchMessage, err := b._createBetBatch(agency_id, scanner)
		if err != nil {
			return err
		}
		err = connSock.SendData([]byte(betBatchMessage))
		if err != nil {
			return err
		}
		err = b._recvConfirmation(connSock)
		if err != nil {
			return err
		}

	}

	err = b._sendDone(connSock)
	if err != nil {
		log.Errorf("action: send_done | result: fail | error: %v", err)
		return err
	}
	return err
}

func (b *BetHandler) _createBetBatch(agency_id int64, scanner *bufio.Scanner) (string, error) {
	var betBatchMessage string

	for i := 0; i < b.MaxBatchAmount && scanner.Scan(); i++ {
		i++
		line := strings.TrimSpace(scanner.Text())

		betMessage := fmt.Sprintf("%d,%s\n", agency_id, line)
		length := len(betMessage)
		if length > MAX_DATA_LEN {
			return "", fmt.Errorf("this line [%s] is too large: %d bytes", line, length)
		}
		betMessage = fmt.Sprintf("%c%s", length, betMessage)
		betBatchMessage += betMessage
	}
	log.Infof("action: create_bet_batch | result: success | batch data: %s", betBatchMessage)

	return betBatchMessage, nil
}

func (b *BetHandler) _sendDone(connSock *network.ConnectionInterface) error {
	err := connSock.SendData([]byte(EOF))
	return err
}

func (b *BetHandler) _recvConfirmation(connSock *network.ConnectionInterface) error {
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
