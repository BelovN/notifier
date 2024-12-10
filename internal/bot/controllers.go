package bot

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type IController interface {
	Run(ctx context.Context, tg *TelegramService, update *tgbotapi.Update) error
}
