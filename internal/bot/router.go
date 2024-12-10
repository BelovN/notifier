package bot

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type IRouter interface {
	Resolve(update *tgbotapi.Update) bool
	RunControllers(ctx context.Context, tg *TelegramService, update *tgbotapi.Update)
}

type Router struct {
	name        string
	resolver    IResolver
	controllers []IController
}

func NewRouter(name string, resolver IResolver, controllers ...IController) *Router {
	return &Router{name: name, resolver: resolver, controllers: controllers}
}

func (r *Router) Resolve(update *tgbotapi.Update) bool {
	return r.resolver.ResolveCondition(update)
}

func (r *Router) RunControllers(ctx context.Context, tg *TelegramService, update *tgbotapi.Update) {
	for _, controller := range r.controllers {
		go func() {
			if err := controller.Run(ctx, tg, update); err != nil {
				fmt.Println(err)
			}
		}()
	}
}
