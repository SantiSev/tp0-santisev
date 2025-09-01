package protocol

import (
	"bufio"
	"encoding/binary"
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
// data size: 2 B
// data: N B in chunks less than 8kB
// . . .
// EOF header

var log = logging.MustGetLogger("log")

const HEADER = "\x02"
const EOF = "\xFF"
const SUCCESS_HEADER_SIZE = 1
const SUCCESS_HEADER = "\x01"
const SUCCESS_MESSAGE_SIZE = 4
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

	scanner := bufio.NewScanner(file)

	for {

		err = connSock.SendData([]byte(HEADER))
		if err != nil {
			return err
		}

		if err := scanner.Err(); err != nil {
			return fmt.Errorf("error reading file: %v", err)
		}

		betBatchMessage, canRead, err := b._createBetBatch(agency_id, scanner)
		if err != nil {
			return err
		}

		batchLen := len(betBatchMessage)

		if batchLen > MAX_BATCH_SIZE {
			return fmt.Errorf("batch length too large for single byte: %d (max: %d)", batchLen, MAX_BATCH_SIZE)
		}
		lenBytes := make([]byte, 2)
		binary.BigEndian.PutUint16(lenBytes, uint16(batchLen))
		log.Infof("action: send_bet_batch | result: success | batch_length: %dB", batchLen)
		err = connSock.SendData(lenBytes)
		if err != nil {
			return err
		}

		err = connSock.SendData([]byte(betBatchMessage))
		if err != nil {
			return err
		}

		if !canRead {
			break
		}

	}

	err = b._sendDone(connSock)
	if err != nil {
		log.Errorf("action: send_done | result: fail | error: %v", err)
		return err
	}

	err = b._recvConfirmation(connSock)
	if err != nil {
		return err
	}
	return err
}

func (b *BetHandler) _createBetBatch(agency_id int64, scanner *bufio.Scanner) (string, bool, error) {
	var betBatchMessage string
	counter := 0

	for counter < b.MaxBatchAmount {

		canRead := scanner.Scan()

		if !canRead {
			return betBatchMessage, false, nil
		}

		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines
		if line == "" {
			continue
		}

		betMessage := fmt.Sprintf("%d,%s,", agency_id, line)
		betBatchMessage += betMessage
		counter++
	}
	return betBatchMessage, true, nil
}

func (b *BetHandler) _sendDone(connSock *network.ConnectionInterface) error {
	err := connSock.SendData([]byte(EOF))
	return err
}

func (b *BetHandler) _recvConfirmation(connSock *network.ConnectionInterface) error {
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
