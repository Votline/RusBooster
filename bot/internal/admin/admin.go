package admin

import (
	"log"
	"strings"
	"strconv"

	"RusBooster/internal/state"
	"RusBooster/internal/words"
	"RusBooster/internal/guide"
)

func HandleCommands(userState *state.UserState, userID int64, userMessage string) string {
	content := make([]string, 0)
	parts := strings.SplitN(userMessage, " ", 3)
	if len(parts) < 2 {
		return "Недостаточно аргументов"
	} else if len(parts) > 2 {
		content = strings.SplitN(parts[2], "$", 2)
	}
	
	command := parts[0]
	taskID, err := strconv.Atoi(parts[1])
	if err != nil {return "Неверный taskID"}

	switch command {
	case "/addword":
		if len(parts) < 4 {
			return "Использование: /addword <taskID> <word> <explanation>"
		}
		return words.SetSomething(content[0], content[1], taskID)

	case "/delword":
		if len(parts) < 3 {
			return "Использование: /delword <taskID> <word>"
		}
		return words.DeleteWord(taskID, content[0])

	case "/addguide":
		if len(parts) < 3 {
			return "Использование: /addguide <taskID> <guide_text>"
		}
		log.Println(content[0])
		return guide.AppendGuide(taskID, content[0])

	case "/delguide":
		if len(parts) < 2 {
			return "Использование: /delguide <taskID>"
		}
		return guide.DeleteGuide(taskID)

	case "/showall":
		if len(parts) < 2 {
			return "Использование: /showall <taskID>"
		} else if len(parts) > 2 {
			return "Найденные слова: \n" + words.ShowAllWords(
				userID, taskID, &userState.PartsOfFindWords,
				&userState.CurrentPageOfFindWords, true, content[0])
		}
		return "Все слова: \n" + words.ShowAllWords(
			userID, taskID, &userState.PartsOfAllWords,
			&userState.CurrentPageOfAllWords, false, "")

	default:
		log.Printf("Неизвестная команда: %s", command)
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
