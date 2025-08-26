package network

import "net"

type ConnectionSocket struct {
	conn net.Conn
}

func NewConnectionSocket(conn net.Conn) *ConnectionSocket {
	return &ConnectionSocket{
		conn: conn,
	}
}

func (c *ConnectionSocket) ConnectClient(serverAddr string, connID string) error {
	conn, err := net.Dial("tcp", serverAddr)
	if err != nil {
		log.Criticalf(
			"action: connect | result: fail | client_id: %v | error: %v",
			connID,
			err,
		)
	}
	c.conn = conn
	return nil
}

func (c *ConnectionSocket) ReceiveData(buffer []byte) (int, error) {
	n, err := c.conn.Read(buffer)
	if err != nil {
		log.Criticalf(
			"action: receive | result: fail | error: %v",
			err,
		)
		return n, err
	}
	return n, nil
}

func (c *ConnectionSocket) SendData(data []byte) error {
	_, err := c.conn.Write(data)
	if err != nil {
		log.Criticalf(
			"action: send | result: fail | error: %v",
			err,
		)
		return err
	}
	return nil
}

func (c *ConnectionSocket) Close() error {
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
