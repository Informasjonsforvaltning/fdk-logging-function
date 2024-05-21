package fdk_logging_function

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
)

func init() {
    functions.HTTP("CspLogging", cspLogging)
}

type CspReportLogEntry struct {
	CspReport CspReport `json:"csp-report"`
}

type CspReport struct {
	BlockedUri          string `json:"blocked-uri"`
	Disposition         string `json:"disposition"`
	DocumentUri         string `json:"document-uri"`
	EffectiveDirective  string `json:"effective-directive"`
	OriginalPolicy      string `json:"original-policy"`
	Referrer            string `json:"referrer"`
	ScriptSample        string `json:"script-sample"`
	StatusCode          int `json:"status-code"`
	ViolatedDirective   string `json:"violated-directive"`
}

func (e CspReportLogEntry) String() string {
	out, err := json.Marshal(e)
	if err != nil {
		log.Printf("json.Marshal: %v", err)
	}
	return string(out)
}

func cspLogging(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	switch r.Method {
	case http.MethodPost:
		cspLogger(w, r)
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

func cspLogger(w http.ResponseWriter, r *http.Request) {
	var cspReportLogEntry CspReportLogEntry

	if err := json.NewDecoder(r.Body).Decode(&cspReportLogEntry); err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Failed to parse csp report.")
		return
	}

	fmt.Println(cspReportLogEntry)
	fmt.Fprintf(w, "Logged csp report")
}
