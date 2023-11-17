package config

import (
	"context"
	"flag"
	"os"
	"strings"
	"sync"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Функция получения заполненной структуры Config
func GetConfig(v *viper.Viper) (Config, error) {
	var config Config
	var absoluteConfigPath, configPath, configName, configType string
	args := os.Args
	for _, arg := range args {

		// Если будет запрошена версия, то она будет выведена
		if arg == "-version" {
			config.Version = "1.0.0"
			return config, nil
		}

		// Если будет задан не стандартный путь к конфигу, то он будет использован
		if strings.Split(arg, "=")[0] == "-configPath" {
			absoluteConfigPath = strings.Split(arg, "=")[1]
			configPath, configName, configType = getParamsConf(absoluteConfigPath)
			break
		}
	}

	// Если не будет задан путь к конфигу, то он будет использован по умолчанию
	if absoluteConfigPath == "" {
		configPath = "./"
		configName = "config"
		configType = "yaml"
	}

	// Попытка чтения конфига
	err := readParametersFromConfig(v, configPath, configName, configType, &config)
	if err != nil {
		return Config{}, err
	}

	// Изменение полей конфига при наличии флагов
	readFlags(&config)

	// Проверка воркеров
	checkWorkers(&config)

	return config, nil
}

func getParamsConf(absoluteConfigPath string) (string, string, string) {

	pathSplit := strings.Split(absoluteConfigPath, "/")
	configNameType := strings.Split(pathSplit[len(pathSplit)-1], ".")
	configPath := strings.Join(pathSplit[:len(pathSplit)-1], "/") + "/"

	return configPath, configNameType[0], configNameType[1]

}

func readParametersFromConfig(
	viper *viper.Viper,
	configPath,
	configName,
	configType string,
	cfg *Config,
) error {

	viper.SetConfigName(configName)
	viper.SetConfigType(configType)
	viper.AddConfigPath(configPath)

	// Попытка чтения конфига
	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	// Попытка заполнение структуры Config полученными данными
	if err := viper.Unmarshal(&cfg); err != nil {
		return err
	}
	return nil
}

func readFlags(cfg *Config) {
	var configPath string

	flag.StringVar(&cfg.Logger.LogLevel, "logLevel", cfg.Logger.LogLevel, "Aviable LogLevael: INFO,DEBUG,TEST")
	flag.StringVar(&cfg.Logger.LogDir, "logDir", cfg.Logger.LogDir, "Full path to save log file")
	flag.StringVar(&cfg.Logger.LogFile, "logFile", cfg.Logger.LogFile, "Name of log file")
	flag.StringVar(&cfg.Logger.LogMode, "logMode", cfg.Logger.LogMode, "Aviable LogMode:stdout,file,empty field")
	flag.BoolVar(&cfg.Logger.RewriteLog, "rewriteLog", cfg.Logger.RewriteLog, "Overwriting a log file")

	flag.DurationVar(&cfg.Sleeps.SleepService, "sleepService", cfg.Sleeps.SleepService, "Sleep time for service")

	flag.BoolVar(&cfg.ApiControl.Use, "apiUse", cfg.ApiControl.Use, "Use api control")
	flag.StringVar(&cfg.ApiControl.ApiHost, "apiHost", cfg.ApiControl.ApiHost, "Host of api control")
	flag.StringVar(&cfg.ApiControl.ApiPort, "apiPort", cfg.ApiControl.ApiPort, "Port of api control")

	flag.StringVar(&configPath, "configPath", "./", "The configPath parameter")
	flag.Parse()
}

func checkWorkers(cfg *Config) {
	var workers []Worker
	if len(cfg.Workers) > 0 {
		for _, v := range cfg.Workers {
			if v.Use {
				workers = append(workers, v)
			}
		}
		for i := 0; i < len(workers); i++ {
			workers[i].Index = i
		}
	}

	cfg.Workers = workers
}

// Функция для отслеживания изменений в конфиг файле
func CheckChangeCfg(
	ctx context.Context,
	viper *viper.Viper,
	cfg *Config,
	log *logrus.Logger,
	wPoolCh,
	stopAppCh,
	cfgChangeCh chan struct{},
	cond *sync.Cond,
	canWork *bool,
) {
	select {
	case <-ctx.Done():

		// Если конфиг изменен
	case <-cfgChangeCh:

		log.Info("Config file was changed")

		// Блокировка старта новых воркеров
		cond.L.Lock()
		*canWork = false
		// Получение всех работ воркеров
		for i := 0; i < len(cfg.Workers); i++ {
			select {
			case <-ctx.Done():
				return
			case <-wPoolCh:

			}

		}
		cfg.Workers = nil

		// Перечитывание и заполнение конфига
		if err := viper.Unmarshal(cfg); err != nil {
			log.Error(err)
		}
		flag.Parse()
		cfg.Logger.RewriteLog = false

		// Вызов Gracefull Shutdown для прототипа приложения,
		//чтобы оно перезапустилось
		select {
		case <-ctx.Done():
		default:
			stopAppCh <- struct{}{}
		}

		checkWorkers(cfg)

		// Разблокировка старта новых воркеров
		*canWork = true
		cond.L.Unlock()

		defer cond.Broadcast()

	}

}
