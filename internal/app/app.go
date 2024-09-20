package app

import (
	"github.com/AskaryanKarine/BMSTU-ds-1/internal/config"
	"github.com/AskaryanKarine/BMSTU-ds-1/internal/repositories/connection"
	"github.com/AskaryanKarine/BMSTU-ds-1/internal/repositories/person"
	"github.com/AskaryanKarine/BMSTU-ds-1/internal/server"
)

type App struct {
	srv *server.Server
	cfg config.Config
}

func New() (*App, error) {
	cfg, err := config.New()
	if err != nil {
		return nil, err
	}

	db, err := connection.OpenPostgres(cfg)
	if err != nil {
		return nil, err
	}

	personStorage := person.NewStorage(db)

	srv := server.New(personStorage)
	return &App{
		srv: srv,
		cfg: cfg,
	}, nil
}

func (a *App) Start() error {
	a.srv.Run(a.cfg.Port)
	return nil
}
