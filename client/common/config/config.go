package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/op/go-logging"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

var log = logging.MustGetLogger("log")

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
			if cliId := envViper.GetString("CLI_ID"); cliId != "" {
				os.Setenv("CLI_ID", cliId)
			}
		}
	}

	if _, err := time.ParseDuration(v.GetString("loop.period")); err != nil {
		return nil, errors.Wrapf(err, "Could not parse loop.period '%s' as time.Duration", v.GetString("loop.period"))
	}

	cliIdStr := os.Getenv("CLI_ID")
	var cliId int64
	if cliIdStr != "" {
		var err error
		cliId, err = strconv.ParseInt(cliIdStr, 10, 64)
		if err != nil {
			return nil, errors.Wrapf(err, "Failed to parse CLI_ID '%s' as int64", cliIdStr)
		}
	}

	clientConfig := &ClientConfig{
		ServerAddress:  v.GetString("server.address"),
		Id:             cliId,
		LoopAmount:     v.GetInt("loop.amount"),
		LoopPeriod:     v.GetDuration("loop.period"),
		LogLevel:       v.GetString("log.level"),
		MaxBatchAmount: v.GetInt("max.batch.amount"),
	}

	return clientConfig, nil
}

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
