package weather

import (
	"context"
	"github.com/BelovN/notifier/internal/bot"
	"github.com/BelovN/notifier/internal/repositories"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type UnSubscribeController struct {
	Controller
}

func NewUnSubscribeController(repo repositories.UserRepository) *UnSubscribeController {
	return &UnSubscribeController{Controller{userRepo: repo}}
}

func (c *UnSubscribeController) Run(ctx context.Context, tg *bot.TelegramService, update *tgbotapi.Update) error {
	username := update.Message.From.UserName
	channelId := update.Message.Chat.ID
	user, err := c.userRepo.GetOrCreateUser(ctx, username, channelId)
	if err != nil {
		return err
	}
	if user.IsSubscribed {
		user.IsSubscribed = false
		if _, err := c.userRepo.Update(ctx, user); err != nil {
			return err
		}
	}
	msgText := "Вы отписались от обновлений"
	if err := tg.SendMessage(update.Message.Chat.ID, msgText); err != nil {
		return err
	}
	return nil
}
