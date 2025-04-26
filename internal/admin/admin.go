package admin

import (
	"RusBooster/internal/guide"
	"RusBooster/internal/state"
	"RusBooster/internal/words"
	"go.uber.org/zap"
	"strconv"
	"strings"
)

func HandleCommands(log *zap.Logger, userState *state.UserState, userId int64, userMessage string) string {
	parts := strings.SplitN(userMessage, " ", 5)
	taskId, errAtoi := strconv.Atoi(parts[2])
	if errAtoi != nil {
		return "Ошибка при попытке преобразовать номер задания в integer"
	}

	if parts[1] == "add" && len(parts) == 5 {
		if parts[0] == "/word" {
			word := parts[3]
			explanation := parts[4]
			return words.SetSomething(log, word, explanation, taskId)
		} else if parts[0] == "/guide" {
			guideBody := parts[3] + " " + parts[4]
			return guide.AppendGuide(log, taskId, guideBody)
		}
	} else if parts[1] == "del" && len(parts) >= 3 {
		if parts[0] == "/word" {
			word := parts[3]
			return words.DeleteWord(log, taskId, word)
		} else if parts[0] == "/guide" {
			return guide.DeleteGuide(log, taskId)
		}
	} else if parts[1] == "find" && len(parts) >= 4 {
		word := parts[3]
		if len(word) > 1 {
			return words.FindWord(log, taskId, word)
		} else if len(word) == 1 {
		}
	} else if parts[1] == "showall" && len(parts) >= 3 {
		if len(parts) >= 4 {
			return "Найденные слова: \n" + words.ShowAllWords(log, userId, taskId,
				&userState.PartsOfFindWords, &userState.CurrentPageOfFindWords, true, parts[3])
		}
		return "Все слова: \n" + words.ShowAllWords(log, userId, taskId,
			&userState.PartsOfAllWords, &userState.CurrentPageOfAllWords, false, "")
	}
	return "Текущие админ команды:\n[/word|/guide] add taskId word explanation\n[/word|/guide] del taskId (word)\n/word find taskId word\n/word showall taskId (letter)"
}

func ShowPreviousWords(currentPage *int, currentSlice *[]string) string {
	if *currentPage > 0 {
		*currentPage -= 1
	} else {
		*currentPage = len(*currentSlice) - 1
	}
	return (*currentSlice)[*currentPage]
}

func ShowNextWords(currentPage *int, currentSlice *[]string) string {
	if *currentPage < len(*currentSlice)-1 {
		*currentPage += 1
	} else {
		*currentPage = 0
	}
	return (*currentSlice)[*currentPage]
}
