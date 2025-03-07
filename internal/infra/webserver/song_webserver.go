package webserver

import (
	"encoding/json"
	"net/http"
	"strconv"
)

type SongWebServer struct{}

func NewSongWebServer() *SongWebServer {
	return &SongWebServer{}
}

func (s *SongWebServer) GetSongsHandler(w http.ResponseWriter, r *http.Request) {

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
