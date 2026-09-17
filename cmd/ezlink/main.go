package main

import (
	"github.com/p1shiA/ezLink/config"
	"github.com/p1shiA/ezLink/internal/utils"
)

func main() {
	// loading env
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	// start logging
	log, err := utils.Init(cfg)
	if err != nil {
		panic(err)
	}

	defer log.Sync()

	log.Info("Starting application")


}