package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Mukxtarov/bot-structure/pkg/telegram"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

func main() {
	//JSONFormatter formats logs into parsable json
	logrus.SetFormatter(new(logrus.JSONFormatter))

	//Load will read your env file(s) and load them into ENV for this process.
	if err := godotenv.Load(); err != nil {
		logrus.Fatalf("Error env open: %s", err.Error())
	}
	logrus.Infoln("Start Program !")

	logrus.Infoln("Bot initialization !")

	CallBot, err := telegram.NewBot()
	if err != nil {
		logrus.Fatalf("Attempt to log into the bot: %s", err.Error())
	}

	logrus.Infoln("Start Bot !")

	stopChan := make(chan struct{})
	go CallBot.StartPolling(stopChan)

	// Intercepting signals
	signChan := make(chan os.Signal, 1)
	// SIGTERM for Server (htop kill 15), SIGINT for Windows (Ctrl + C)
	signal.Notify(signChan, syscall.SIGTERM, syscall.SIGINT)
	// Wait signal
	<-signChan
	// Stop bot
	close(stopChan)
	// Wait for all functions to complete.
	// Due to the bot no longer accepting new messages, functions will not be called
	time.Sleep(2 * time.Second)

	logrus.Infoln("Stop robot")

	time.Sleep(500 * time.Millisecond)

}
