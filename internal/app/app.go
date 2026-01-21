package app

import (
	"log"
	"net/http"
)

type App struct {
	ServiceProvider *ServiceProvider
}

func NewApp() *App {
	return &App{}
}

func (s *App) initServiceProvider() {
	s.ServiceProvider = newServiceProvider()
}

func (s *App) Run() error {
	s.initServiceProvider()

	r := s.ServiceProvider.Router()

	log.Printf("starting server at %s", s.ServiceProvider.HTTPCfg().Address())
	err := http.ListenAndServe(s.ServiceProvider.HTTPCfg().Address(), r)
	if err != nil {
		return err
	}
	return err
}
