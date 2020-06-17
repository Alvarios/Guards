package server

import (
	"github.com/Alvarios/guards/guards"
	"github.com/gorilla/mux"
	"net/http"
	"time"
)

type Server struct {
	Log    *guards.Guards
	Router *mux.Router
}

func NewServer(log *guards.Guards, r *mux.Router) Server {
	Server := Server{
		Log:    log,
		Router: r,
	}
	return Server
}

func (s *Server) Hello(w http.ResponseWriter, r *http.Request) {
	s.Log.InvalidRequest(r, nil, "created")
}

func (s *Server) Run() {

	s.Router.Handle("/hello", s.Log.C.ThenFunc(s.Hello)).Methods("GET")
	server := http.Server{
		Handler:      s.Router,
		Addr:         ":8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	server.ListenAndServe()
}
