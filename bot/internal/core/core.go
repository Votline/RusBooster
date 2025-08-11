package core

import (
	"fmt"
	"time"
	"strings"
	"strconv"
	"math/rand"

	"RusBooster/internal/stat"
	"RusBooster/internal/state"
	"RusBooster/internal/utils"
	"RusBooster/internal/words"
)

func getWordsForTask(currentTask int) map[string]string {
	wordsForTask := make(map[string]string)
	words.GetAll(wordsForTask, currentTask)
	return wordsForTask
}

func MakeUserTask(userId int64, userState *state.UserState) string {
	rand.Seed(time.Now().UnixNano())
	currentTask, err := stat.GetSomething(userId, "current_task")
	if err != nil {
		return utils.GetReturnText(false)
	}
	switch currentTask {
	case 9, 10:
		return Number9to12(getWordsForTask(currentTask), 3, userState, userId)
	case 11:
		return Number9to12(getWordsForTask(currentTask), 2, userState, userId)
	}
	return ""
}

func SendAnswer(userId int64, userAnswer int, userState *state.UserState) string {
	defer state.SetUserState(userState, userId)
	userState.IsActive = true
	emoji := checkAnswer(userId, userAnswer, userState)
	currentTask, _ := stat.GetSomething(userId, "current_task")
	msg := fmt.Sprintf("%sОтвет на задание: №%d: %s\n%s",
		emoji, currentTask, userState.OutputAnswer, userState.OutputTask)
	return msg
}

func findAnswer(explanations string, userState *state.UserState, userId int64) string {
	userState.Answer = 0
	userState.OutputAnswer = ""
	userState.OutputTask = ""
	userState.Explanations = explanations
	rows := strings.Split(explanations, "@")
	defer state.SetUserState(userState, userId)

	for i, row := range rows {
		cleanRow := strings.TrimSuffix(row, "$")
		lines := strings.Split(cleanRow, "$")

		bases := []rune{}
		for _, line := range lines {
			word := strings.Fields(line)
			if len(word) == 0 {
				continue
			}
			bases = append(bases, utils.FindUpperCase(word[0]))
		}
		if utils.AllEqual(bases) && bases[0] != 0 {
			userState.Answer += i + 1
			userState.OutputTask += "\n" + rows[i]
			userState.OutputAnswer += strconv.Itoa(i + 1)
		}
	}
	return userState.OutputAnswer + userState.OutputTask
}

func checkAnswer(userId int64, userAnswer int, userState *state.UserState) string {
	var emoji string
	currentTask, _ := stat.GetSomething(userId, "current_task")
	currentScore, _ := stat.GetSomething(userId, "current_score")
	worstTaskScore, _ := stat.GetSomething(userId, "worst_task_score")
	worstTaskResult, _ := stat.GetSomething(userId, "worst_task_result")
	bestTaskScore, _ := stat.GetSomething(userId, "best_task_score")
	bestTaskResult, _ := stat.GetSomething(userId, "best_task_result")
	streak , _ := stat.GetSomething(userId, "streak")
	timeZone, _ := stat.GetSomething(userId, "time_zone")
	lastActiveDate , _ := stat.GetSomething(userId, "last_active_date")
	
	now := time.Now().UTC().Add(time.Duration(timeZone) * time.Hour)
	currentDate := int(now.Unix() / 86400)

	if utils.SumDigits(userAnswer) == userState.Answer {
		if utils.CheckEquals(strconv.Itoa(userAnswer), userState.OutputAnswer) {
			currentScore += 1
			emoji = "✅"
		} else {
			currentScore -= 1
			emoji = "❌"
		}
	} else {
		currentScore -= 1
		emoji = "❌"
	}

	if currentScore < worstTaskScore {
		worstTaskScore = currentScore
		worstTaskResult = currentTask
	} else if currentScore > bestTaskScore {
		bestTaskScore = currentScore
		bestTaskResult = currentTask
	}

	if currentDate - lastActiveDate == 1 || lastActiveDate == 0{
		streak += 1
		lastActiveDate = currentDate
		userState.OutputTask += fmt.Sprintf("\nВаш рекорд теперь: %d🎉", streak)
	}
	
	stat.SetSomething(userId, currentScore, "current_score")
	stat.SetSomething(userId, worstTaskScore, "worst_task_score")
	stat.SetSomething(userId, worstTaskResult, "worst_task_result")
	stat.SetSomething(userId, bestTaskScore, "best_task_score")
	stat.SetSomething(userId, bestTaskResult, "best_task_result")
	stat.SetSomething(userId, streak, "streak")
	stat.SetSomething(userId, lastActiveDate, "last_active_date")
	return emoji
}

func Number9to12(wordsForTask map[string]string, howMuchWords int, userState *state.UserState, userId int64) string {
	var message strings.Builder
	var explanations strings.Builder

	words := make([]string, 0, len(wordsForTask))
	for word := range wordsForTask {
		words = append(words, word)
	}

	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(words), func(i, j int) {
		words[i], words[j] = words[j], words[i]
	})

	for i := 0; i < howMuchWords && i < len(words); i++ {
		word := words[i]
		explanation := wordsForTask[word]
		message.WriteString(fmt.Sprintf("%d. %s\n", i+1, word))
		explanations.WriteString(fmt.Sprintf("%s$", explanation))
		if i < howMuchWords-1 && i < len(words)-1 {
			explanations.WriteString("@")
		}
	}

	findAnswer(explanations.String(), userState, userId)
	if len(userState.OutputAnswer) < 2 {
		return Number9to12(wordsForTask, howMuchWords, userState, userId)
	}
	return message.String()
}
