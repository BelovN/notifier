package weather

import (
	"context"
	bot "github.com/BelovN/notifier/internal/bot"
	"github.com/BelovN/notifier/internal/repositories"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type SubscribeController struct {
	Controller
}

func NewSubscribeController(repo repositories.UserRepository) *SubscribeController {
	return &SubscribeController{Controller{userRepo: repo}}
}

func (c *SubscribeController) Run(ctx context.Context, tg *bot.TelegramService, update *tgbotapi.Update) error {
	username := update.Message.From.UserName
	channelId := update.Message.Chat.ID
	user, err := c.Controller.userRepo.GetOrCreateUser(ctx, username, channelId)
	if err != nil {
		return err
	}
	if !user.IsSubscribed {
		user.IsSubscribed = true
		if _, err := c.Controller.userRepo.Update(ctx, user); err != nil {
			return err
		}
	}
	msgText := "Вы подписались на обновления"
	if err := tg.SendMessage(update.Message.Chat.ID, msgText); err != nil {
		return err
	}
	return nil
}
