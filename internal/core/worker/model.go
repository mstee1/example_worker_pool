package worker

import (
	"github.com/mstee1/example_worker_pool/internal/config"
	"github.com/sirupsen/logrus"
)

type workerApp struct {
	cfgw  cfgWorker
	chans chans
	log   *logrus.Logger
}

type cfgWorker struct {
	logger config.Logger
	worker config.Worker
	sleeps config.Sleeps
}

type chans struct {
	errCh   chan config.ErrCh
	wPoolCh chan struct{}
}

func NewWorkerApp(
	cfgLog config.Logger,
	cfgWork config.Worker,
	sleeps config.Sleeps,
	errCh chan config.ErrCh,
	wPoolCh chan struct{},
	log *logrus.Logger,
) WorkerApp {

	return &workerApp{
		cfgw: cfgWorker{
			logger: cfgLog,
			worker: cfgWork,
			sleeps: sleeps,
		},

		chans: chans{
			errCh,
			wPoolCh,
		},

		log: log,
	}
}
