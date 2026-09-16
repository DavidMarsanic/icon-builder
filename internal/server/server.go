// Package server exposes icon-composer over a small loopback JSON HTTP
// API for the embedded browser-based UI. Composing an icon is a pure,
// sub-millisecond string operation — no SSE/job machinery, a plain
// request/response is enough (same call this family of apps makes for
// any instant operation). Shared server plumbing (loopback bind,
// idle-timeout shutdown, job events/cancel/reveal/open routes — unused
// here beyond idle shutdown) comes from brightencode-appkit.
package server

import (
	"context"
	"net/http"

	appkit "github.com/DavidMarsanic/brightencode-appkit/server"
	"github.com/DavidMarsanic/icon-builder/web"
)

type Server struct {
	*appkit.Server
}

func New(ctx context.Context) *Server {
	return &Server{Server: appkit.New(ctx, 0)}
}

func (s *Server) Start(port int) (string, error) {
	return s.Server.Start(port, web.Static, func(mux *http.ServeMux) {
		mux.HandleFunc("GET /api/icons", s.handleIcons)
		mux.HandleFunc("GET /api/icon-glyph", s.handleIconGlyph)
		mux.HandleFunc("GET /api/suggest", s.handleSuggest)
		mux.HandleFunc("POST /api/compose", s.handleCompose)
	})
}
