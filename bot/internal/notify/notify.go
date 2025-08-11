package notify

import (
	"log"
	"time"

	tele "gopkg.in/telebot.v3"

	"RusBooster/internal/stat"
	"RusBooster/internal/utils"
)

type notificationSystem struct {
	bot       *tele.Bot
	scheduler *time.Ticker
}

func NewNotificationSystem(bot *tele.Bot) {
	ns := &notificationSystem{bot: bot}
	go ns.startScheduler()
}

func (ns *notificationSystem) startScheduler() {
	ns.scheduler = time.NewTicker(30 * time.Minute)
	for range ns.scheduler.C {
		userIds := stat.GetAllChatIDs()
		for _, id := range userIds {
			go ns.sendNotification(id)
		}
	}
}

func (ns *notificationSystem) sendNotification(userId int64) {
	lastActiveDate, _ := stat.GetSomething(userId, "last_active_date")
	timeZone, _ := stat.GetSomething(userId, "time_zone")
	now := time.Now().UTC().Add(time.Duration(timeZone) * time.Hour)
	currentDate := int(now.Unix() / 86400)
	period := ns.calculatePeriod(userId)

	if currentDate - lastActiveDate == 1 {
		time.AfterFunc(time.Duration(period)*time.Hour, func() {
			msg, err := ns.bot.Send(tele.ChatID(userId), "Не забудь выполнить задание чтобы увеличить свой рекорд!")
			if err != nil {
				log.Printf("Ошибка при попытке отправить сообщение для пользователя %d: %v", userId, err)
			}
			utils.ActionAfter(
				func() error { return ns.bot.Delete(msg) }, 10,
				"Ошибка при попытке удалить уведомление")
			ns.sendNotification(userId)
		})
	} else if currentDate - lastActiveDate > 1 {
		stat.SetSomething(userId, 0, "streak")
		stat.SetSomething(userId, 0, "last_active_date")
	}
}

func (ns *notificationSystem) calculatePeriod(userId int64) int {
	streak, _ := stat.GetSomething(userId, "streak")
	if streak < 3 {
		return 24
	} else if streak < 7 {
		return 12
	} else if streak < 14 {
		return 8
	} else {
		return 6
	}
}
