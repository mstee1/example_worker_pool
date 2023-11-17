package core

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/mstee1/example_worker_pool/internal/config"
	"github.com/mstee1/example_worker_pool/internal/core/worker"
	"github.com/mstee1/example_worker_pool/internal/logger"
	"github.com/mstee1/example_worker_pool/internal/transport/rest"
)

func (a *app) GracefulShutdownApp(ctx context.Context, cancel context.CancelFunc) {

	select {
	case <-ctx.Done():

	case <-a.stopAppCh:
		cancel()
	}

	time.Sleep(time.Second)
	a.log.Info("Stop service\n")
	close(a.stopAppCh)
	close(a.errCh)
	close(a.wPoolCh)

}

func (a *app) Run(ctx context.Context, cond *sync.Cond, canWork *bool) {
	select {
	case <-ctx.Done():
	default:

		a.log.Info("Start service")

		api := rest.NewApi(
			ctx,
			&a.cfg.ApiControl,
			&a.cfg.Logger,
			a.log,
			a.errCh,
		)
		go api.StartApiControl()

		var workers []worker.WorkerApp
		for i, _ := range a.cfg.Workers {

			workerLog, err := logger.NewLog(
				&a.cfg.Logger,
				a.cfg.Workers[i].Name,
			)

			if err != nil {
				select {
				case <-ctx.Done():
					return
				default:
					a.errCh <- config.ErrCh{
						Name: "main",
						Log:  a.log,
						Err:  err,
					}
					continue
				}
			}

			worker := worker.NewWorkerApp(
				a.cfg.Logger,
				a.cfg.Workers[i],
				a.cfg.Sleeps,
				a.errCh,
				a.wPoolCh,
				workerLog,
			)

			workers = append(workers, worker)
			go worker.GracefulShutdownWorker(ctx)
		}

		if len(workers) < 1 {
			select {
			case <-ctx.Done():
				return
			default:
				a.errCh <- config.ErrCh{
					Name: "main",
					Log:  a.log,
					Err:  fmt.Errorf("no workers"),
				}

				select {
				case <-ctx.Done():
				case <-time.After(2 * time.Second):
					a.stopAppCh <- struct{}{}
				}

				return
			}
		}

		for _, w := range workers {

			worker := w

			go func() {

				for {
					select {
					case <-ctx.Done():
						return
					default:
						worker.Run(ctx, cond, canWork)
					}
				}
			}()
		}
	}
}
