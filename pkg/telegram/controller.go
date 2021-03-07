package telegram

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

// Consts ...
const (
	START = "start"
)

// start ...
func (bot *Bot) start(msg *tgbotapi.Message) {
	message := tgbotapi.NewMessage(msg.Chat.ID, "Pong !")
	message.ReplyToMessageID = msg.MessageID
	bot.messages <- message
}
