package main

import (
	"os"
	"log"
	"time"

	tele "gopkg.in/telebot.v3"

	"RusBooster/internal/bot"
	"RusBooster/internal/notify"
)

func main() {
	log.SetFlags(log.Lshortfile)
	
	token := os.Getenv("TOKEN")

	b, err := tele.NewBot(tele.Settings{
		Token: token,
		Poller: &tele.LongPoller{
			Timeout:        10 * time.Second,
		},
	})
	if err != nil {
		log.Fatalf("Не удалось запустить бота: %v", err.Error())
	}

	b.Handle(tele.OnText, func(con tele.Context) error {
		return bot.HandleText(b, con)
	})

	b.Handle(tele.OnCallback, func(con tele.Context) error {
		return bot.HandleCallback(b, con)
	})

	notify.NewNotificationSystem(b)
	b.Start()
}
