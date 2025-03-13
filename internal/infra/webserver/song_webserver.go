package webserver

import (
	"encoding/json"
	"net/http"
	"strconv"

	middlewares_ratelimiter "github.com/renebizelli/ratelimiter/internal/infra/middlewares/ratelimiter"
)

type SongWebServer struct {
	mux         *http.ServeMux
	rateLimiter *middlewares_ratelimiter.RateLimiterMiddleware
}

func NewSongWebServer(mux *http.ServeMux, rateLimiter *middlewares_ratelimiter.RateLimiterMiddleware) *SongWebServer {
	return &SongWebServer{mux: mux, rateLimiter: rateLimiter}
}

func (l *SongWebServer) RegisterRoutes() {
	l.mux.Handle("GET /songs", l.rateLimiter.Limiter(http.HandlerFunc(l.getSongsHandler)))
}

func (s *SongWebServer) getSongsHandler(w http.ResponseWriter, r *http.Request) {

	type song struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}

	songs := []song{}

	for i := 0; i < 100; i++ {
		songs = append(songs, song{ID: i + 1, Name: "Song " + strconv.Itoa(i+1)})
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(songs)
}
