// Package server exposes icon-composer over a small loopback JSON HTTP
// API for the embedded browser-based UI. Composing an icon is a pure,
// sub-millisecond string operation — no SSE/job machinery, a plain
// request/response is enough (same call this family of apps makes for
// any instant operation).
package server

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"sync/atomic"
	"time"

	"github.com/DavidMarsanic/icon-builder/web"
)

const idleTimeout = 30 * time.Minute

type Server struct {
	lastActivity atomic.Int64
}

func New() *Server {
	s := &Server{}
	s.touch()
	return s
}

func (s *Server) Start(port int) (string, error) {
	ln, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", port))
	if err != nil {
		return "", fmt.Errorf("starting local server: %w", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/icons", s.handleIcons)
	mux.HandleFunc("GET /api/icon-glyph", s.handleIconGlyph)
	mux.HandleFunc("GET /api/suggest", s.handleSuggest)
	mux.HandleFunc("POST /api/compose", s.handleCompose)
	mux.Handle("GET /", http.FileServer(http.FS(web.Static)))

	httpSrv := &http.Server{Handler: s.trackActivity(mux)}
	go func() {
		_ = httpSrv.Serve(ln)
	}()
	go s.watchIdle()

	return "http://" + ln.Addr().String(), nil
}

func (s *Server) trackActivity(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.touch()
		next.ServeHTTP(w, r)
	})
}

func (s *Server) touch() {
	s.lastActivity.Store(time.Now().Unix())
}

func (s *Server) watchIdle() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()
	for range ticker.C {
		idleFor := time.Now().Unix() - s.lastActivity.Load()
		if idleFor > int64(idleTimeout.Seconds()) {
			os.Exit(0)
		}
	}
}
