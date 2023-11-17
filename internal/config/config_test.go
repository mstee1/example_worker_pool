package config

import (
	"os"
	"testing"
	"time"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
)

func TestVersionFlag(t *testing.T) {

	v := viper.New()

	os.Args = []string{"-version"}

	cfg, err := GetConfig(v)
	require.NoError(t, err)

	require.NotNil(t, cfg.Version)
}

func TestGetConfig(t *testing.T) {

	t.Run("Test flags", func(t *testing.T) {

		v := viper.New()

		os.Args = []string{
			"-configPath=../../config.yaml",
			"-logLevel=INFO",
			"-logDir=./test",
			"-logFile=test.log",
			"-logMode=file",
			"-rewriteLog=false",
			"-sleepService=5s",
			"-apiUse=true",
			"-apiHost=127.0.0.1",
			"-apiPort=8080",
		}

		cfg, err := GetConfig(v)
		require.NoError(t, err)

		require.Equal(t, "INFO", cfg.Logger.LogLevel)
		require.Equal(t, "./test", cfg.Logger.LogDir)
		require.Equal(t, "test.log", cfg.Logger.LogFile)
		require.Equal(t, "file", cfg.Logger.LogMode)
		require.Equal(t, false, cfg.Logger.RewriteLog)
		require.Equal(t, 5*time.Second, cfg.Sleeps.SleepService)
		require.Equal(t, true, cfg.ApiControl.Use)
		require.Equal(t, "127.0.0.1", cfg.ApiControl.ApiHost)
		require.Equal(t, "8080", cfg.ApiControl.ApiPort)

	})

	t.Run("Test empty config path", func(t *testing.T) {
		v := viper.New()

		os.Args = []string{}

		cfg, err := GetConfig(v)

		require.Error(t, err)
		require.Equal(t, Config{}, cfg)
	})

}

func TestGetParams(t *testing.T) {

	cPath, cName, cType := getParamsConf("/opt/test/configs/service/conf.yaml")

	require.Equal(t, "/opt/test/configs/service/", cPath)
	require.Equal(t, "conf", cName)
	require.Equal(t, "yaml", cType)
}

func TestCheckWorkers(t *testing.T) {

	c := Config{
		Workers: nil,
	}
	cOld := c

	checkWorkers(&c)
	require.Equal(t, c, cOld)

	c = Config{
		Workers: []Worker{
			{
				Index: 0,
				Use:   true,
			},
		},
	}
	cOld = c
	checkWorkers(&c)
	require.Equal(t, c, cOld)

	c = Config{
		Workers: []Worker{
			{
				Index: 0,
				Use:   false,
			},
			{
				Index: 1,
				Use:   true,
			},
		},
	}
	cOld = c
	checkWorkers(&c)
	require.NotEqual(t, c, cOld)

}

func TestReadParameters(t *testing.T) {

	t.Run("Read real file", func(t *testing.T) {

		v := viper.New()
		cPath, cName, cType := getParamsConf("../../config.yaml")
		cfg := Config{}
		cfg2 := cfg

		err := readParametersFromConfig(v, cPath, cName, cType, &cfg)

		require.NoError(t, err)
		require.NotEqual(t, cfg, cfg2)

	})

	t.Run("Read not exist file", func(t *testing.T) {
		path := "../../bad.test"
		cPath, cName, cType := getParamsConf(path)
		cfg := Config{}
		cfg2 := cfg
		v := viper.New()
		err := readParametersFromConfig(v, cPath, cName, cType, &cfg)

		require.Error(t, err)
		require.Equal(t, cfg, cfg2)
	})

	t.Run("Read empty file", func(t *testing.T) {
		path := "../../test.yaml"
		f, _ := os.Create(path)
		f.WriteString("logger: nil")
		f.Close()
		defer os.Remove(path)
		cPath, cName, cType := getParamsConf(path)
		cfg := Config{}
		cfg2 := cfg
		v := viper.New()
		err := readParametersFromConfig(v, cPath, cName, cType, &cfg)

		require.Error(t, err)
		require.Equal(t, cfg, cfg2)
	})

}
