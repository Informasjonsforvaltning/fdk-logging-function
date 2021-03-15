package fdk_logging_function

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type LogEntry struct {
	Message     string `json:"message"`
	Severity    string `json:"severity"`
	Trace       string `json:"trace"`
	Environment string `json:"environment"`
	Name        string `json:"name"`
}

func (e LogEntry) String() string {
	out, err := json.Marshal(e)
	if err != nil {
		log.Printf("json.Marshal: %v", err)
	}
	return string(out)
}

func isInvalid(e LogEntry) bool {
	return e.Message == "" || e.Severity == ""
}

func Logging(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		logger(w, r)
	case http.MethodOptions:
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Max-Age", "3600")
		w.WriteHeader(http.StatusNoContent)
		return
	default:
		http.Error(w, "405 - Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func logger(w http.ResponseWriter, r *http.Request) {
	var logEntry LogEntry

	if err := json.NewDecoder(r.Body).Decode(&logEntry); err != nil || isInvalid(logEntry) {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Failed to parse log message.")
		return
	}

	fmt.Println(logEntry)
	fmt.Fprintf(w, "Logged message: %s", logEntry)
}
