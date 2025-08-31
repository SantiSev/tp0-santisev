package network

import (
	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("log")

type ConnectionManager struct {
}

func NewConnectionManager() *ConnectionManager {
	return &ConnectionManager{}
}

func (c *ConnectionManager) Connect(serverAddr string) (*ConnectionInterface, error) {

	connSocket := NewConnectionInterface(nil)
	err := connSocket.Connect(serverAddr)
	if err != nil {
		return nil, err
	}
	return connSocket, nil
}
