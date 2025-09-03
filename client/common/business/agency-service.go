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
	agency_id      uint8
	agencyFile     *os.File
	has_data       bool
	scanner        *bufio.Scanner
	maxBatchAmount int
}

func NewAgencyService(agencyFilePath string, maxBatchAmount int, agency_id uint8) (*AgencyService, error) {
	file, err := os.Open(agencyFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	return &AgencyService{
		agency_id:      agency_id,
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

		if line == "" {
			continue
		}
		betMessage := fmt.Sprintf("%d,%s,", a.agency_id, line)
		betBatchMessage += betMessage
	}
	return betBatchMessage, nil
}

func (a *AgencyService) HasData() bool {
	return a.has_data
}

func (a *AgencyService) ShowResults(results string) {

	amountWinners := len(strings.Split(results, ","))

	log.Infof("action: consulta_ganadores | result: success | cant_ganadores: %d:", amountWinners)
	log.Infof("action: mostrar_ganadores | result: success | ganadores: %s:", results)
}
