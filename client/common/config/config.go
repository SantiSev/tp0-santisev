package config

import (
	"os"
	"strings"

	"github.com/7574-sistemas-distribuidos/docker-compose-init/client/common/client"
	"github.com/op/go-logging"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

func InitConfig() (*client.ClientConfig, error) {
	loadEnvVars()

	v := viper.New()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	configPath := os.Getenv("CLI_CONFIG_FILEPATH")

	v.SetConfigFile(configPath)
	if err := v.ReadInConfig(); err != nil {
		return nil, errors.Wrapf(err, "failed to read config file %s", configPath)
	}

	clientConfig := &client.ClientConfig{
		ServerAddress:  v.GetString("server.address"),
		Id:             uint8(v.GetInt("CLI_ID")),
		LoopAmount:     v.GetInt("loop.amount"),
		LoopPeriod:     v.GetDuration("loop.period"),
		LogLevel:       v.GetString("log.level"),
		MaxBatchAmount: v.GetInt("batch.maxAmount"),
		AgencyFilePath: v.GetString("CLI_AGENCY_FILEPATH"),
	}

	return clientConfig, nil
}

func loadEnvVars() {
	// if the env vars are present, there is no need to load .env
	if os.Getenv("CLI_ID") != "" &&
		os.Getenv("CLI_AGENCY_FILEPATH") != "" &&
		os.Getenv("CLI_CONFIG_FILEPATH") != "" {
		return
	}

	if _, err := os.Stat(".env"); err == nil {
		envViper := viper.New()
		envViper.SetConfigFile(".env")
		if err := envViper.ReadInConfig(); err == nil {
			for _, key := range envViper.AllKeys() {
				envKey := strings.ToUpper(key)
				if os.Getenv(envKey) == "" {
					os.Setenv(envKey, envViper.GetString(key))
				}
			}
		}
	}
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

	logging.SetBackend(backendLeveled)
	return nil
}
