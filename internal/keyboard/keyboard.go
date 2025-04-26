package keyboard

import (
	"RusBooster/internal/stat"
	"RusBooster/internal/state"
	"RusBooster/internal/utils"
	"fmt"
	"go.uber.org/zap"
	tele "gopkg.in/telebot.v3"
)

func ShowWordsMenu(userState *state.UserState, currentPage *int, currentSlice []string) *tele.ReplyMarkup {
	selector := &tele.ReplyMarkup{}
	allPagesText := fmt.Sprintf("[%d/%d]",
		*currentPage+1,
		len(currentSlice))

	btnShowPrevious := selector.Data("<", "ShowPreviousWords")
	btnShowNext := selector.Data(">", "ShowNextWords")
	btnAllPages := selector.Data(allPagesText, "AllPages")

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

func TimeZoneMenu(log *zap.Logger, con tele.Context) (string, *tele.ReplyMarkup) {
	userId := con.Sender().ID
	timeZone, err := stat.GetSomething(log, userId, "time_zone")
	if err != nil {
		log.Error("Ошибка при получении часового пояса", zap.Error(err))
		return utils.GetReturnText(false), nil
	}
	timeZoneForm := utils.GetTimeZoneForm(timeZone)
	msg := fmt.Sprintf("Текущий часовой пояс: %s\nНапишите свою разницу во времени от МСК:", timeZoneForm)
	selector := &tele.ReplyMarkup{}

	btnToMain := selector.Data("В главное меню", "ToMain")

	selector.Inline(selector.Row(btnToMain))
	return msg, selector
}

func StatisticMenu(log *zap.Logger, con tele.Context) (string, *tele.ReplyMarkup) {
	msg := stat.GetStatistic(log, con.Sender().FirstName, con.Sender().ID)
	selector := &tele.ReplyMarkup{}

	btnSpecifyTimeZone := selector.Data("Указать часовой пояс", "SpecifyTimeZone")
	btnToMain := selector.Data("Вернуться в главное меню", "ToMain")

	selector.Inline(
		selector.Row(btnSpecifyTimeZone),
		selector.Row(btnToMain),
	)
	return msg, selector
}

func SelectMenu(log *zap.Logger, userId int64) *tele.ReplyMarkup {
	selector := &tele.ReplyMarkup{}

	worstTaskScore, err := stat.GetSomething(log, userId, "worst_task_result")
	if err != nil {
		log.Error("Ошибка при получении данных из rsdb", zap.Error(err))
		return nil
	}
	btnWorstTask := selector.Data(fmt.Sprintf("Наихудшая успеваимость: №%d", worstTaskScore), " ")
	btnToMain := selector.Data("Вернуться в главное меню", "ToMain")

	selector.Inline(
		selector.Row(btnWorstTask),
		selector.Row(btnToMain),
	)
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
