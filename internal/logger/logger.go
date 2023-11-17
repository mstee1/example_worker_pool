package logger

import (
	"github.com/mstee1/example_worker_pool/internal/config"
	"github.com/sirupsen/logrus"
)

func NewLog(
	cfgLog *config.Logger,
	workerName string,
) (*logrus.Logger, error) {

	err := checkLogDir(cfgLog.LogDir)
	if err != nil {
		return &logrus.Logger{}, err
	}

	logger, err := createLog(cfgLog, workerName)
	if err != nil {
		return &logrus.Logger{}, err
	}

	return logger, nil
}
