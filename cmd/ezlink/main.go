package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/p1shiA/ezLink/config"
	"github.com/p1shiA/ezLink/internal/api"
	"github.com/p1shiA/ezLink/internal/bot"
	"github.com/p1shiA/ezLink/internal/bot/handlers"
	"github.com/p1shiA/ezLink/internal/core/repository"
	"github.com/p1shiA/ezLink/internal/core/service"
	"github.com/p1shiA/ezLink/internal/database"
	"github.com/p1shiA/ezLink/internal/payment/zarinpal"
	"github.com/p1shiA/ezLink/internal/utils"
	"go.uber.org/zap"
)

func main() {
	// 1. Loading env
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	// 2. Start logging
	log, err := utils.Init(cfg)
	if err != nil {
		panic(err)
	}
	defer log.Sync()

	log.Info("Starting ezLink Application")

	// 3. Connect to Database
	db, err := database.InitDB(cfg.DbUrl)
	if err != nil {
		log.Fatal("Failed to connect to PostgreSQL", zap.Error(err))
	}

	// 4. Run SQL migrations
	if err := database.RunMigrations(context.Background(), db); err != nil {
		log.Fatal("Failed to run database migrations", zap.Error(err))
	}

	// 5. Seed default plans and packages if missing
	if err := database.SeedDefaults(context.Background(), db); err != nil {
		log.Warn("Failed to seed database defaults", zap.Error(err))
	}

	// 5. Initialize Unit-Of-Work / TxManager & Repositories
	txManager := repository.NewGormTxManager(db)
	userRepo := repository.NewUserRepository(db)
	planRepo := repository.NewPlanRepository(db)
	subRepo := repository.NewSubscriptionRepository(db)
	txRepo := repository.NewTransactionRepository(db)
	fileLogRepo := repository.NewFileLogRepository(db)
	adminRepo := repository.NewAdminRepository(db)

	// 6. Initialize Payment Gateway (Zarinpal v4)
	gateway := zarinpal.NewClient(zarinpal.Config{
		MerchantID: cfg.PaymentConfig.MerchantID,
		Sandbox:    cfg.PaymentConfig.Sandbox,
	})

	// 7. Initialize Services
	userService := service.NewUserService(userRepo)
	planService := service.NewPlanService(planRepo)
	subService := service.NewSubscriptionService(subRepo, txManager)
	txService := service.NewTransactionService(
		txRepo, planRepo, subRepo, subService, gateway, txManager, cfg.PaymentConfig.CallbackURL,
	)
	fileLogService := service.NewFileLogService(fileLogRepo, subService, txManager)
	adminService := service.NewAdminService(adminRepo)

	// 8. Start Telegram Bot
	log.Info("Starting Telegram Bot")
	client, err := bot.Run(cfg)
	if err != nil {
		log.Fatal("Failed to start the bot", zap.Error(err))
	}

	// 9. Register Bot Handlers with Services
	services := &handlers.Services{
		User:         userService,
		Plan:         planService,
		Subscription: subService,
		Transaction:  txService,
		FileLog:      fileLogService,
		Admin:        adminService,
	}

	bh := handlers.NewBotHandler(client, services, log)
	handlers.RegisterBotHandlers(bh)
	log.Info("Bot started successfully", zap.String("bot_name", client.Self.FirstName))

	// 10. Start Gin HTTP Callback Server
	notifier := bot.NewBotNotifier(client, log)
	apiServer := api.NewServer(cfg, txService, notifier, client.Self.Username, log)
	go func() {
		if err := apiServer.Start(); err != nil {
			log.Error("HTTP Server error", zap.Error(err))
		}
	}()

	// 11. Graceful Shutdown on interrupt signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigChan
		log.Info("Shutting down ezLink...")

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := apiServer.Shutdown(shutdownCtx); err != nil {
			log.Error("HTTP Server shutdown error", zap.Error(err))
		}

		client.Stop()
		os.Exit(0)
	}()

	client.Idle()
}