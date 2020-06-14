package http_server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"

	"go.uber.org/zap"
)

func Test_accessLogMiddleware(t *testing.T) {
	core, o := observer.New(zapcore.InfoLevel)
	logger := zap.New(core)
	tests := map[string]struct {
		logger      *zap.Logger
		wantMessage string
	}{
		"OK": {
			logger:      logger,
			wantMessage: "access",
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			m := accessLogMiddleware(test.logger)
			h := m(http.HandlerFunc(mockHandler))

			req, err := http.NewRequest("POST", "/mock", nil)
			if err != nil {
				t.Fatal(err)
			}
			rr := httptest.NewRecorder()
			defer func() { _ = rr.Result().Body.Close() }()
			h.ServeHTTP(rr, req)

			entries := o.All()
			if len(entries) == 0 {
				t.Fatal("Log entries are empty")
			}
			if got, want := entries[0].Message, test.wantMessage; got != want {
				t.Errorf("unexpected log message: got = %v, want = %v", got, want)
			}
		})
	}
}

func mockHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
