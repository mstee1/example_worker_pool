package logger

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/mstee1/example_worker_pool/internal/config"
	"github.com/sirupsen/logrus"
	easy "github.com/t-tomalak/logrus-easy-formatter"
)

func initLogLevel(level string) logrus.Level {
	switch level {
	case "FATAL":
		return logrus.FatalLevel
	case "ERROR":
		return logrus.ErrorLevel
	case "WARN":
		return logrus.WarnLevel
	case "INFO":
		return logrus.InfoLevel
	case "DEBUG":
		return logrus.DebugLevel
	case "TRACE":
		return logrus.TraceLevel
	case "TEST":
		return logrus.DebugLevel
	default:
		return logrus.InfoLevel
	}
}

func checkLogDir(logDir string) error {

	_, err := os.Stat(logDir)
	if err != nil {
		err := os.MkdirAll(logDir, 0755)
		if err != nil {
			return err
		}
	}
	return nil
}

func createLog(
	cfgLog *config.Logger,
	workerName string,
) (*logrus.Logger, error) {

	if cfgLog.LogMode == "stdout" {
		log := &logrus.Logger{
			Out:   io.MultiWriter(os.Stdout),
			Level: initLogLevel(cfgLog.LogLevel),
			Formatter: &easy.Formatter{
				TimestampFormat: "2006-01-02 15:04:05",
				LogFormat:       "[%lvl%]: (%time%) %msg%\n",
			},
		}
		return log, nil
	}

	logFile := cfgLog.LogDir + "/example.log"
	if workerName == "main" {
		if cfgLog.LogFile != "" {
			logFile = cfgLog.LogDir + fmt.Sprintf("/%s.log",
				strings.Split(cfgLog.LogFile, ".")[0])
		}
	} else {
		logFile = cfgLog.LogDir + fmt.Sprintf("/example_%s.log", workerName)
		if cfgLog.LogFile != "" {
			logFile = cfgLog.LogDir + fmt.Sprintf("/%s_%s.log",
				strings.Split(cfgLog.LogFile, ".")[0], workerName)
		}
	}

	if cfgLog.RewriteLog {
		if _, err := os.Stat(logFile); err == nil {
			err := os.Remove(logFile)
			if err != nil {
				return &logrus.Logger{}, err
			}

		}
	}
	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		logrus.Fatalf("cannot open file for logging: %v - %s", err, logFile)
		return &logrus.Logger{}, err
	}

	if cfgLog.LogMode == "file" {
		log := &logrus.Logger{
			Out:   io.MultiWriter(file),
			Level: initLogLevel(cfgLog.LogLevel),
			Formatter: &easy.Formatter{
				TimestampFormat: "2006-01-02 15:04:05",
				LogFormat:       "[%lvl%]: (%time%) %msg%\n",
			},
		}
		return log, nil
	}

	log := &logrus.Logger{
		Out:   io.MultiWriter(file, os.Stdout),
		Level: initLogLevel(cfgLog.LogLevel),
		Formatter: &easy.Formatter{
			TimestampFormat: "2006-01-02 15:04:05",
			LogFormat:       "[%lvl%]: (%time%) %msg%\n",
		},
	}

	return log, nil
}
