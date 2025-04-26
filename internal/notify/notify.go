package notify

import (
	"RusBooster/internal/stat"
	"RusBooster/internal/utils"
	"go.uber.org/zap"
	tele "gopkg.in/telebot.v3"
	"time"
)

type notificationSystem struct {
	bot       *tele.Bot
	scheduler *time.Ticker
}

func NewNotificationSystem(bot *tele.Bot, log *zap.Logger) {
	ns := &notificationSystem{bot: bot}
	go ns.startScheduler(log)
}

func (ns *notificationSystem) startScheduler(log *zap.Logger) {
	ns.scheduler = time.NewTicker(30 * time.Minute)
	for range ns.scheduler.C {
		userIds := stat.GetAllChatIDs(log)
		for _, id := range userIds {
			go ns.sendNotification(id, log)
		}
	}
}

func (ns *notificationSystem) sendNotification(userId int64, log *zap.Logger) {
	lastActiveDate, _ := stat.GetSomething(log, userId, "last_active_date")
	timeZone, _ := stat.GetSomething(log, userId, "time_zone")
	now := time.Now().UTC().Add(time.Duration(timeZone) * time.Hour)
	currentDate := int(now.Unix() / 86400)
	period := ns.calculatePeriod(userId)

	if currentDate - lastActiveDate == 1 {
		time.AfterFunc(time.Duration(period)*time.Hour, func() {
			msg, err := ns.bot.Send(tele.ChatID(userId), "Не забудь выполнить задание чтобы увеличить свой рекорд!")
			if err != nil {
				log.Warn("Ошибка при попытке отправить сообщение для пользователя: ", zap.Int64("userId: ", userId))
			}
			utils.ActionAfter(log,
				func() error { return ns.bot.Delete(msg) }, 10,
				"Ошибка при попытке удалить уведомление")
			ns.sendNotification(userId, log)
		})
	} else if currentDate - lastActiveDate > 1 {
		stat.SetSomething(log, userId, 0, "streak")
		stat.SetSomething(log, userId, 0, "last_active_date")
	}
}

func (ns *notificationSystem) calculatePeriod(userId int64) int {
	loc, _ := time.LoadLocation("Europe/Moscow")
	now := time.Now().In(loc)
	hour := now.Hour()

	switch {
	case hour >= 0 && hour < 12:
		return 6
	case hour >= 12 && hour < 21:
		return 3
	case hour >= 21 && hour < 24:
		return 1
	default:
		return 3
	}
}
