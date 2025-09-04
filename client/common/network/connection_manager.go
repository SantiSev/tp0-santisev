package network

import (
	"fmt"
	"time"

	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("log")

const CONNECTION_RETRIES = 5
const CONNECTION_SLEEP_AMOUNT = 1000

type ConnectionManager struct {
}

func NewConnectionManager() *ConnectionManager {
	return &ConnectionManager{}
}

func (c *ConnectionManager) Connect(serverAddr string) (*ConnectionInterface, error) {

	for attempt := 1; attempt <= CONNECTION_RETRIES; attempt++ {

		connSocket := NewConnectionInterface()
		err := connSocket.Connect(serverAddr)
		if err == nil {
			return connSocket, nil
		}

		log.Warningf("action: connect | result: fail | attempt: %d/%d | server: %s",
			attempt, CONNECTION_RETRIES, serverAddr)

		sleepDuration := time.Duration(CONNECTION_SLEEP_AMOUNT)
		log.Infof("action: retry_connect | duration: %v", sleepDuration)
		time.Sleep(sleepDuration)
	}

	return nil, fmt.Errorf("failed to connect to %s", serverAddr)
}
