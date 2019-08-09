package archive_builder_bot

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

type UpdateHandler interface {
	Handle(update tgbotapi.Update)
}

func NewMessageHandler(b *tgbotapi.BotAPI) UpdateHandler {
	return &handler{
		bot: b,
	}
}

type handler struct {
	bot *tgbotapi.BotAPI
}

func (h *handler) Handle(update tgbotapi.Update) {

}
