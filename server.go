package app

import (
	"net/http"

	"github.com/sirupsen/logrus"
)

type Server struct {
	httpserver *http.Server
}

func (s *Server) Run(port string, handler http.Handler) error {
	logrus.Println("Запуск сервера... ")
	s.httpserver = &http.Server{
		Addr:    ":" + port,
		Handler: handler,
	}
	return s.httpserver.ListenAndServe()
}
