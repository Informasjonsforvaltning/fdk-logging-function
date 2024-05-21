package fdk_logging_function

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
)

func init() {
    functions.HTTP("ErrorLogging", errorLogging)
}

type ErrorLogEntry struct {
	Message     string `json:"message"`
	Severity    string `json:"severity"`
	Namespace   string `json:"namespace"`
	Trace       string `json:"trace"`
	Name        string `json:"name"`
	Location    string `json:"location"`
	Application string `json:"application"`
	Image       string `json:"image"`
}

func (e ErrorLogEntry) String() string {
	out, err := json.Marshal(e)
	if err != nil {
		log.Printf("json.Marshal: %v", err)
	}
	return string(out)
}

func isInvalidErrorLogEntry(e ErrorLogEntry) bool {
	return e.Message == "" || e.Severity == "" || e.Application == "" || e.Image == ""
}

func errorLogging(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	switch r.Method {
	case http.MethodPost:
		errorLogger(w, r)
	case http.MethodOptions:
		w.Header().Set("Access-Control-Allow-Methods", "POST")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Max-Age", "3600")
		w.WriteHeader(http.StatusNoContent)
		return
	default:
		http.Error(w, "405 - Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func errorLogger(w http.ResponseWriter, r *http.Request) {
	var errorLogEntry ErrorLogEntry

	if err := json.NewDecoder(r.Body).Decode(&errorLogEntry); err != nil || isInvalidErrorLogEntry(errorLogEntry) {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Failed to parse error log message.")
		return
	}

	fmt.Println(errorLogEntry)
	fmt.Fprintf(w, "Logged error message: %s", errorLogEntry)
}
