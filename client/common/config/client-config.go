package config

import (
	"time"
)

type ClientConfig struct {
	Id             int64
	ServerAddress  string
	LoopAmount     int
	LoopPeriod     time.Duration
	LogLevel       string
	MaxBatchAmount int
}

func (c *ClientConfig) PrintConfig() {
	log.Infof("action: config | result: success | client_id: %s | server_address: %s | loop_amount: %v | loop_period: %v | log_level: %s | max_batch_amount: %d",
		c.Id,
		c.ServerAddress,
		c.LoopAmount,
		c.LoopPeriod,
		c.LogLevel,
		c.MaxBatchAmount,
	)
}
