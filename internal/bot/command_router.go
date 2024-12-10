package bot

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type CommandRouter struct {
	Router
}

func NewCommandRouter(name string, controllers ...IController) *CommandRouter {
	resolver := NewCommandResolver(name)
	return &CommandRouter{Router: *NewRouter(name, resolver, controllers...)}
}

func (r *CommandRouter) Resolve(update *tgbotapi.Update) bool {
	return r.Router.Resolve(update)
}

func (r *CommandRouter) RunControllers(ctx context.Context, tg *TelegramService, update *tgbotapi.Update) {
	r.Router.RunControllers(ctx, tg, update)
}
