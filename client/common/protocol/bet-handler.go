package protocol

import (
	"bytes"
	"fmt"

	"github.com/7574-sistemas-distribuidos/docker-compose-init/client/common/network" // odio los imports de golang :D
	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("log")

type Bet struct {
	AgencyId  int64
	FirstName string
	LastName  string
	Document  int64
	Birthdate string
	Number    int64
}

type BetConfirmation struct {
	Status string
	Error  string
}

func NewBet(agencyId int64, firstName, lastName string, document int64, birthdate string, number int64) *Bet {
	return &Bet{
		AgencyId:  agencyId,
		FirstName: firstName,
		LastName:  lastName,
		Document:  document,
		Birthdate: birthdate,
		Number:    number,
	}
}

func (b *Bet) to_string() (string, error) {
	return fmt.Sprintf("%d,%s,%s,%d,%s,%d", b.AgencyId, b.FirstName, b.LastName, b.Document, b.Birthdate, b.Number), nil
}

type BetHandler struct {
}

func NewBetHandler() *BetHandler {
	return &BetHandler{}
}

func (b *BetHandler) SendBet(bet Bet, connSock *network.ConnectionSocket) error {
	data, err := bet.to_string()
	if err != nil {
		return err
	}
	err = connSock.SendData([]byte(data))
	return err
}

func (b *BetHandler) RecvBetConfirmation(connSock *network.ConnectionSocket) error {
	data := make([]byte, 256)

	n, err := connSock.ReceiveData(data)
	if err != nil {
		return err
	}

	response := string(bytes.TrimRight(data[:n], "\x00"))

	log.Infof("Received response: %s", response)

	return nil
}
