package business

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("log")

type AgencyService struct {
	agencyFile     *os.File
	has_data       bool
	scanner        *bufio.Scanner
	maxBatchAmount int
}

func NewAgencyService(agencyFilePath string, maxBatchAmount int) (*AgencyService, error) {
	file, err := os.Open(agencyFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %v", err)
	}

	return &AgencyService{
		agencyFile:     file,
		has_data:       true,
		scanner:        bufio.NewScanner(file),
		maxBatchAmount: maxBatchAmount,
	}, nil
}

func (a *AgencyService) ReadBets(batchSize int) (string, error) {

	var betBatchMessage string

	for i := 0; i < a.maxBatchAmount; i++ {

		canRead := a.scanner.Scan()

		if err := a.scanner.Err(); err != nil {
			return betBatchMessage, fmt.Errorf("error reading file: %v", err)
		}

		if !canRead {
			a.has_data = false
			return betBatchMessage, nil
		}

		line := strings.TrimSpace(a.scanner.Text())

		if line == "" || !is_valid_bet(line) {
			continue
		}

		betMessage := fmt.Sprintf("%s,", line)
		betBatchMessage += betMessage
	}
	return betBatchMessage, nil
}

func is_valid_bet(bet string) bool {
	parts := strings.Split(bet, ",")
	return len(parts) == 5
}

func (a *AgencyService) HasData() bool {
	return a.has_data
}

func (a *AgencyService) ShowResults(results string) {

	var amountWinners int
	if strings.TrimSpace(results) == "" {
		amountWinners = 0
	} else {
		amountWinners = len(strings.Split(results, ","))
	}

	log.Infof("action: consulta_ganadores | result: success | cant_ganadores: %d", amountWinners)
	log.Infof("action: mostrar_ganadores | result: success | ganadores: %s", results)
}

func (a *AgencyService) Close() error {
	if a.agencyFile != nil {
		return a.agencyFile.Close()
	}
	return nil
}
