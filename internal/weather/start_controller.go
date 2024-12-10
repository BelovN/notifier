package weather

import (
	"context"
	"github.com/BelovN/notifier/internal/bot"
	"github.com/BelovN/notifier/internal/repositories"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type StartController struct {
	Controller
}

func NewStartController(repo repositories.UserRepository) *StartController {
	return &StartController{Controller{repo}}
}

func (c *StartController) Run(ctx context.Context, tg *bot.TelegramService, update *tgbotapi.Update) error {
	username := update.Message.From.UserName
	channelId := update.Message.Chat.ID
	if _, err := c.Controller.userRepo.GetOrCreateUser(ctx, username, channelId); err != nil {
		return err
	}
	msgText := "Добро пожаловать в Бота, теперь вы подписаны на погоду в Белграде"

	if err := tg.SendMessage(update.Message.Chat.ID, msgText); err != nil {
		return err
	}
	return nil
}
