package bot

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

type IResolver interface {
	ResolveCondition(update *tgbotapi.Update) bool
}

type CommandResolver struct {
	name string
}

func NewCommandResolver(name string) *CommandResolver {
	return &CommandResolver{name}
}

func (r *CommandResolver) ResolveCondition(update *tgbotapi.Update) bool {
	if update.Message == nil {
		return false
	}
	return update.Message.Command() == r.name
}
