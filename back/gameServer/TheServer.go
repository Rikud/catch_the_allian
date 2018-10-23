package gameServer

import (
	"IT-Berries_Go_server/controllers"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"time"
)

type TheServer struct {
	router *mux.Router
	server *http.Server
}

func (s *TheServer) bindHandlers() {
	for path, handle := range controllers.Handlers {
		s.router.Handle(path, controllers.AccessControl(handle))
	}
}

func (s *TheServer) Prepare() {
	s.router = mux.NewRouter()
	s.bindHandlers()
}

func (s *TheServer) Start() {
	s.server = &http.Server{
		Handler: s.router,
		Addr:    "0.0.0.0:8080",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Fatal(s.server.ListenAndServe())
}
