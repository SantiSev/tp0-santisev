package client

type ClientConfig struct {
	Id            uint8
	ServerAddress string
	LogLevel      string
	Bet           string
}

func (c *ClientConfig) PrintConfig() {
	log.Infof("action: config | result: success | client_id: %d | server_address: %s | log_level: %s | bet: %s",
		c.Id,
		c.ServerAddress,
		c.LogLevel,
		c.Bet,
	)
}
