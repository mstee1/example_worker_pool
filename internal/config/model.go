package config

import (
	"time"

	"github.com/sirupsen/logrus"
)

type Config struct {
	Logger     Logger     `mapstructure:"logger"`
	Sleeps     Sleeps     `mapstructure:"sleeps"`
	ApiControl ApiControl `mapstructure:"apiControl"`
	Workers    []Worker   `mapstructure:"workers"`
	Version    string
}

type Logger struct {
	LogLevel   string `mapstructure:"logLevel"`
	LogDir     string `mapstructure:"logDir"`
	LogFile    string `mapstructure:"logFile"`
	LogMode    string `mapstructure:"logMode"`
	RewriteLog bool   `mapstructure:"rewriteLog"`
}

type Sleeps struct {
	SleepService time.Duration `mapstructure:"sleepService"`
}

type ApiControl struct {
	Use     bool   `yaml:"use"`
	ApiHost string `yaml:"apiHost"`
	ApiPort string `yaml:"apiPort"`
}

type Worker struct {
	Index int
	Use   bool   `mapstructure:"use"`
	Name  string `mapstructure:"name"`
}

type ErrCh struct {
	Name string
	Log  *logrus.Logger
	Err  error
}
