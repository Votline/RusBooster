package bot

import (
	"log"
	"strings"
	"strconv"

	tele "gopkg.in/telebot.v3"

	"RusBooster/internal/stat"
	"RusBooster/internal/state"
	"RusBooster/internal/utils"
	"RusBooster/internal/keyboard"
)

func getTargets(message string, userState *state.UserState) (string, *int, *[]string, *tele.ReplyMarkup) {
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

func setUserField(userId int64, message string, column string) (int, error) {
	userNumber, errUtils := utils.ToInt(message)
	if errUtils != nil {
		return -1, errUtils
	}
	if err := stat.SetSomething(userId, userNumber, column); err != nil {
		log.Printf("Ошибка при попытке установить значение пользователя: %v", err)
		return -1, err
	}
	return userNumber, nil
}

func deleteOriginalBotMsg(con tele.Context) {
	originalBotMsg := &tele.StoredMessage{
		ChatID:    con.Chat().ID,
		MessageID: strconv.Itoa(con.Message().ID - 1)}
	utils.ActionAfter(func() error {
		if originalBotMsg.MessageID == "" {
			return nil
		}
		return con.Bot().Delete(originalBotMsg)
	}, 3,
		"Ошибка при попытке удалить первоначальное сообщение бота")
}

func deleteMsgAfter(con tele.Context, temporaryText string) {
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
	utils.ActionAfter(func() error {
		if deleteMessage.MessageID == "" {
			return nil
		}
		return con.Bot().Delete(deleteMessage)
	}, 3,
		"Ошибка при попытке удалить ответное сообщение бота")
	if anotherMessage.MessageID != "" {
		utils.ActionAfter(func() error {
			if anotherMessage.MessageID == "" {
				return nil
			}
			return con.Bot().Delete(anotherMessage)
		}, 3,
			"Ошибка при попытке удалить сообщение пользователя")
	}
}

func containsRune(text string, item string) bool {
	for _, r := range text {
		if string(r) == item {
			return true
		}
	}
	return false
}
