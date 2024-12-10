package main

import (
	"context"
	"github.com/BelovN/notifier/internal/bot"
	"github.com/BelovN/notifier/internal/config"
	"github.com/BelovN/notifier/internal/hface"
	"github.com/BelovN/notifier/internal/meteo"
	"github.com/BelovN/notifier/internal/repositories"
	"github.com/BelovN/notifier/internal/weather"
	"github.com/robfig/cron/v3"
	"log"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg := config.LoadConfigFromEnv()

	err, userRepo := repositories.NewSqliteUserRepository(cfg.DbPath)
	if err != nil {
		log.Fatalf("failed to initialize user repository: %v", err)
	}

	err, tgService := bot.NewTelegramService(cfg.TelegramApiToken)
	if err != nil {
		log.Fatalf("failed to initialize Telegram service: %v", err)
	}

	meteoService := meteo.NewService(ctx, nil)

	hfaceService := hface.NewHfaceService(ctx, cfg.HfaceApiToken, nil)

	c := cron.New()
	c.AddFunc("0 9 * * *", func() {

	})

	tgService.AddRouters(
		bot.NewCommandRouter("start", weather.NewStartController(userRepo)),
		bot.NewCommandRouter("subscribe", weather.NewSubscribeController(userRepo)),
		bot.NewCommandRouter("unsubscribe", weather.NewUnSubscribeController(userRepo)),
		bot.NewCommandRouter("weather", weather.NewMeteoController(*meteoService, *hfaceService)),
	)
	tgService.Run(ctx)
}
