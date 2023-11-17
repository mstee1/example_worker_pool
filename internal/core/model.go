package core

import (
	"github.com/mstee1/example_worker_pool/internal/config"
	"github.com/sirupsen/logrus"
)

type app struct {
	cfg       config.Config
	log       *logrus.Logger
	stopAppCh chan struct{}
	wPoolCh   chan struct{}
	errCh     chan config.ErrCh
}

func NewApp(
	cfg config.Config,
	log *logrus.Logger,
	stopAppCh chan struct{},
	wPoolCh chan struct{},
	errCh chan config.ErrCh,
) App {

	return &app{
		cfg:       cfg,
		log:       log,
		stopAppCh: stopAppCh,
		wPoolCh:   wPoolCh,
		errCh:     errCh,
	}

}
