package client

import (
	"time"
)

type ClientConfig struct {
	Id             uint8
	ServerAddress  string
	LoopAmount     int
	LoopPeriod     time.Duration
	LogLevel       string
	MaxBatchAmount int
	AgencyFilePath string
}

func (c *ClientConfig) PrintConfig() {
	log.Infof("action: config | result: success | client_id: %d | server_address: %s | loop_amount: %v | loop_period: %v | log_level: %s | max_batch_amount: %d | agency_file_path: %s",
		c.Id,
		c.ServerAddress,
		c.LoopAmount,
		c.LoopPeriod,
		c.LogLevel,
		c.MaxBatchAmount,
		c.AgencyFilePath,
	)
}
