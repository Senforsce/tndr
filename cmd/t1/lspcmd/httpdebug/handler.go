package httpdebug

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/url"

	"github.com/senforsce/tndr"
	"github.com/senforsce/tndr/cmd/t1/lspcmd/proxy"
	"github.com/senforsce/tndr/cmd/t1/visualize"
)

var log *slog.Logger

func NewHandler(l *slog.Logger, s *proxy.Server) http.Handler {
	m := http.NewServeMux()
	log = l
	m.HandleFunc("/tndr", func(w http.ResponseWriter, r *http.Request) {
		uri := r.URL.Query().Get("uri")
		c, ok := s.TndrSource.Get(uri)
		if !ok {
			Error(w, "uri not found", http.StatusNotFound)
			return
		}
		String(w, c.String())
	})
	m.HandleFunc("/sourcemap", func(w http.ResponseWriter, r *http.Request) {
		uri := r.URL.Query().Get("uri")
		sm, ok := s.SourceMapCache.Get(uri)
		if !ok {
			Error(w, "uri not found", http.StatusNotFound)
			return
		}
		JSON(w, sm.SourceLinesToTarget)
	})
	m.HandleFunc("/go", func(w http.ResponseWriter, r *http.Request) {
		uri := r.URL.Query().Get("uri")
		c, ok := s.GoSource[uri]
		if !ok {
			Error(w, "uri not found", http.StatusNotFound)
			return
		}
		String(w, c)
	})
	m.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		uri := r.URL.Query().Get("uri")
		if uri == "" {
			// List all URIs.
			if err := list(s.TndrSource.URIs()).Render(r.Context(), w); err != nil {
				Error(w, "failed to list URIs", http.StatusInternalServerError)
			}
			return
		}
		// Assume we've got a URI.
		tndrSource, ok := s.TndrSource.Get(uri)
		if !ok {
			if !ok {
				Error(w, "uri not found in document contents", http.StatusNotFound)
				return
			}
		}
		goSource, ok := s.GoSource[uri]
		if !ok {
			if !ok {
				Error(w, "uri not found in document contents", http.StatusNotFound)
				return
			}
		}
		sm, ok := s.SourceMapCache.Get(uri)
		if !ok {
			Error(w, "uri not found", http.StatusNotFound)
			return
		}
		if err := visualize.HTML(uri, tndrSource.String(), goSource, sm).Render(r.Context(), w); err != nil {
			Error(w, "failed to visualize HTML", http.StatusInternalServerError)
		}
	})
	return m
}

func getMapURL(uri string) tndr.SafeURL {
	return withQuery("/", uri)
}

func getSourceMapURL(uri string) tndr.SafeURL {
	return withQuery("/sourcemap", uri)
}

func getTndrURL(uri string) tndr.SafeURL {
	return withQuery("/tndr", uri)
}

func getGoURL(uri string) tndr.SafeURL {
	return withQuery("/go", uri)
}

func withQuery(path, uri string) tndr.SafeURL {
	q := make(url.Values)
	q.Set("uri", uri)
	u := &url.URL{
		Path:     path,
		RawPath:  path,
		RawQuery: q.Encode(),
	}
	return tndr.SafeURL(u.String())
}

func JSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	if err := enc.Encode(v); err != nil {
		log.Error("failed to write JSON response", slog.Any("error", err))
	}
}

func String(w http.ResponseWriter, s string) {
	if _, err := io.WriteString(w, s); err != nil {
		log.Error("failed to write string response", slog.Any("error", err))
	}
}

func Error(w http.ResponseWriter, msg string, status int) {
	w.WriteHeader(status)
	if _, err := io.WriteString(w, msg); err != nil {
		log.Error("failed to write error response", slog.Any("error", err))
	}
}
