package telegram

import (
	"os"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/sirupsen/logrus"
)

// Bot ...
type Bot struct {
	BotAPI   *tgbotapi.BotAPI
	messages chan tgbotapi.Chattable
}

// NewBot ...
func NewBot() (*Bot, error) {
	var err error
	var bot Bot

	bot.BotAPI, err = tgbotapi.NewBotAPI(os.Getenv("API_TOKEN"))
	if err != nil {
		return nil, err
	}

	bot.BotAPI.Debug = false
	bot.BotAPI.Buffer = 12 * 15
	bot.messages = make(chan tgbotapi.Chattable, 300)

	return &bot, nil
}

// StartPolling ...
func (bot *Bot) StartPolling(stopChan chan struct{}) {
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60
	updateChannel, err := bot.BotAPI.GetUpdatesChan(updateConfig)
	if err != nil {
		logrus.Fatalf("StartPolling error receive GetUpdatesChan: %s", err.Error())
	}

	go bot.sendWrapper(500)

	for update := range updateChannel {
		select {
		case <-stopChan:
			{
				logrus.Fatalf("Stop Long Pooling !")
				return
			}
		default:
			{
				go bot.distributeUpdate(update)
			}
		}
	}
}

// distributeUpdate ...
func (bot *Bot) distributeUpdate(update tgbotapi.Update) {

	if update.Message.IsCommand() != false {
		if !bot.distributeCommand(update.Message) {
			message := tgbotapi.NewMessage(update.Message.Chat.ID, "Error command!")
			message.ReplyToMessageID = update.Message.MessageID
			bot.messages <- message
		}
	}

	if update.Message.Text != "" && update.Message.IsCommand() == false {
		if !bot.distributeMessage(update.Message) {
			message := tgbotapi.NewMessage(update.Message.Chat.ID, "Error message!")
			message.ReplyToMessageID = update.Message.MessageID
			bot.messages <- message
		}
	}
}

// distributeCommand ...
func (bot *Bot) distributeCommand(message *tgbotapi.Message) bool {
	var isRightCommand = true

	command := message.Command()
	if command == "" {
		isRightCommand = false
	}

	switch command {
	case START:
		{
			go bot.start(message)
		}
	default:
		{
			isRightCommand = false
		}
	}

	return isRightCommand
}

// distributeMessage ...
func (bot *Bot) distributeMessage(message *tgbotapi.Message) bool {
	if message.Text == "" {
		return false
	}

	//

	return true
}

// Sender ...
func (bot *Bot) Sender(msg tgbotapi.Chattable) {
	if _, err := bot.BotAPI.Send(msg); err != nil {
		panic(err)
	}
}

// sendWrapper ...
func (bot *Bot) sendWrapper(millisecond uint64) {
	nanosecond := millisecond * 1000000
	rate := time.Duration(nanosecond)
	limiter := time.Tick(rate)
	for message := range bot.messages {
		<-limiter
		go bot.Sender(message)
	}
}
