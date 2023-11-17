package rest

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func (ap *api) StartApiControl() {
	if ap.cfg.Use {
		go ap.getEror()

		adrApi := ap.cfg.ApiHost + ":" + ap.cfg.ApiPort
		r := mux.NewRouter()
		r.HandleFunc("/api/status/", ap.getStatus).Methods("POST")
		ap.log.Info(fmt.Sprintf("Start work api at http://%s:%s/api/status/",
			ap.cfg.ApiHost, ap.cfg.ApiPort))

		server := &http.Server{
			Addr:    adrApi,
			Handler: r,
		}
		go func() {
			<-ap.ctx.Done()
			server.Shutdown(ap.ctx)
			server = nil
		}()
		server.ListenAndServe()
		return
	}
	ap.getEror()
}
