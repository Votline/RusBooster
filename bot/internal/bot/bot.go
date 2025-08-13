package bot

import (
	"fmt"
	"log"
	"strings"
	"strconv"

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
			if userState.IsChecking == true {
				utils.ActionAfter(
					func() error {
						return con.Bot().Delete(&tele.StoredMessage{
							ChatID:    con.Chat().ID,
							MessageID: strconv.Itoa(con.Message().ID),
						})
					}, 3, "Ошибка при попытке удалить сообщение пользователя case Проверить Знания")
				deleteOriginalBotMsg(con)
			}
			userState.IsChecking = true
			text := core.MakeUserTask(userId, userState)
			menu := keyboard.MakeTaskKeyboard()
			message, _ := con.Bot().Send(con.Chat(), text, menu)
			utils.ActionAfter(func() error {
				if message == nil {
					return nil
				}
				return con.Bot().Delete(&tele.StoredMessage{
					ChatID:    con.Chat().ID,
					MessageID: strconv.Itoa(message.ID),
				})
			}, 30, "Ошибка при попытке удалить сообщение с заданием")
			return nil
		} else {
			return con.Delete()
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
		message, _ := con.Bot().Send(con.Chat(), text, menu)
		utils.ActionAfter(func() error {
			if message == nil {
				return nil
			}
			return con.Bot().Delete(&tele.StoredMessage{
				ChatID:    con.Chat().ID,
				MessageID: strconv.Itoa(message.ID),
			})
		}, 30, "Ошибка при попытке удалить сообщение с гайдом")
		return nil
	default:
		if userId == 5459965917 && containsRune(con.Text(), "/") {
			text := admin.HandleCommands(userState, userId, con.Text())
			if strings.Contains(con.Text(), "showall") {
				_, _, _, menu := getTargets(text, userState)
				message, _ := con.Bot().Send(con.Chat(), text, menu)
				utils.ActionAfter(func() error {
					if message == nil {
						return nil
					}
					return con.Bot().Delete(&tele.StoredMessage{
						ChatID:    con.Chat().ID,
						MessageID: strconv.Itoa(message.ID),
					})
				}, 30, "Ошибка при попытке удалить сообщение с админской командой")
				return nil
			}
			return con.Send(text)
		} else if userState.IsChecking {
			userValue, errUtils := utils.ToInt(con.Text())
			if errUtils != nil {
				log.Printf("Ошибка при попытке преобразовать сообщение пользователя в integer: %v", errUtils)
				deleteMsgAfter(con, utils.GetReturnText(false))
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
				deleteMsgAfter(con, utils.GetReturnText(false))
				return nil
			}
			if !(userValue > -16 && userValue < 16) {
				deleteMsgAfter(con, "Введите значение от -15 до 15")
				return nil
			}
			userState.IsSetting = false
			timeZoneForm := utils.GetTimeZoneForm(userValue)
			msg := fmt.Sprintf("Успешно! Ваш часовой пояс изменён на %s", timeZoneForm)
			deleteMsgAfter(con, msg)
		} else if userState.IsChoosing {
			userValue, err := setUserField(userId, con.Text(), "current_task")
			if err != nil {
				log.Printf("Ошибка при попытке внести значение пользователя в current_task: %v", err)
				msg := utils.GetReturnText(false)
				deleteMsgAfter(con, msg)
				return nil
			}
			if !(userValue > 0 && userValue < 27) {
				deleteMsgAfter(con, "Введите значение от 1 до 26")
				return nil
			}
			userState.IsChoosing = false

			deleteOriginalBotMsg(con)
			msg := "Успешно! Текущее задание: №" + con.Text()
			deleteMsgAfter(con, msg)
		}
		return con.Send("Выберите опцию: ", keyboard.MainMenu())
	}
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
		userMsgID := strconv.Itoa(con.Message().ID)
		message := tele.StoredMessage{ChatID: con.Chat().ID, MessageID: userMsgID}

		deleteMsgAfter(con, "")
		utils.ActionAfter(func() error {
			if message.MessageID == "" {
				return nil
			}
			return con.Bot().Delete(message)
		}, 3, "Ошибка при попытке удалить сообщение пользователя case SpecifyTimeZone")
		userState.IsSetting = true
		text, menu := keyboard.TimeZoneMenu(con)

		return con.Edit(text, menu)
	case "ShowAllExplanations":
		text := userState.Explanations
		return con.Send(text)
	case "ShowPreviousWords", "ShowNextWords":
		text, targetPage, targetSlice, _ := getTargets(con.Message().Text, userState)
		if len(*targetSlice) == 0 {
			deleteMsgAfter(con, utils.GetReturnText(false))
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
