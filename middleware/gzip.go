package middleware

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"github.com/SyntaxCrew/harmony"
	"io"
	"net"
	"net/http"
	"sync"
)

const gzipScheme = "gzip"

type (
	// GzipConfig defines the config for Gzip middleware.
	GzipConfig struct {
		// Skipper defines a function to skip middleware.
		Skipper Skipper

		// Gzip compression level.
		// Optional. Default value -1.
		// -2 means use Huffman-only compression.
		// -1 means use default compression level.
		// 1 (BestSpeed) to 9 (BestCompression).
		Level int

		// MinLength is the minimum length required for compression.
		// Optional. Default value 0.
		MinLength int
	}

	gzipResponseWriter struct {
		io.Writer
		http.ResponseWriter
		code                int
		minLength           int
		isWroteHeader       bool
		isWroteBody         bool
		isMinLengthExceeded bool
		buffer              *bytes.Buffer
	}
)

// Gzip returns a middleware which compresses HTTP response using gzip compression
func Gzip(gzipCfg ...*GzipConfig) harmony.MiddlewareFunc {
	cfg := &GzipConfig{
		Skipper: defaultSkipper,
		Level:   gzip.DefaultCompression,
	}
	if len(gzipCfg) > 0 {
		cfg = gzipCfg[0]
	}

	gzipPool := sync.Pool{
		New: func() any {
			w, err := gzip.NewWriterLevel(io.Discard, cfg.Level)
			if err != nil {
				return err
			}
			return w
		},
	}
	bufferPool := sync.Pool{
		New: func() any {
			return &bytes.Buffer{}
		},
	}
	return func(next harmony.HandlerFunc) harmony.HandlerFunc {
		return func(ctx harmony.Context) error {
			if cfg.Skipper(ctx) {
				return next(ctx)
			}
			rw := ctx.ResponseWriter()
			rw.Header().Add(harmony.HeaderVary, harmony.HeaderAcceptEncoding)

			gp := gzipPool.Get()
			gw, ok := gp.(*gzip.Writer)
			if !ok {
				return harmony.NewHTTPError(http.StatusInternalServerError, gp.(error).Error())
			}
			defer func() {
				gw.Reset(rw)
				gzipPool.Put(gw)
				_ = gw.Close()
			}()

			bp := bufferPool.Get()
			buf, ok := bp.(*bytes.Buffer)
			if !ok {
				return harmony.NewHTTPError(http.StatusInternalServerError, bp.(error).Error())
			}
			defer func() {
				buf.Reset()
				bufferPool.Put(buf)
			}()

			grw := newGzipResponseWriter(rw, gw, cfg.MinLength, buf)
			defer func() {
				if grw.isWroteHeader {
					rw.WriteHeader(grw.code)
				}

				if !grw.isWroteBody {
					if rw.Header().Get(harmony.HeaderContentEncoding) == gzipScheme {
						rw.Header().Del(harmony.HeaderContentEncoding)
					}

				} else if !grw.isMinLengthExceeded {
					_, _ = grw.buffer.WriteTo(rw)
				}
			}()
			return next(ctx)
		}
	}
}

// WriteHeader implements http.ResponseWriter.
func (grw *gzipResponseWriter) WriteHeader(code int) {
	if grw.isWroteHeader {
		return
	}
	grw.Header().Del(harmony.HeaderContentLength)
	grw.isWroteHeader = true
	grw.code = code
	grw.ResponseWriter.WriteHeader(code)
}

// Write implements io.Writer.
func (grw *gzipResponseWriter) Write(b []byte) (int, error) {
	if grw.isWroteBody {
		return 0, nil
	}

	if grw.Header().Get(harmony.HeaderContentType) == "" {
		grw.Header().Set(harmony.HeaderContentType, http.DetectContentType(b))
	}
	grw.isWroteBody = true

	if !grw.isMinLengthExceeded {
		n, err := grw.buffer.Write(b)
		if err != nil {
			return 0, err
		}

		if grw.buffer.Len() >= grw.minLength {
			grw.isMinLengthExceeded = true
			grw.Header().Add(harmony.HeaderContentEncoding, gzipScheme)
			grw.WriteHeader(grw.code)
			return grw.Writer.Write(grw.buffer.Bytes())
		}
		return n, nil
	}
	return grw.Writer.Write(b)
}

// Flush implements http.Flusher.
func (grw *gzipResponseWriter) Flush() {
	if !grw.isMinLengthExceeded {
		// Enforce compression because we will not know how much more data will come
		grw.isMinLengthExceeded = true
		grw.Header().Set(harmony.HeaderContentEncoding, gzipScheme) // Issue #806
		if grw.isWroteHeader {
			grw.ResponseWriter.WriteHeader(grw.code)
		}

		_, _ = grw.Writer.Write(grw.buffer.Bytes())
	}

	_ = grw.Writer.(*gzip.Writer).Flush()
	if flusher, ok := grw.ResponseWriter.(http.Flusher); ok {
		flusher.Flush()
	}
}

// Hijack implements http.Hijacker.
func (grw *gzipResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return grw.ResponseWriter.(http.Hijacker).Hijack()
}

func newGzipResponseWriter(rw http.ResponseWriter, w io.Writer, minLength int, buffer *bytes.Buffer) *gzipResponseWriter {
	return &gzipResponseWriter{Writer: w, ResponseWriter: rw, minLength: minLength, buffer: buffer}
}
