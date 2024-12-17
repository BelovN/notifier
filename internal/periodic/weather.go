package periodic

import (
	"context"
	"github.com/BelovN/notifier/internal/bot"
	"github.com/BelovN/notifier/internal/hface"
	"github.com/BelovN/notifier/internal/meteo"
	"github.com/BelovN/notifier/internal/repositories"
	"github.com/robfig/cron/v3"
)

const (
	DefaultTimeSheet = "0 9 * * *"
)

type PeriodicWeather struct {
	CronTimeSheet string
	meteoService  meteo.Service
	hfaceService  hface.HfaceService
	userRepo      repositories.UserRepository
	bot           bot.TelegramService
	ctx           context.Context
}

func NewPeriodicWeather(
	cronTimeSheet string,
	meteoService meteo.Service,
	hfaceService hface.HfaceService,
	userRepo repositories.UserRepository,
	bot bot.TelegramService,
	ctx context.Context,
) *PeriodicWeather {

	if cronTimeSheet == "" {
		cronTimeSheet = DefaultTimeSheet
	}
	return &PeriodicWeather{
		cronTimeSheet,
		meteoService,
		hfaceService,
		userRepo,
		bot,
		ctx,
	}
}

func (w *PeriodicWeather) periodicSync() {
	filters := map[string]interface{}{
		"is_subscribed": true,
	}

	users, err := w.userRepo.FilterUsers(w.ctx, filters)
	if err != nil {
		return
	}

	weather, err := w.meteoService.GetCurrentWeather()
	if err != nil {
		return
	}

	response, err := w.hfaceService.GetAIAnswer(weather.ToString())
	if err != nil {
		return
	}

	for _, user := range users {
		w.bot.SendMessage(user.ChannelId, response)
	}
}

func (w *PeriodicWeather) Run() error {
	c := cron.New()
	if _, err := c.AddFunc(w.CronTimeSheet, w.periodicSync); err != nil {
		return err
	}
	return nil
}
