package bot

import (
	"log"
	"strings"

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

func containsRune(text string, item string) bool {
	for _, r := range text {
		if string(r) == item {
			return true
		}
	}
	return false
}
