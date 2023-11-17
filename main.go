package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/mstee1/example_worker_pool/internal/config"
	"github.com/mstee1/example_worker_pool/internal/core"
	"github.com/mstee1/example_worker_pool/internal/logger"
	"github.com/spf13/viper"
)

func main() {

	v := viper.New()

	cfg, err := config.GetConfig(v)
	if err != nil {
		panic(err)
	}

	if cfg.Version != "" {
		fmt.Println(cfg.Version)
		return
	}

	cfgChangeCh := make(chan struct{})
	v.OnConfigChange(func(e fsnotify.Event) {

		select {
		case cfgChangeCh <- struct{}{}:
		default:
		}

	})
	v.WatchConfig()

	ctx, cancel := context.WithCancel(context.Background())
	mu := sync.Mutex{}
	cond := sync.NewCond(&mu)
	canWork := true

	go createApp(
		ctx,
		cfg,
		v,
		cfgChangeCh,
		cond,
		&canWork,
	)

	sigCh := make(chan os.Signal)

	gracefulShutdownMain(sigCh, cancel, cfgChangeCh)
}

func gracefulShutdownMain(
	sigCh chan os.Signal,
	cancel context.CancelFunc,
	cfgChangeCh chan struct{},
) {
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM) //отслеживает команды завершения

	sign := <-sigCh
	fmt.Printf("Got signal: %v, exiting\n", sign)
	cancel()

	close(sigCh)
	close(cfgChangeCh)

	time.Sleep(time.Second * 2)

}

func createApp(
	ctx context.Context,
	cfg config.Config,
	viper *viper.Viper,
	cfgChangeCh chan struct{},
	cond *sync.Cond,
	canWork *bool,
) {

	for {

		select {
		case <-ctx.Done():
			return
		default:

			stopAppCh := make(chan struct{})
			wPoolCh := make(chan struct{}, len(cfg.Workers))

			ctxApp, cancelApp := context.WithCancel(ctx)

			logApp, err := logger.NewLog(
				&cfg.Logger,
				"main",
			)
			if err != nil {
				fmt.Println(err)
				time.Sleep(5 * time.Second)
				continue
			}

			go config.CheckChangeCfg(
				ctxApp,
				viper,
				&cfg,
				logApp,
				wPoolCh,
				stopAppCh,
				cfgChangeCh,
				cond,
				canWork,
			)

			if len(cfg.Workers) == 0 {
				logApp.Error("no workers")
				select {
				case <-ctx.Done():
				case <-stopAppCh:
				case <-time.After(cfg.Sleeps.SleepService):
					cancelApp()
				}
				continue
			}

			for i := 0; i < len(cfg.Workers); i++ {
				wPoolCh <- struct{}{}
			}

			app := core.NewApp(
				cfg,
				logApp,
				stopAppCh,
				wPoolCh,
				make(chan config.ErrCh),
			)

			go app.Run(ctxApp, cond, canWork)
			app.GracefulShutdownApp(ctxApp, cancelApp)

		}
	}
}
