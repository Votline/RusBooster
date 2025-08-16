package keyboard

import (
	"fmt"
	"log"

	tele "gopkg.in/telebot.v3"

	"RusBooster/internal/stat"
	"RusBooster/internal/state"
	"RusBooster/internal/utils"
)

func ShowWordsMenu(userState *state.UserState, currentPage *int, currentSlice []string) *tele.ReplyMarkup {
	selector := &tele.ReplyMarkup{}
	allPagesText := fmt.Sprintf("[%d/%d]",
		*currentPage+1,
		len(currentSlice))

	btnShowPrevious := selector.Data("<", "ShowPreviousWords")
	btnShowNext := selector.Data(">", "ShowNextWords")
	btnAllPages := selector.Data(allPagesText, "Ignore")

	selector.Inline(
		selector.Row(btnShowPrevious, btnShowNext),
		selector.Row(btnAllPages),
	)
	return selector
}

func MakeAnswerKeyboard() *tele.ReplyMarkup {
	selector := &tele.ReplyMarkup{}
	btnShowAllExpl := selector.Data("Показать пояснения всех слов", "ShowAllExplanations")
	selector.Inline(selector.Row(btnShowAllExpl))
	return selector
}

func MakeTaskKeyboard() *tele.ReplyMarkup {
	selector := &tele.ReplyMarkup{}
	btnToMain := selector.Data("Вернуться в главное меню", "ToMain")
	selector.Inline(selector.Row(btnToMain))
	return selector
}

func TimeZoneMenu(con tele.Context) (string, *tele.ReplyMarkup) {
	userId := con.Sender().ID
	timeZone, err := stat.GetSomething(userId, "time_zone")
	if err != nil {
		log.Printf("Ошибка при получении часового пояса: %v", err)
		timeZone = 0
	}
	timeZoneForm := utils.GetTimeZoneForm(timeZone)
	text := fmt.Sprintf("Текущий часовой пояс: %s\nВведите новое значение от -15 до 15:", timeZoneForm)
	selector := &tele.ReplyMarkup{}
	btnCancel := selector.Data("Отмена", "Cancel")
	selector.Inline(selector.Row(btnCancel))
	return text, selector
}

func StatisticMenu(con tele.Context) (string, *tele.ReplyMarkup) {
	userId := con.Sender().ID
	userName := con.Sender().Username
	if userName == "" {
		userName = con.Sender().FirstName
	}
	text := stat.GetStatistic(userName, userId)
	selector := &tele.ReplyMarkup{}
	btnBack := selector.Data("Назад", "Back")
	selector.Inline(selector.Row(btnBack))
	return text, selector
}

func SelectMenu(userId int64) *tele.ReplyMarkup {
	_, err := stat.GetSomething(userId, "current_task")
	if err != nil {
		log.Printf("Ошибка при получении данных из rsdb: %v", err)
	}
	selector := &tele.ReplyMarkup{}
	btnCancel := selector.Data("Отмена", "Cancel")
	selector.Inline(selector.Row(btnCancel))
	return selector
}

func MainMenu() *tele.ReplyMarkup {
	menu := &tele.ReplyMarkup{ResizeKeyboard: true}

	btnSelectTask := menu.Text("Выбрать задание")
	btnTestKnowledge := menu.Text("Проверить знания")
	btnStatistic := menu.Text("Статистика")
	btnShowGuide := menu.Text("Гайд к заданию")

	menu.Reply(
		menu.Row(btnSelectTask, btnTestKnowledge, btnStatistic),
		menu.Row(btnShowGuide),
	)
	return menu
}
