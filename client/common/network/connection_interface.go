package network

import (
	"io"
	"net"
)

type ConnectionInterface struct {
	conn net.Conn
}

func NewConnectionInterface(conn net.Conn) *ConnectionInterface {
	return &ConnectionInterface{
		conn: conn,
	}
}

func (c *ConnectionInterface) Connect(serverAddr string, connID string) error {
	log.Infof("connecting . . . ")
	conn, err := net.Dial("tcp", serverAddr)
	log.Infof("connection established, starting bet transmission")
	if err != nil {
		log.Criticalf(
			"action: connect | result: fail | client_id: %v | error: %v",
			connID,
			err,
		)
		return err
	}
	c.conn = conn
	return nil
}
func (c *ConnectionInterface) ReceiveData(buffer []byte) error {

	_, err := io.ReadFull(c.conn, buffer)
	if err != nil {
		log.Criticalf(
			"action: receive | result: fail | expected: %d bytes | error: %v",
			len(buffer), err,
		)
		return err
	}
	log.Debugf("action: receive | result: success | bytes: %d", len(buffer))
	return nil
}

func (c *ConnectionInterface) SendData(data []byte) error {
	log.Debugf("action: send | preparing to send | bytes: %d", len(data))
	totalWritten := 0
	for totalWritten < len(data) {
		n, err := c.conn.Write(data[totalWritten:])
		if err != nil {
			log.Criticalf(
				"action: send | result: fail | written: %d/%d | error: %v",
				totalWritten, len(data), err,
			)
			return err
		}
		totalWritten += n
	}
	log.Debugf("action: send | result: success | bytes: %d", totalWritten)
	return nil
}

func (c *ConnectionInterface) Close() error {
	if c.conn != nil {
		err := c.conn.Close()
		if err != nil {
			log.Criticalf(
				"action: close | result: fail | error: %v",
				err,
			)
			return err
		}
	}
	log.Infof("action: close | result: success")
	return nil
}
