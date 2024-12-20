package periodic

import (
	"context"
	"fmt"
	"github.com/BelovN/notifier/internal/bot"
	"github.com/BelovN/notifier/internal/hface"
	"github.com/BelovN/notifier/internal/meteo"
	"github.com/BelovN/notifier/internal/repositories"
	"github.com/robfig/cron/v3"
	"log"
)

type Weather struct {
	CronTimeSheet string
	meteoService  meteo.Service
	hfaceService  hface.Service
	userRepo      repositories.UserRepository
	bot           bot.TelegramService
	ctx           context.Context
}

func NewPeriodicWeather(
	cronTimeSheet string,
	meteoService meteo.Service,
	hfaceService hface.Service,
	userRepo repositories.UserRepository,
	bot bot.TelegramService,
	ctx context.Context,
) *Weather {

	return &Weather{
		cronTimeSheet,
		meteoService,
		hfaceService,
		userRepo,
		bot,
		ctx,
	}
}

func (w *Weather) periodicSync() {
	fmt.Println("RUN PERIODIC")

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
		err = w.bot.SendMessage(user.ChannelId, response)
		if err != nil {
			log.Println("error")
		}
	}
}

func (w *Weather) Run() error {
	c := cron.New()
	if _, err := c.AddFunc(w.CronTimeSheet, w.periodicSync); err != nil {
		return err
	}
	c.Start()
	return nil
}
