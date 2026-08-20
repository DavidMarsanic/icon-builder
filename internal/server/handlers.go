package server

import (
	"encoding/json"
	"net/http"

	"github.com/DavidMarsanic/icon-composer/compose"
)

func (s *Server) handleIcons(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{"icons": compose.IconNames()})
}

func (s *Server) handleIconGlyph(w http.ResponseWriter, r *http.Request) {
	svg, err := compose.RawSVG(r.URL.Query().Get("name"))
	if err != nil {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", "image/svg+xml")
	_, _ = w.Write([]byte(svg))
}

func (s *Server) handleSuggest(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	writeJSON(w, http.StatusOK, map[string]any{"letters": compose.Initials(name, 3)})
}

type composeRequest struct {
	Icon    string `json:"icon"`
	Letters string `json:"letters"`
	Name    string `json:"name"`
	From    string `json:"from"`
	To      string `json:"to"`
}

func (s *Server) handleCompose(w http.ResponseWriter, r *http.Request) {
	var req composeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "invalid request body"})
		return
	}

	svg, err := compose.Compose(compose.Spec{
		Icon:    req.Icon,
		Letters: req.Letters,
		From:    req.From,
		To:      req.To,
		Seed:    req.Name,
	})
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{"svg": svg})
}

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}
