package bot

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TelegramService struct {
	bot     *tgbotapi.BotAPI
	routers []IRouter
}

func NewTelegramService(token string) (error, *TelegramService) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return err, nil
	}
	return nil, &TelegramService{bot: bot}
}

func (tg *TelegramService) AddRouters(routers ...IRouter) {
	tg.routers = append(tg.routers, routers...)
}

func (tg *TelegramService) Run(ctx context.Context) {
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 30
	updates := tg.bot.GetUpdatesChan(updateConfig)
	for {
		select {
		case update := <-updates:
			tg.onUpdate(ctx, &update)
		case <-ctx.Done():
			return
		}
	}
}

func (tg *TelegramService) SendMessage(chatId int64, msg string) error {
	message := tgbotapi.NewMessage(chatId, msg)
	if _, err := tg.bot.Send(message); err != nil {
		return err
	}
	return nil
}

func (tg *TelegramService) onUpdate(ctx context.Context, update *tgbotapi.Update) {
	for _, router := range tg.routers {
		if router.Resolve(update) {
			router.RunControllers(ctx, tg, update)
		}
	}
}
