package protocol

import "fmt"

type Bet struct {
	AgencyId  int64
	FirstName string
	LastName  string
	Document  int64
	Birthdate string
	Number    int64
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

func (b *Bet) To_string() string {
	return fmt.Sprintf("%d,%s,%s,%d,%s,%d", b.AgencyId, b.FirstName, b.LastName, b.Document, b.Birthdate, b.Number)
}
