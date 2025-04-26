package main

import (
	"RusBooster/internal/bot"
	"RusBooster/internal/notify"
	"bufio"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	tele "gopkg.in/telebot.v3"
	"os"
)

func getToken() string {
	tokenFile, tokenErr := os.OpenFile("token.txt", os.O_RDONLY, 0644)
	if tokenErr != nil {
		panic(tokenErr)
	}
	defer tokenFile.Close()

	scanner := bufio.NewScanner(tokenFile)
	if scanner.Scan() {
		return scanner.Text()
	} else {
		panic("Токен не найден")
	}
}

func createLogger() *zap.Logger {
	file, fileErr := os.OpenFile("logs/fatal.log", os.O_WRONLY, 0644)
	if fileErr != nil {
		panic(fileErr)
	}

	consoleCore := zapcore.NewCore(
		zapcore.NewConsoleEncoder(zap.NewProductionEncoderConfig()),
		zapcore.AddSync(os.Stdout),
		zapcore.InfoLevel,
	)
	fatalCore := zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
		zapcore.AddSync(file),
		zapcore.WarnLevel,
	)
	combinedCore := zapcore.NewTee(consoleCore, fatalCore)
	log := zap.New(combinedCore, zap.AddCaller())
	return log
}

func main() {
	token := getToken()
	log := createLogger()
	defer log.Sync()

	webhookURL := os.Getenv("WEBHOOK_URL")
	if webhookURL == "" {
		log.Fatal("Webhook_url не задан")
	}

	b, err := tele.NewBot(tele.Settings{
		Token: token,
		Poller: &tele.Webhook{
			Listen:   "0.0.0.0:8080",
			Endpoint: &tele.WebhookEndpoint{PublicURL: webhookURL},
		},
	})
	if err != nil {
		log.Fatal("Не удалось подключиться к cloudflare", zap.Error(err))
	}
	log.Info("Бот запущен на вебхуках!")
	b.Handle(tele.OnText, func(con tele.Context) error {
		return bot.HandleText(log, b, con)
	})
	b.Handle(tele.OnCallback, func(con tele.Context) error {
		return bot.HandleCallback(log, b, con)
	})
	notify.NewNotificationSystem(b, log)
	b.Start()
}
