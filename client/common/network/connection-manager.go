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

func (c *ConnectionManager) Connect(serverAddr string, connID string) *ConnectionSocket {

	connSocket := NewConnectionSocket(nil)
	connSocket.ConnectClient(serverAddr, connID)
	return connSocket
}
