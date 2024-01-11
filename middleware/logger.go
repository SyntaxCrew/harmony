package middleware

import (
	"github.com/SyntaxCrew/harmony"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	defaultLoggerFormat = `{remote_ip} - {host} "{method} {path} {protocol}" {status} {latency}`
)

type (
	// LoggerConfig defines the config for Logger middleware.
	LoggerConfig struct {
		// Skipper defines a function to skip middleware.
		Skipper Skipper

		// Format defines the logging format with defined variables.
		// Optional. Default value {remote_ip} - {host} "{method} {path} {protocol}" {status} {latency}
		// possible variables:
		// - {remote_ip}
		// - {host}
		// - {method}
		// - {path}
		// - {protocol}
		// - {status}
		// - {latency}
		Format string
	}

	loggerResponseWriter struct {
		http.ResponseWriter
		code    int
		latency time.Duration
	}
)

// Logger returns a middleware which logs HTTP requests.
func Logger(loggerCfg ...*LoggerConfig) harmony.MiddlewareFunc {
	cfg := &LoggerConfig{
		Skipper: defaultSkipper,
		Format:  defaultLoggerFormat,
	}

	if len(loggerCfg) > 0 {
		cfg = loggerCfg[0]
	}

	return func(next harmony.HandlerFunc) harmony.HandlerFunc {
		return func(ctx harmony.Context) error {
			if cfg.Skipper(ctx) {
				return next(ctx)
			}

			req := ctx.Request()
			cfg.Format = strings.ReplaceAll(cfg.Format, "{remote_ip}", req.RemoteAddr)
			cfg.Format = strings.ReplaceAll(cfg.Format, "{host}", req.Host)
			cfg.Format = strings.ReplaceAll(cfg.Format, "{method}", req.Method)
			cfg.Format = strings.ReplaceAll(cfg.Format, "{path}", req.URL.Path)
			cfg.Format = strings.ReplaceAll(cfg.Format, "{protocol}", req.Proto)

			lrw := &loggerResponseWriter{ResponseWriter: ctx.ResponseWriter()}
			ctx.SetResponseWriter(lrw)
			_ = next(ctx)
			cfg.Format = strings.ReplaceAll(cfg.Format, "{status}", strconv.Itoa(lrw.code))
			cfg.Format = strings.ReplaceAll(cfg.Format, "{latency}", lrw.latency.String())
			log.Println(cfg.Format)
			return nil
		}
	}
}

// WriteHeader implements http.ResponseWriter.
func (w *loggerResponseWriter) WriteHeader(code int) {
	now := time.Now()
	w.code = code
	w.ResponseWriter.WriteHeader(code)
	w.latency = time.Since(now)
}
