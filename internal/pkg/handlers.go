package pkg

import (
	"errors"
	"net/http"
)

func Return404(w http.ResponseWriter, r *http.Request) {
	var err = errors.New("route does not exists")
	HandleErrorInHTTPRequest(404, err, w)
}

type HTTPHandler func(w http.ResponseWriter, r *http.Request)

func HandleGet(handler HTTPHandler) HTTPHandler {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			handler(w, r)
			return
		} else if r.Method == http.MethodOptions {
			WriteEmptyResponse(w)
			return
		}
		Return404(w, r)
	}
}