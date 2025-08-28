package common

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/7574-sistemas-distribuidos/docker-compose-init/client/common/protocol"
	"github.com/op/go-logging"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

// ClientConfig Configuration used by the client
type ClientConfig struct {
	Id            string
	ServerAddress string
	LoopAmount    int
	LoopPeriod    time.Duration
	Bet           *protocol.Bet
	LogLevel      string
}

func (c *ClientConfig) createBetFromEnv() error {
	agencyID, err := strconv.ParseInt(os.Getenv("CLI_ID"), 10, 64)
	if err != nil {
		log.Errorf("action: parse_ID | result: fail | client_id: %v | error: %v", c.Id, err)
		return err
	}
	betNumber, err := strconv.ParseInt(os.Getenv("CLI_BET_NUMBER"), 10, 64)
	if err != nil {
		log.Errorf("action: parse_bet_number | result: fail | client_id: %v | error: %v", c.Id, err)
		return err
	}
	clientDocument, err := strconv.ParseInt(os.Getenv("CLI_DOCUMENT"), 10, 64)
	if err != nil {
		log.Errorf("action: parse_client_document | result: fail | client_id: %v | error: %v", c.Id, err)
		return err
	}

	c.Bet = protocol.NewBet(
		agencyID,
		os.Getenv("CLI_FIRST_NAME"),
		os.Getenv("CLI_LAST_NAME"),
		clientDocument,
		os.Getenv("CLI_BIRTHDATE"),
		betNumber,
	)

	return nil
}

func (c *ClientConfig) PrintConfig() {
	log.Infof("action: config | result: success | client_id: %s | server_address: %s | loop_amount: %v | loop_period: %v | log_level: %s | bet info: %s",
		c.Id,
		c.ServerAddress,
		c.LoopAmount,
		c.LoopPeriod,
		c.LogLevel,
		c.Bet.To_string(),
	)
}

// InitConfig Function that uses viper library to parse configuration parameters.
// Viper is configured to read variables from both environment variables and the
// config file ./config.yaml. Environment variables takes precedence over parameters
// defined in the configuration file. If some of the variables cannot be parsed,
// an error is returned
func InitConfig() (*ClientConfig, error) {
	v := viper.New()
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	v.SetConfigFile("./config.yaml")
	if err := v.ReadInConfig(); err != nil {
		return nil, errors.Wrap(err, "Failed to read config.yaml")
	}

	if _, err := os.Stat(".env"); err == nil {
		envViper := viper.New()
		envViper.SetConfigFile(".env")
		if err := envViper.ReadInConfig(); err != nil {
			fmt.Printf("Warning: .env file exists but could not be read: %v\n", err)
		} else {
			fmt.Printf(".env file loaded successfully.\n")
			// Set CLI_ environment variables from .env file
			for key, value := range envViper.AllSettings() {
				if strings.HasPrefix(strings.ToUpper(key), "CLI_") {
					os.Setenv(strings.ToUpper(key), fmt.Sprintf("%v", value))
				}
			}
		}
	}

	if _, err := time.ParseDuration(v.GetString("loop.period")); err != nil {
		return nil, errors.Wrapf(err, "Could not parse loop.period '%s' as time.Duration", v.GetString("loop.period"))
	}

	// Create client configuration
	clientConfig := &ClientConfig{
		ServerAddress: v.GetString("server.address"),
		Id:            os.Getenv("CLI_ID"),
		LoopAmount:    v.GetInt("loop.amount"),
		LoopPeriod:    v.GetDuration("loop.period"),
		LogLevel:      v.GetString("log.level"),
	}

	if err := clientConfig.createBetFromEnv(); err != nil {
		return nil, errors.Wrap(err, "Failed to create bet from environment variables")
	}

	return clientConfig, nil
}

// InitLogger Receives the log level to be set in go-logging as a string. This method
// parses the string and set the level to the logger. If the level string is not
// valid an error is returned
func InitLogger(logLevel string) error {
	baseBackend := logging.NewLogBackend(os.Stdout, "", 0)
	format := logging.MustStringFormatter(
		`%{time:2006-01-02 15:04:05} %{level:.5s}     %{message}`,
	)
	backendFormatter := logging.NewBackendFormatter(baseBackend, format)

	backendLeveled := logging.AddModuleLevel(backendFormatter)
	logLevelCode, err := logging.LogLevel(logLevel)
	if err != nil {
		return err
	}
	backendLeveled.SetLevel(logLevelCode, "")

	// Set the backends to be used.
	logging.SetBackend(backendLeveled)
	return nil
}
