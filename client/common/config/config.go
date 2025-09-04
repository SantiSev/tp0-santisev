package config

import (
	"os"
	"strings"

	"github.com/7574-sistemas-distribuidos/docker-compose-init/client/common/client"
	"github.com/op/go-logging"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

const CONFIG_FILE_PATH = "./config.yaml"

func InitConfig() (*client.ClientConfig, error) {

	v := viper.New()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	configPath := CONFIG_FILE_PATH

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
