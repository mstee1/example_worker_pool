package rest

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

func (ap *api) getStatus(w http.ResponseWriter, r *http.Request) {

	select {
	case <-ap.ctx.Done():
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusServiceUnavailable)
		message := &jsonResponse{
			CountErrors: 0,
		}
		json.NewEncoder(w).Encode(message)
		return
	default:
		var req requsetApi
		w.Header().Set("Content-Type", "application/json")
		json.NewDecoder(r.Body).Decode(&req)

		if req.Message != "get status" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)

		message := &jsonResponse{
			CountErrors: ap.countErrors,
		}

		json.NewEncoder(w).Encode(message)

		ap.countErrors = 0
		ap.counter++
		if ap.counter == 10 {
			ap.counter = 0

			if _, err := os.Stat(ap.errorFile); err == nil {
				err := os.Remove(ap.errorFile)
				if err != nil {
					ap.log.Error(err)
				}
			}

			_, err := os.Create(ap.errorFile)
			if err != nil {
				ap.log.Error(err)
			}
		}

	}
}

func (ap *api) getEror() {

	for {
		select {
		case <-ap.ctx.Done():
			return
		case message := <-ap.errChan:
			message.Log.Error(message.Err.Error())

			timeError := time.Now().Format("02-01-2006 15:04:05")

			ap.countErrors++

			errStr := fmt.Sprintf("workerName: %s, timeError: %s, body: %s", message.Name, timeError, message.Err.Error())

			ap.saveFile(errStr)

		}
	}

}

func (ap *api) saveFile(errStr string) {

	file, err := os.OpenFile(ap.errorFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		ap.log.Error(err)
		return
	}
	defer file.Close()

	_, err = file.WriteString(errStr + "\n")
	if err != nil {
		ap.log.Error(err)
		return
	}

	ap.log.Debug(fmt.Sprintf("Success write error to file %s", ap.errorFile))
}
