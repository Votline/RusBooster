package bot

import (
	"fmt"
	"log"
	"strings"

	tele "gopkg.in/telebot.v3"

	"RusBooster/internal/core"
	"RusBooster/internal/stat"
	"RusBooster/internal/state"
	"RusBooster/internal/utils"
	"RusBooster/internal/guide"
	"RusBooster/internal/admin"
	"RusBooster/internal/keyboard"
)

func HandleText(bot *tele.Bot, con tele.Context) error {
	userId := con.Sender().ID
	userState, err := state.GetUserState(userId)
	defer state.SetUserState(userState, userId)
	if err != nil {
		return con.Send(utils.GetReturnText(false))
	}

	switch con.Text() {
	case "Выбрать задание":
		if !userState.IsChecking && !userState.IsSetting {
			userState.IsChoosing = true
			return con.Send("Напишите номер упражнения от 1 до 26: ", keyboard.SelectMenu(userId))
		} else {
			return con.Delete()
		}
	case "Проверить знания":
		if !userState.IsSetting && !userState.IsChoosing {
			userState.IsChecking = true
			
			text := core.MakeUserTask(userId, userState)
			menu := keyboard.MakeTaskKeyboard()
		
			return con.Send(text, menu)
		}
	case "Статистика":
		if !userState.IsSetting && !userState.IsChoosing && !userState.IsChecking {
			text, menu := keyboard.StatisticMenu(con)
			return con.Send(text, menu)
		} else {
			return con.Delete()
		}
	case "Гайд к заданию":
		taskId, _ := stat.GetSomething(userId, "current_task")
		text := (guide.ShowGuide(taskId,
			&userState.CurrentPageOfGuide, &userState.PartsOfGuide))
		menu := keyboard.ShowWordsMenu(userState, new(int), userState.PartsOfGuide)
		return con.Send(text, menu)
	default:
		if userId == 5459965917 && containsRune(con.Text(), "/") {
			text := admin.HandleCommands(userState, userId, con.Text())
			if strings.Contains(con.Text(), "showall") {
				_, _, _, menu := getTargets(text, userState)
				return con.Send(text, menu)
			}
			return con.Send(text)
		} else if userState.IsChecking {
			userValue, errUtils := utils.ToInt(con.Text())
			if errUtils != nil {
				log.Printf("Ошибка при попытке преобразовать сообщение пользователя в integer: %v", errUtils)
				return nil
			}
			userState.IsChecking = false
			text, err := core.SendAnswer(userId, userValue, userState)
			if err != nil {
				text = (utils.GetReturnText(false) + "\nДанные не сохранены")
				return con.Send(text)
			}
			menu := keyboard.MakeAnswerKeyboard()
			return con.Send(text, menu)
		} else if userState.IsSetting {
			userValue, err := setUserField(userId, con.Text(), "time_zone")
			if err != nil {
				return nil
			}
			if !(userValue > -16 && userValue < 16) {
				return nil
			}
			userState.IsSetting = false
			timeZoneForm := utils.GetTimeZoneForm(userValue)
			msg := fmt.Sprintf("Успешно! Ваш часовой пояс изменён на %s", timeZoneForm)
			return con.Send(msg)
		} else if userState.IsChoosing {
			userValue, err := setUserField(userId, con.Text(), "current_task")
			if err != nil {
				log.Printf("Ошибка при попытке внести значение пользователя в current_task: %v", err)
				msg := utils.GetReturnText(false)
				return con.Send(msg)
			}
			if !(userValue > 0 && userValue < 27) {
				return nil
			}
			userState.IsChoosing = false

			msg := "Успешно! Текущее задание: №" + con.Text()
			return con.Send(msg)
		}
		return con.Send("Выберите опцию: ", keyboard.MainMenu())
	}
	return nil
}

func HandleCallback(bot *tele.Bot, con tele.Context) error {
	userId := con.Sender().ID
	data := strings.TrimPrefix(con.Callback().Data, "\f")
	userState, err := state.GetUserState(userId)
	defer func() {
		con.Respond()
		if err := state.SetUserState(userState, userId); err != nil {
			log.Printf("Ошибка сохранения данных в Redis: %v", err)
		}
	}()

	if err != nil {
		log.Printf("Ошибка при попытке получить userState: %v", err)
		return con.Send(utils.GetReturnText(false))
	}

	switch data {
	case "ToMain":
		userState.IsChecking = false
		userState.IsSetting = false
		userState.IsChoosing = false
		return con.Send("Выберите опцию: ", keyboard.MainMenu())
	
	case "Cancel":
		userState.IsChecking = false
		userState.IsSetting = false
		userState.IsChoosing = false
		return con.Send("Выберите опцию: ", keyboard.MainMenu())
	
	case "Back":
		return con.Send("Выберите опцию: ", keyboard.MainMenu())
	
	case "SpecifyTimeZone":
		userState.IsSetting = true
		text, menu := keyboard.TimeZoneMenu(con)
		return con.Edit(text, menu)

	case "ShowAllExplanations":
		text := userState.Explanations
		return con.Send(text)
	
	case "ShowPreviousWords", "ShowNextWords":
		text, targetPage, targetSlice, _ := getTargets(con.Message().Text, userState)
		if len(*targetSlice) == 0 {
			return nil
		}
		if data == "ShowPreviousWords" {
			text += admin.ShowPreviousWords(targetPage, targetSlice)
		} else {
			text += admin.ShowNextWords(targetPage, targetSlice)
		}
		_, _, _, menu := getTargets(con.Message().Text, userState)
		return con.Edit(text, menu)
	case "Ignore":
		return nil
	default:
		return con.Send("Неизвестная команда")
	}
}
