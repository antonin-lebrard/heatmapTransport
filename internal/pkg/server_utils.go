package pkg

import (
	"encoding/json"
	"log"
	"net/http"
)

var notMarshableErrorResponse = []byte(`cannot marshal json error, the error will appear in the logs of heatmap server`)

func formatError(errToSend error) []byte {
	var errResp = errToSend.Error()

	rawResp, err := json.Marshal(errResp)
	if err != nil {
		log.Println("not marshable error: json.Marshal error:", err.Error(), "the error to send:", err.Error())
		return notMarshableErrorResponse
	}
	return rawResp
}

func HandleErrorInHTTPRequest(statusIfError int, err error, w http.ResponseWriter) bool {
	if err == nil {
		return false
	}
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Methods", "GET,PUT,POST,DELETE,OPTIONS")
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type, Authorization, Content-Length, X-Requested-With, X-Access-Token")
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(statusIfError)
	_, err = w.Write(formatError(err))
	if err != nil {
		log.Println("error when writing error to http response:", err.Error())
	}
	return true
}

func WriteResponse(data interface{}, w http.ResponseWriter) {
	response, err := json.Marshal(data)
	if HandleErrorInHTTPRequest(500, err, w) { return }

	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Methods", "GET,PUT,POST,DELETE,OPTIONS")
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type, Authorization, Content-Length, X-Requested-With, X-Access-Token")
	w.Header().Add("Content-Type", "application/json")
	_, err = w.Write(response)
	if err != nil {
		log.Println("error when writing response to http response:", err.Error())
	}
}

func WriteEmptyResponse(w http.ResponseWriter) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Methods", "GET,PUT,POST,DELETE,OPTIONS")
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type, Authorization, Content-Length, X-Requested-With, X-Access-Token")
	w.Header().Add("Content-Type", "application/json")
	_, err := w.Write([]byte{})
	if err != nil {
		log.Println("error when writing response to http response:", err.Error())
	}
}