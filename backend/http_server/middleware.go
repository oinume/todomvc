package http_server

import (
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

func accessLogMiddleware(logger *zap.Logger) mux.MiddlewareFunc {
	return func(h http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			writerProxy := WrapWriter(w)
			h.ServeHTTP(writerProxy, r)
			func() {
				end := time.Now()
				status := writerProxy.Status()
				if status == 0 {
					status = http.StatusOK
				}
				//trackingID := ""
				//if v, err := context_data.GetTrackingID(r.Context()); err == nil {
				//	trackingID = v
				//}

				// 180.76.15.26 - - [31/Jul/2016:13:18:07 +0000] "GET / HTTP/1.1" 200 612 "-" "Mozilla/5.0 (compatible; Baiduspider/2.0; +http://www.baidu.com/search/spider.html)"
				logger.Info(
					"access",
					zap.String("date", start.Format(time.RFC3339)),
					zap.String("method", r.Method),
					zap.String("url", r.URL.String()),
					zap.Int("status", status),
					zap.Int("bytes", writerProxy.BytesWritten()),
					zap.String("remoteAddr", getRemoteAddress(r)),
					zap.String("userAgent", r.Header.Get("User-Agent")),
					zap.String("referer", r.Referer()),
					zap.Duration("elapsedMsec", end.Sub(start)/time.Millisecond),
					//zap.String("trackingID", trackingID),
				)
			}()
		}
		return http.HandlerFunc(fn)
	}
}

func getRemoteAddress(req *http.Request) string {
	xForwardedFor := req.Header.Get("X-Forwarded-For")
	if xForwardedFor == "" {
		return (strings.Split(req.RemoteAddr, ":"))[0]
	}
	return strings.TrimSpace((strings.Split(xForwardedFor, ","))[0])
}
