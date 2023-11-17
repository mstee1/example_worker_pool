package worker

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

func (w *workerApp) GracefulShutdownWorker(ctx context.Context) {
	<-ctx.Done()
	w.log.Info(fmt.Sprintf("Worker %s stopped\n", w.cfgw.worker.Name))
}

func (w *workerApp) Run(ctx context.Context, cond *sync.Cond, canWork *bool) {

	cond.L.Lock()
	if !*canWork {
		cond.Wait()
	}
	cond.L.Unlock()

	select {
	case <-ctx.Done():
		return
	case <-w.chans.wPoolCh:

		w.log.Info(fmt.Sprintf("The new cycle is start - %s", w.cfgw.worker.Name))

		startCycle := time.Now()

		defer func() {
			sleepTime := w.cfgw.sleeps.SleepService - time.Since(startCycle)
			select {
			case <-ctx.Done():
				return
			default:
				w.chans.wPoolCh <- struct{}{}
				h := int(sleepTime.Hours())
				m := int(sleepTime.Minutes()) % 60
				s := int(sleepTime.Seconds()) % 60
				w.log.Info(fmt.Sprintf("Cycle %s completed, start sleep -"+
					" %d Hours, %d minutes, %d seconds",
					w.cfgw.worker.Name, h, m, s))

				select {
				case <-ctx.Done():
				case <-time.After(sleepTime):
				}
			}
		}()

		w.sendErr(
			ctx,
			errors.New("Test error"),
		)

	}
}
