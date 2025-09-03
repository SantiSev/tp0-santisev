package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/7574-sistemas-distribuidos/docker-compose-init/client/common/client"
	"github.com/op/go-logging"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

const CONFIG_FILE_PATH = "./config.yaml"

func InitConfig() (*client.ClientConfig, error) {
	loadEnvVars()

	v := viper.New()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	configPath := CONFIG_FILE_PATH

	v.SetConfigFile(configPath)
	if err := v.ReadInConfig(); err != nil {
		return nil, errors.Wrapf(err, "failed to read config file %s", configPath)
	}

	client_id := v.GetInt("CLI_ID")
	bet_first_name := v.GetString("CLIENT_FIRST_NAME")
	bet_last_name := v.GetString("CLIENT_LAST_NAME")
	bet_birthdate := v.GetString("CLIENT_BIRTHDATE")
	bet_document := v.GetString("CLIENT_DOCUMENT")
	bet_number := v.GetInt("BET_NUMBER")

	betString := fmt.Sprintf("%d,%s,%s,%s,%s,%d", client_id, bet_first_name, bet_last_name, bet_document, bet_birthdate, bet_number)

	clientConfig := &client.ClientConfig{
		ServerAddress: v.GetString("server.address"),
		Id:            uint8(v.GetInt("CLI_ID")),
		LogLevel:      v.GetString("log.level"),
		Bet:           betString,
	}

	return clientConfig, nil
}

func loadEnvVars() {
	// if the env vars are present, there is no need to load .env
	if os.Getenv("CLI_ID") != "" &&
		os.Getenv("CLIENT_FIRST_NAME") != "" &&
		os.Getenv("CLIENT_LAST_NAME") != "" &&
		os.Getenv("CLIENT_BIRTHDATE") != "" &&
		os.Getenv("CLIENT_DOCUMENT") != "" &&
		os.Getenv("BET_NUMBER") != "" {
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
