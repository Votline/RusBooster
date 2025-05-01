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
	var rest string
	parts := strings.SplitN(userMessage, " ", 4)
	if len(parts) < 3 {
		return "Неверный формат сообщения, parts < 3"
	} else if len(parts) > 3{
		rest = parts[3]
	}

	command := parts[0]
	action := parts[1]
	content := strings.SplitN(rest, " - ", 2)

	taskId, errAtoi := strconv.Atoi(parts[2])
	if errAtoi != nil {
		return "Ошибка при попытке преобразовать номер задания в integer"
	}
	
	if action == "add" && len(parts) == 4 {
		if command == "/word" {
			if len(content) <= 1 {
				return "Неверный формат сообщения, content <= 1"
			}
			word := content[0]
			explanation := content[1]
			return words.SetSomething(log, word, explanation, taskId)
		} else if command == "/guide" {
			return guide.AppendGuide(log, taskId, rest)
		}
	} else if action == "del" && len(parts) >= 4 {
		if command == "/word" {
			word := content[0]
			return words.DeleteWord(log, taskId, word)
		} else if parts[0] == "/guide" {
			return guide.DeleteGuide(log, taskId)
		}
	} else if action == "find" && len(parts) >= 4 {
		word := content[0]
		if len(word) > 1 {
			text, _ := words.FindWord(log, taskId, word)
			return text
		} else if len(word) == 1 {
		}
	} else if action == "showall" && len(parts) >= 3 {
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
