package protocol

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"os"
	"strings"

	"github.com/7574-sistemas-distribuidos/docker-compose-init/client/common/network"
	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("log")

const HEADER = "\x02"
const EOF = "\xFF"
const SUCCESS_HEADER_SIZE = 1
const WINNER_HEADER_SIZE = 1
const SUCCESS_HEADER = "\x01"
const SUCCESS_MESSAGE_SIZE = 4
const WINNERS_HEADER = "\x03"
const MAX_BATCH_SIZE = 8192 // 8 kB
const WINNER_COUNT_SIZE = 2

type BetHandler struct {
	MaxBatchAmount int
	bets           map[int16][]string
}

func NewBetHandler(maxBatchAmount int) *BetHandler {
	return &BetHandler{
		MaxBatchAmount: maxBatchAmount,
		bets:           make(map[int16][]string),
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

		betBatchMessage, canRead, err := b._createBetBatch(agency_id, scanner)
		if err != nil {
			return err
		}

		batchLen := len(betBatchMessage)

		if batchLen > MAX_BATCH_SIZE {
			return fmt.Errorf("batch length too large for single byte: %d (max: %d)", batchLen, MAX_BATCH_SIZE)
		}

		if !canRead || batchLen == 0 {
			break
		}

		err = connSock.SendData([]byte(HEADER))
		if err != nil {
			return err
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

func (b *BetHandler) GetLotteryResults(connSock *network.ConnectionInterface) error {
	headerData := make([]byte, WINNER_HEADER_SIZE)
	err := connSock.ReceiveData(headerData)
	if err != nil {
		return fmt.Errorf("failed to receive header: %v", err)
	}
	success_header := string(headerData)

	if success_header == WINNERS_HEADER {
		log.Infof("action: winner_confirmation | result: success")
		winnerCountBytes := make([]byte, WINNER_COUNT_SIZE)
		err := connSock.ReceiveData(winnerCountBytes)
		if err != nil {
			return fmt.Errorf("failed to receive length data: %v", err)
		}
		winnerCount := binary.BigEndian.Uint16(winnerCountBytes)

		log.Infof("action: consulta_ganadores | result: success | cant_ganadores: %d", winnerCount)

	} else {
		log.Errorf("action: winner_confirmation | result: fail")
		return fmt.Errorf("an error occurred tallying up the winners")
	}

	return nil
}

func (b *BetHandler) _createBetBatch(agency_id int64, scanner *bufio.Scanner) (string, bool, error) {
	var betBatchMessage string

	for i := 0; i < b.MaxBatchAmount; i++ {

		canRead := scanner.Scan()

		if err := scanner.Err(); err != nil {
			return "", false, fmt.Errorf("error reading file: %v", err)
		}

		if !canRead {
			return betBatchMessage, false, nil
		}

		line := strings.TrimSpace(scanner.Text())

		if line == "" {
			continue
		}
		betMessage := fmt.Sprintf("%d,%s,", agency_id, line)
		betBatchMessage += betMessage
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
