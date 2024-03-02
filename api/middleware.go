package api

import (
	"adams549659584/go-proxy-bingai/common"
	"bufio"
	"fmt"
	"net"
	"net/http"
	"strings"
)

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func newResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{w, http.StatusOK}
}

func (w *responseWriter) WriteHeader(code int) {
	w.statusCode = code
	w.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	hijacker, ok := rw.ResponseWriter.(http.Hijacker)
	if !ok {
		return nil, nil, fmt.Errorf("the ResponseWriter doesn't support the Hijacker interface")
	}
	return hijacker.Hijack()
}

func Middleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if common.SB && !strings.HasPrefix(r.URL.Path, "/web/") && r.URL.Path != "/" {
			w.WriteHeader(http.StatusUnavailableForLegalReasons)
			return
		}
		wr := newResponseWriter(w)
		next(wr, r)
		if strings.HasPrefix(r.URL.Path, "/web/") || r.URL.Path != "/" {
			common.Logger.Debug("%s - %s %s - %d - %s", strings.Split(r.RemoteAddr, ":")[0], r.Method, r.URL.Path, wr.statusCode, r.Header.Get("User-Agent"))
		} else {
			common.Logger.Info("%s - %s %s - %d - %s", strings.Split(r.RemoteAddr, ":")[0], r.Method, r.URL.Path, wr.statusCode, r.Header.Get("User-Agent"))
		}
	})
}
