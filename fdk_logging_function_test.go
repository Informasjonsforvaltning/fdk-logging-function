package fdk_logging_function

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandlers(t *testing.T) {
	const validError = `{"message":"m","severity":"ERROR","application":"a","image":"i"}`
	tooLarge := `{"message":"` + strings.Repeat("x", maxBodyBytes) + `","severity":"ERROR","application":"a","image":"i"}`

	tests := []struct {
		name    string
		handler http.HandlerFunc
		method  string
		body    string
		want    int
	}{
		{"error preflight", errorLogging, http.MethodOptions, "", http.StatusNoContent},
		{"error wrong method", errorLogging, http.MethodGet, "", http.StatusMethodNotAllowed},
		{"error valid", errorLogging, http.MethodPost, validError, http.StatusNoContent},
		{"error missing field", errorLogging, http.MethodPost, `{"message":"m"}`, http.StatusBadRequest},
		{"error bad json", errorLogging, http.MethodPost, `{`, http.StatusBadRequest},
		{"error body too large", errorLogging, http.MethodPost, tooLarge, http.StatusBadRequest},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, "/", strings.NewReader(tt.body))
			rec := httptest.NewRecorder()

			tt.handler(rec, req)

			if rec.Code != tt.want {
				t.Errorf("status = %d, want %d", rec.Code, tt.want)
			}
			if got := rec.Header().Get("Access-Control-Allow-Origin"); got != "*" {
				t.Errorf("Access-Control-Allow-Origin = %q, want *", got)
			}
		})
	}
}
