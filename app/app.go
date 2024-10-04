package app

import (
	"log"

	"go-flattern-prices/internal/configs"
	"go-flattern-prices/internal/controller"
	"go-flattern-prices/internal/logger"
	"go-flattern-prices/internal/store"
)

type app struct {
	controller controller.Controller
}

type App interface {
	Start()
}

func New() (App, error) {
	a := &app{}

	cfg, err := configs.MustLoad()
	if err != nil {
		return nil, err
	}

	logger.Init(cfg.LogLevel, cfg.LogLevelFile, cfg.LogsDir)

	s, err := store.New(cfg.DBFile)
	if err != nil {
		return nil, err
	}

	a.controller = controller.New(cfg, s)
	return a, nil
}

func (a *app) Start() {
	err := a.controller.DoJob()
	if err != nil {
		log.Fatal(err)
	}
}
