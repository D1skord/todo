package core_http_middleware

import (
	"context"
	"net/http"
	"time"

	core_logger "github.com/D1skord/todo/internal/core/logger"
	core_http_response "github.com/D1skord/todo/internal/core/transport/http/response"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

const (
	requestIDHeader = "X-Request-Id"
)

func RequestId() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := r.Header.Get(requestIDHeader)
			if requestID == "" {
				requestID = uuid.NewString()
			}

			r.Header.Set(requestIDHeader, requestID)
			r.Header.Set(requestIDHeader, requestID)

			next.ServeHTTP(w, r)
		})
	}
}

func Logger(log *core_logger.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := r.Header.Get(requestIDHeader)

			l := log.With(
				zap.String("request_id", requestID),
				zap.String("url", r.URL.String()),
				zap.String("method", r.Method),
			)

			ctx := context.WithValue(r.Context(), "log", l)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func Panic() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			logger := core_logger.FromContext(ctx)
			responseHandler := core_http_response.NewHTTPResponseHandler(logger, w)

			defer func() {
				if err := recover(); err != nil {
					responseHandler.PanicResponse(err, "during handle http request got unexpected panic")
				}
			}()

			next.ServeHTTP(w, r)

		})
	}
}

func Trace(log *core_logger.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger := core_logger.FromContext(r.Context())
			rw := core_http_response.NewResponseWriter(w)

			before := time.Now()

			logger.Debug(
				">>> incoming HTTP request",
				zap.String("method", r.Method),
				zap.Time("time", before.UTC()),
			)

			next.ServeHTTP(rw, r)

			logger.Debug(
				">>> done HTTP request",
				zap.Duration("latency", time.Since(before)),
				zap.Int("status_code", rw.GetStatusCodeOrPanic()),
			)
		})
	}
}
