package network

import (
	"io"
	"net"
)

type ConnectionInterface struct {
	conn net.Conn
}

func NewConnectionInterface() *ConnectionInterface {
	return &ConnectionInterface{}
}

func (c *ConnectionInterface) Connect(serverAddr string) error {
	conn, err := net.Dial("tcp", serverAddr)
	if err != nil {
		log.Criticalf(
			"action: connect | result: fail | error: %v",
			err,
		)
		return err
	}
	log.Debugf("action: connect | result: success | server: %s", serverAddr)
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
	log.Debugf("action: close | result: success")
	return nil
}
