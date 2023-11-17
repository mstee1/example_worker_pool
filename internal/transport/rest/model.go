package rest

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/mstee1/example_worker_pool/internal/config"
	"github.com/sirupsen/logrus"
)

type api struct {
	ctx         context.Context
	cfg         *config.ApiControl
	log         *logrus.Logger
	errorFile   string
	countErrors int
	counter     int
	errChan     chan config.ErrCh
}

type requsetApi struct {
	Message string `json:"message"`
}

type jsonResponse struct {
	CountErrors int `json:"countErrors"`
}

func NewApi(
	ctx context.Context,
	cfg *config.ApiControl,
	cfgLog *config.Logger,
	log *logrus.Logger,
	errChan chan config.ErrCh,
) Api {

	errorFile := cfgLog.LogDir + "/example_errs.log"
	if cfgLog.LogFile != "" {
		errorFile = cfgLog.LogDir + fmt.Sprintf("/%s_errs.log", strings.Split(cfgLog.LogFile, ".")[0])
	}

	if _, err := os.Stat(errorFile); err != nil {
		_, err := os.Create(errorFile)
		if err != nil {
			log.Error(err)
		}
	}

	return &api{
		ctx:         ctx,
		cfg:         cfg,
		log:         log,
		errorFile:   errorFile,
		countErrors: 0,
		counter:     0,
		errChan:     errChan,
	}
}
