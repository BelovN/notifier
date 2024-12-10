package weather

import (
	"context"
	"github.com/BelovN/notifier/internal/bot"
	"github.com/BelovN/notifier/internal/hface"
	"github.com/BelovN/notifier/internal/meteo"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type MeteoController struct {
	meteoService meteo.Service
	hfaceService hface.HfaceService
}

func NewMeteoController(meteoService meteo.Service, hfaceService hface.HfaceService) *MeteoController {
	return &MeteoController{meteoService, hfaceService}
}

func (c *MeteoController) Run(ctx context.Context, tg *bot.TelegramService, update *tgbotapi.Update) error {
	weather, err := c.meteoService.GetCurrentWeather()
	if err != nil {
		return err
	}

	response, err := c.hfaceService.GetAIAnswer(weather.ToString())
	if err != nil {
		return nil
	}

	if err := tg.SendMessage(update.Message.Chat.ID, response); err != nil {
		return err
	}
	return nil
}
