package admin

import (
	"fmt"
	"strings"
	"strconv"

	"RusBooster/internal/state"
	"RusBooster/internal/words"
	"RusBooster/internal/guide"
)

func HandleCommands(userState *state.UserState, userId int64, userMessage string) string {
	parts := strings.Fields(userMessage)
	if len(parts) < 2 {
		return "Недостаточно аргументов"
	}

	command := parts[1]
	switch command {
	case "addword":
		if len(parts) < 5 {
			return "Использование: /admin addword <taskId> <word> <explanation>"
		}
		taskId, err := strconv.Atoi(parts[2])
		if err != nil {
			return "Неверный taskId"
		}
		word := parts[3]
		explanation := strings.Join(parts[4:], " ")
		return words.SetSomething(word, explanation, taskId)

	case "delword":
		if len(parts) < 4 {
			return "Использование: /admin delword <taskId> <word>"
		}
		taskId, err := strconv.Atoi(parts[2])
		if err != nil {
			return "Неверный taskId"
		}
		word := parts[3]
		return words.DeleteWord(taskId, word)

	case "addguide":
		if len(parts) < 4 {
			return "Использование: /admin addguide <taskId> <guide_text>"
		}
		taskId, err := strconv.Atoi(parts[2])
		if err != nil {
			return "Неверный taskId"
		}
		guideText := strings.Join(parts[3:], " ")
		return guide.AppendGuide(taskId, guideText)

	case "delguide":
		if len(parts) < 3 {
			return "Использование: /admin delguide <taskId>"
		}
		taskId, err := strconv.Atoi(parts[2])
		if err != nil {
			return "Неверный taskId"
		}
		return guide.DeleteGuide(taskId)

	case "showall":
		if len(parts) < 3 {
			return "Использование: /admin showall <taskId>"
		}
		_, err := strconv.Atoi(parts[2])
		if err != nil {
			return "Неверный taskId"
		}
		return fmt.Sprintf("Все слова: \n")

	default:
		return "Неизвестная команда"
	}
}

func ShowPreviousWords(targetPage *int, targetSlice *[]string) string {
	if *targetPage > 0 {
		*targetPage--
	}
	return (*targetSlice)[*targetPage]
}

func ShowNextWords(targetPage *int, targetSlice *[]string) string {
	if *targetPage < len(*targetSlice)-1 {
		*targetPage++
	}
	return (*targetSlice)[*targetPage]
}
