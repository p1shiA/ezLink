package main

import (
	"github.com/p1shiA/ezLink/config"
	"github.com/p1shiA/ezLink/internal/bot/handlers"
	"github.com/p1shiA/ezLink/internal/bot"
	"github.com/p1shiA/ezLink/internal/utils"
	"go.uber.org/zap"
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

	log.Info("Starting the bot")
	client, err := bot.Run(cfg)
	if err != nil {
		log.Fatal("Failed to start the bot", zap.Error(err))
	}

	bh := handlers.NewBotHandler(client, log)
	handlers.RegisterBotHandlers(bh)

	log.Info("Bot started successfully", zap.String("bot_name", client.Self.FirstName))


	client.Idle()


}