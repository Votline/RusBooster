package bot

import (
	"RusBooster/internal/admin"
	"RusBooster/internal/core"
	"RusBooster/internal/guide"
	"RusBooster/internal/keyboard"
	"RusBooster/internal/stat"
	"RusBooster/internal/state"
	"RusBooster/internal/utils"
	"fmt"
	"go.uber.org/zap"
	tele "gopkg.in/telebot.v3"
	"strconv"
	"strings"
)

func getTargets(log *zap.Logger, message string, userState *state.UserState) (string, *int, *[]string, *tele.ReplyMarkup) {
	if strings.HasPrefix(message, "Найденные слова: \n") {
		menu := keyboard.ShowWordsMenu(userState,
			&userState.CurrentPageOfFindWords, userState.PartsOfFindWords)
		return "Найденные слова: \n",
			&userState.CurrentPageOfFindWords, &userState.PartsOfFindWords, menu
	} else if strings.HasPrefix(message, "Все слова: \n") {
		menu := keyboard.ShowWordsMenu(userState,
			&userState.CurrentPageOfAllWords, userState.PartsOfAllWords)
		return "Все слова: \n",
			&userState.CurrentPageOfAllWords, &userState.PartsOfAllWords, menu
	}
	menu := keyboard.ShowWordsMenu(userState,
		&userState.CurrentPageOfGuide, userState.PartsOfGuide)
	return "", &userState.CurrentPageOfGuide, &userState.PartsOfGuide, menu
}

func setUserField(log *zap.Logger, userId int64, message string, column string) (int, error) {
	userNumber, errUtils := utils.ToInt(message)
	if errUtils != nil {
		return -1, errUtils
	}
	if err := stat.SetSomething(log, userId, userNumber, column); err != nil {
		log.Error("Ошибка при попытке установить значение пользователя", zap.Error(err))
		return -1, err
	}
	return userNumber, nil
}

func deleteOriginalBotMsg(log *zap.Logger, con tele.Context) {
	originalBotMsg := &tele.StoredMessage{
		ChatID:    con.Chat().ID,
		MessageID: strconv.Itoa(con.Message().ID - 1)}
	utils.ActionAfter(log, func() error { return con.Bot().Delete(originalBotMsg) }, 3,
		"Ошибка при попытке удалить первоначальное сообщение бота")

}

func deleteMsgAfter(log *zap.Logger, con tele.Context, temporaryText string) {
	errMessage, _ := con.Bot().Send(con.Chat(), temporaryText)
	deleteMessage := tele.StoredMessage{
		ChatID:    con.Chat().ID,
		MessageID: ""}
	anotherMessage := deleteMessage
	if temporaryText == "" {
		deleteMessage.MessageID = strconv.Itoa(con.Message().ID)
	} else {
		deleteMessage.MessageID = strconv.Itoa(errMessage.ID)
		anotherMessage.MessageID = strconv.Itoa(errMessage.ID - 1)
	}
	utils.ActionAfter(log, func() error { return con.Bot().Delete(deleteMessage) }, 3,
		"Ошибка при попытке удалить ответное сообщение бота")
	if anotherMessage.MessageID != "" {
		utils.ActionAfter(log, func() error { return con.Bot().Delete(anotherMessage) }, 3,
			"Ошибка при попытке удалить сообщение пользователя")
	}
}

func HandleText(log *zap.Logger, bot *tele.Bot, con tele.Context) error {
	userId := con.Sender().ID
	userState, err := state.GetUserState(log, userId)
	defer state.SetUserState(log, userState, userId)
	if err != nil {
		return con.Send(utils.GetReturnText(false))
	}

	switch con.Text() {
	case "Выбрать задание":
		if !userState.IsChecking && !userState.IsSetting {
			userState.IsChoosing = true
			return con.Send("Напишите номер упражнения от 1 до 26: ", keyboard.SelectMenu(log, userId))
		} else {
			return con.Delete()
		}
	case "Проверить знания":
		if !userState.IsSetting && !userState.IsChoosing {
			if userState.IsChecking == true {
				utils.ActionAfter(log,
					func() error {
						return con.Bot().Delete(&tele.StoredMessage{
							ChatID:    con.Chat().ID,
							MessageID: strconv.Itoa(con.Message().ID - 2),
						})
					}, 3, "Ошибка при попытке удалить сообщение пользователя case Проверить Знания")
				deleteOriginalBotMsg(log, con)
			}
			userState.IsChecking = true
			text := core.MakeUserTask(log, userId, userState)
			menu := keyboard.MakeTaskKeyboard()
			message, err := con.Bot().Send(con.Chat(), text, menu)
			utils.ActionAfter(log, func() error {return con.Bot().Delete(message)}, 3600, "Ошибка при попытке удалить check-сообщение")
			return err
		} else {
			return con.Delete()
		}
	case "Статистика":
		if !userState.IsSetting && !userState.IsChoosing && !userState.IsChecking {
			text, menu := keyboard.StatisticMenu(log, con)
			return con.Send(text, menu)
		} else {
			return con.Delete()
		}
	case "Гайд к заданию":
		taskId, _ := stat.GetSomething(log, userId, "current_task")
		text := (guide.ShowGuide(log, taskId,
			&userState.CurrentPageOfGuide, &userState.PartsOfGuide))
		menu := keyboard.ShowWordsMenu(userState, new(int), userState.PartsOfGuide)
		message, err := con.Bot().Send(con.Chat(), text, menu)
		utils.ActionAfter(log, func() error {return con.Bot().Delete(message)}, 3600, "Ошибка при попытке удалить guide-сообщение" )
		return err
	case "/start":
		userState.IsChoosing, userState.IsChecking, userState.IsSetting = false, false, false
		return nil
	default:
		if userId == 5459965917 && utils.ContainsRune(con.Text(), "/") {
			text := admin.HandleCommands(log, userState, userId, con.Text())
			if strings.Contains(con.Text(), "showall") {
				_, _, _, menu := getTargets(log, text, userState)
				message, err := con.Bot().Send(con.Chat(), text, menu)
				utils.ActionAfter(log, func() error {return con.Bot().Delete(message)}, 3600, "Ошибка при попытке удалить admin-сообщение")
				return err 
			}
			return con.Send(text)
		}
		if userState.IsChecking {
			userValue, errUtils := utils.ToInt(con.Text())
			if errUtils != nil {
				log.Error("Ошибка при попытке преобразовать сообщение пользователя в integer")
				deleteMsgAfter(log, con, utils.GetReturnText(false))
				return nil
			}
			userState.IsChecking = false
			text := core.SendAnswer(log, userId, userValue, userState)
			menu := keyboard.MakeAnswerKeyboard()
			return con.Send(text, menu)
		} else if userState.IsSetting {
			userValue, err := setUserField(log, userId, con.Text(), "time_zone")
			if err != nil {
				deleteMsgAfter(log, con, utils.GetReturnText(false))
				return nil
			}
			if !(userValue > -16 && userValue < 16) {
				deleteMsgAfter(log, con, "Введите значение от -15 до 15")
				return nil
			}
			userState.IsSetting = false

			timeZoneForm := utils.GetTimeZoneForm(userValue)
			msg := fmt.Sprintf("Успешно! Ваш часовой пояс изменён на %s", timeZoneForm)
			deleteMsgAfter(log, con, msg)
		} else if userState.IsChoosing {
			userValue, err := setUserField(log, userId, con.Text(), "current_task")
			if err != nil {
				log.Error("Ошибка при попытке внести значение пользователя в current_task")
				msg := utils.GetReturnText(false)
				deleteMsgAfter(log, con, msg)
				return nil
			}
			if !(userValue > 0 && userValue < 27) {
				deleteMsgAfter(log, con, "Введите значение от 1 до 26")
				return nil
			}
			userState.IsChoosing = false

			deleteOriginalBotMsg(log, con)
			msg := "Успешно! Текущее задание: №" + con.Text()
			deleteMsgAfter(log, con, msg)
		}
		return con.Send("Выберите опцию: ", keyboard.MainMenu())
	}
}

func HandleCallback(log *zap.Logger, bot *tele.Bot, con tele.Context) error {
	userId := con.Sender().ID
	data := strings.TrimPrefix(con.Callback().Data, "\f")
	userState, err := state.GetUserState(log, userId)
	defer func() {
		con.Respond()
		if err := state.SetUserState(log, userState, userId); err != nil {
			log.Error("Ошибка сохранения данных в Redis", zap.Error(err))
		}
	}()

	if err != nil {
		log.Error("Ошибка при попытке получить userState")
		return con.Send(utils.GetReturnText(false))
	}
	if userState == nil {
		userState = &state.UserState{}
	}

	switch data {
	case "ToMain":
		userState.IsChoosing, userState.IsChecking, userState.IsSetting = false, false, false
		userMsgID := strconv.Itoa(int(con.Message().ID - 1))
		message := tele.StoredMessage{ChatID: con.Chat().ID, MessageID: userMsgID}

		deleteMsgAfter(log, con, "")
		utils.ActionAfter(log, func() error { return con.Bot().Delete(message) }, 3,
			"Ошибка при попытке удалить сообщение бота")

		return con.Send("Выберите опцию:", keyboard.MainMenu())
	case "SpecifyTimeZone":
		if userState.IsChoosing || userState.IsChecking || userState.IsSetting {
			return nil
		}
		userState.IsSetting = true
		text, menu := keyboard.TimeZoneMenu(log, con)

		return con.Edit(text, menu)
	case "ShowAllExplanations":
		text := strings.ReplaceAll(strings.ReplaceAll(userState.Explanations, "@", ""), "$", "")
		return con.Send(text)
	case "ShowPreviousWords", "ShowNextWords":
		text, targetPage, targetSlice, _ := getTargets(log, con.Message().Text, userState)
		if len(*targetSlice) == 0 {
			deleteMsgAfter(log, con, utils.GetReturnText(false))
			return nil
		}
		if data == "ShowPreviousWords" {
			text += admin.ShowPreviousWords(targetPage, targetSlice)
		} else {
			text += admin.ShowNextWords(targetPage, targetSlice)
		}
		_, _, _, menu := getTargets(log, con.Message().Text, userState)
		return con.Edit(text, menu)
	case "Ignore":
		return nil
	default:
		userState.IsChoosing, userState.IsChecking, userState.IsSetting = false, false, false
		return con.Send("Выберите опцию", keyboard.MainMenu())
	}
}
