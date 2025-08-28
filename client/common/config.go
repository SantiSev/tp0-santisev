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
	ID            string
	ServerAddress string
	LoopAmount    int
	LoopPeriod    time.Duration
	Bet           *protocol.Bet
	LogLevel      string
}

func (c *ClientConfig) createBetFromEnv() error {
	agencyID, err := strconv.ParseInt(os.Getenv("CLI_AGENCY_ID"), 10, 64)
	if err != nil {
		log.Errorf("action: parse_agency_id | result: fail | client_id: %v | error: %v", c.ID, err)
		return err
	}
	betNumber, err := strconv.ParseInt(os.Getenv("CLI_BET_NUMBER"), 10, 64)
	if err != nil {
		log.Errorf("action: parse_bet_number | result: fail | client_id: %v | error: %v", c.ID, err)
		return err
	}
	clientDocument, err := strconv.ParseInt(os.Getenv("CLI_DOCUMENT"), 10, 64)
	if err != nil {
		log.Errorf("action: parse_client_document | result: fail | client_id: %v | error: %v", c.ID, err)
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
		c.ID,
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

	// Configure viper to read env variables with the CLI_ prefix
	v.AutomaticEnv()
	v.SetEnvPrefix("cli")
	// Use a replacer to replace env variables underscores with points. This let us
	// use nested configurations in the config file and at the same time define
	// env variables for the nested configurations
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Add env variables supported
	v.BindEnv("id")
	v.BindEnv("server", "address")
	v.BindEnv("loop", "period")
	v.BindEnv("loop", "amount")
	v.BindEnv("log", "level")

	// Add additional env variables from .env file
	v.BindEnv("agency_id")
	v.BindEnv("bet_number")
	v.BindEnv("document")
	v.BindEnv("first_name")
	v.BindEnv("last_name")
	v.BindEnv("birthdate")

	// Try to read configuration from config file. If config file
	// does not exists then ReadInConfig will fail but configuration
	// can be loaded from the environment variables so we shouldn't
	// return an error in that case
	v.SetConfigFile("./config.yaml")
	if err := v.ReadInConfig(); err != nil {
		fmt.Printf("Configuration could not be read from config file.")
	}

	// Parse time.Duration variables and return an error if those variables cannot be parsed

	if _, err := time.ParseDuration(v.GetString("loop.period")); err != nil {
		return nil, errors.Wrapf(err, "Could not parse CLI_LOOP_PERIOD env var as time.Duration.")
	}

	clientConfig := ClientConfig{
		ServerAddress: v.GetString("server.address"),
		ID:            v.GetString("id"),
		LoopAmount:    v.GetInt("loop.amount"),
		LoopPeriod:    v.GetDuration("loop.period"),
		LogLevel:      v.GetString("log.level"),
	}

	v.SetConfigFile(".env")
	if err := v.ReadInConfig(); err != nil {
		fmt.Printf("Configuration could not be read from config file.")
	}

	v.SetConfigFile(".env")
	if err := v.ReadInConfig(); err != nil {
		fmt.Printf("Configuration could not be read from .env file.")
	} else {
		// Set environment variables from .env file so os.Getenv() can access them
		for _, key := range v.AllKeys() {
			os.Setenv(strings.ToUpper(key), v.GetString(key))
		}
	}

	// Create the bet from environment variables after reading .env file
	if err := clientConfig.createBetFromEnv(); err != nil {
		return nil, errors.Wrap(err, "Failed to create bet from environment variables")
	}

	return &clientConfig, nil
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
