package telegram

import (
	"os"

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

func (bot *Bot) StartPolling(stopChan chan struct{}) {
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60
	updateChannel, err := bot.BotAPI.GetUpdatesChan(updateConfig)
	if err != nil {
		logrus.Fatalf("StartPolling error receive GetUpdatesChan: %s", err.Error())
	}

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

func (bot *Bot) distributeCommand(message *tgbotapi.Message) bool {
	var isRightCommand = true

	command := message.Command()
	if command == "" {
		return false
	}

	switch command {
	case "start":
		{
			//
		}
	default:
		{
			return false
		}
	}

	return isRightCommand
}

func (bot *Bot) dictributeMessage(message *tgbotapi.Message) bool {
	if message.Text == "" {
		return false
	}

	//

	return true
}

func (bot *Bot) Sender(msg tgbotapi.Chattable) {
	if _, err := bot.BotAPI.Send(msg); err != nil {
		panic(err)
	}
}
