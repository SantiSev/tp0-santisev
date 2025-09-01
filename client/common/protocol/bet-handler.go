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
// [bet length] [bet data] ;
// [bet length] [bet data] ;
// . . .
// EOF header

var log = logging.MustGetLogger("log")

const HEADER = "\x02"
const EOF = "\xFF"
const SUCCESS_HEADER = "\x01"
const MAX_DATA_LEN = 255
const SUCCESS_MESSAGE_SIZE = 4
const MAX_BATCH_SIZE = 8192 // 8 kBFAIL = b"\xff"

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

		batchLen := len(betBatchMessage)

		if batchLen > MAX_DATA_LEN {
			return fmt.Errorf("batch length too large for single byte: %d (max: %d)", batchLen, MAX_DATA_LEN)
		}
		lenBytes := make([]byte, 4)
		binary.BigEndian.PutUint32(lenBytes, uint32(batchLen))
		err = connSock.SendData(lenBytes)
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
	// Read the header (1 byte)
	headerData := make([]byte, len(SUCCESS_HEADER))
	err := connSock.ReceiveData(headerData)
	if err != nil {
		return fmt.Errorf("failed to receive header: %v", err)
	}

	success_header := string(headerData)

	if success_header == SUCCESS_HEADER {
		countBytes := make([]byte, 4)
		err = connSock.ReceiveData(countBytes)
		if err != nil {
			return fmt.Errorf("failed to receive bet count: %v", err)
		}

		betCount := int(binary.BigEndian.Uint32(countBytes))

		log.Infof("action: batch_confirmation | result: success | bets_processed: %d", betCount)

	} else {
		log.Errorf("action: batch_confirmation | result: fail | header: %x", headerData)
		return fmt.Errorf("batch processing failed")
	}

	return nil
}
