package business

import (
	"fmt"
	"strings"
)

type AgencyService struct {
	bet       string
	agency_id uint8
}

func NewAgencyService(bet string, agency_id uint8) (*AgencyService, error) {

	return &AgencyService{
		bet:       bet,
		agency_id: agency_id,
	}, nil
}

func (a *AgencyService) ReadBets() (string, error) {

	if !a.validateBet() {
		return "", fmt.Errorf("invalid bet format: %s", a.bet)
	}

	return a.bet, nil
}

func (a *AgencyService) validateBet() bool {

	bets := strings.Split(a.bet, ",")

	return len(bets) == 6
}
