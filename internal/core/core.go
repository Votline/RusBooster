package core

import (
	"RusBooster/internal/stat"
	"RusBooster/internal/state"
	"RusBooster/internal/utils"
	"RusBooster/internal/words"
	"fmt"
	"go.uber.org/zap"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

func getWordsForTask(log *zap.Logger, currentTask int) map[string]string {
	wordsForTask := make(map[string]string)
	words.GetAll(log, wordsForTask, currentTask)
	return wordsForTask
}

func MakeUserTask(log *zap.Logger, userId int64, userState *state.UserState) string {
	rand.Seed(time.Now().UnixNano())
	currentTask, err := stat.GetSomething(log, userId, "current_task")
	if err != nil {
		return utils.GetReturnText(false)
	}
	switch currentTask {
	case 9, 10:
		return Number9to12(log, getWordsForTask(log, currentTask), 3, userState, userId)
	case 11:
		return Number9to12(log, getWordsForTask(log, currentTask), 2, userState, userId)
	}
	return ""
}

func SendAnswer(log *zap.Logger, userId int64, userAnswer int, userState *state.UserState) string {
	defer state.SetUserState(log, userState, userId)
	userState.IsActive = true
	emoji := checkAnswer(log, userId, userAnswer, userState)
	currentTask, _ := stat.GetSomething(log, userId, "current_task")
	msg := fmt.Sprintf("%sОтвет на задание: №%d: %s\n%s",
		emoji, currentTask, userState.OutputAnswer, userState.OutputTask)
	return msg
}

func findAnswer(log *zap.Logger, explanations string, userState *state.UserState, userId int64) string {
	userState.Answer = 0
	userState.OutputAnswer = ""
	userState.OutputTask = ""
	userState.Explanations = explanations
	rows := strings.Split(explanations, "@")
	defer state.SetUserState(log, userState, userId)

	for i, row := range rows {
		lines := strings.Split(row, "$")

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
			userState.OutputTask += "\n" + strings.ReplaceAll(rows[i], "$", "")
			userState.OutputAnswer += strconv.Itoa(i + 1)
		}
	}
	return userState.OutputAnswer + userState.OutputTask
}

func checkAnswer(log *zap.Logger, userId int64, userAnswer int, userState *state.UserState) string {
	var emoji string
	currentTask, _ := stat.GetSomething(log, userId, "current_task")
	currentScore, _ := stat.GetSomething(log, userId, "current_score")
	worstTaskScore, _ := stat.GetSomething(log, userId, "worst_task_score")
	worstTaskResult, _ := stat.GetSomething(log, userId, "worst_task_result")
	bestTaskScore, _ := stat.GetSomething(log, userId, "best_task_score")
	bestTaskResult, _ := stat.GetSomething(log, userId, "best_task_result")
	streak , _ := stat.GetSomething(log, userId, "streak")
	timeZone, _ := stat.GetSomething(log, userId, "time_zone")
	lastActiveDate , _ := stat.GetSomething(log, userId, "last_active_date")
	
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
	
	stat.SetSomething(log, userId, currentScore, "current_score")
	stat.SetSomething(log, userId, worstTaskScore, "worst_task_score")
	stat.SetSomething(log, userId, worstTaskResult, "worst_task_result")
	stat.SetSomething(log, userId, bestTaskScore, "best_task_score")
	stat.SetSomething(log, userId, bestTaskResult, "best_task_result")
	stat.SetSomething(log, userId, streak, "streak")
	stat.SetSomething(log, userId, lastActiveDate, "last_active_date")
	return emoji
}

func Number9to12(log *zap.Logger, wordsForTask map[string]string, howMuchWords int, userState *state.UserState, userId int64) string {
	var message strings.Builder
	var explanations strings.Builder

	baseMap := make(map[rune][]string)
	tempAllWords := make([]string, 0, len(wordsForTask))
	for word, explanation := range wordsForTask {
		combined := fmt.Sprintf("%s|%s", word, explanation)
		base := utils.FindUpperCase(explanation)
		if base != 0 {
			baseMap[base] = append(baseMap[base], combined)
		}
		tempAllWords = append(tempAllWords, combined)
	}

	chosenBases := make([]rune, 0, 5)
	for base := range baseMap {
		chosenBases = append(chosenBases, base)
		if len(chosenBases) == 5 {
			break
		}
	}
	if len(chosenBases) < 5 {
		randomBase := chosenBases[rand.Intn(len(chosenBases))]
		chosenBases = append(chosenBases, randomBase)
	}

	rand.Shuffle(len(chosenBases), func(i, j int) {
		chosenBases[i], chosenBases[j] = chosenBases[j], chosenBases[i]
	})

	var selectedWords []string
	for _, base := range chosenBases {
		neededWords := baseMap[base]
		rand.Shuffle(len(neededWords), func(i, j int) {
			neededWords[i], neededWords[j] = neededWords[j], neededWords[i]
		})

		selected := make([]string, 0, len(neededWords))
		selected = append(selected, neededWords...)

		for _, word := range neededWords {
			if !utils.Contains(selected, word) && len(selected) < howMuchWords {
				selected = append(selected, word)
			}
		}

		if len(selected) < howMuchWords {
			for _, word := range tempAllWords {
				if !utils.Contains(selected, word) && !utils.Contains(selectedWords, word) && len(selected) < howMuchWords {
					selected = append(selected, word)
				}
			}
		}
		selectedWords = append(selectedWords, selected[:howMuchWords]...)

	}

	rand.Shuffle(len(selectedWords), func(i, j int) {
		selectedWords[i], selectedWords[j] = selectedWords[j], selectedWords[i]
	})

	message.WriteString("Укажите варианты ответов, в которых во всех словах одного ряда пропущена одна и та же буква. Запишите номера ответов.\n")

	index := 0
	for i := 0; i < 5; i++ {
		message.WriteString(fmt.Sprintf("%d) ", i+1))

		for j := 0; j < howMuchWords; j++ {
			item := selectedWords[index]
			index++

			parts := strings.SplitN(item, "|", 2)
			word, explanation := parts[0], parts[1]

			message.WriteString(word + " ")
			explanations.WriteString(explanation + "\n")
			if j < howMuchWords-1 {
				explanations.WriteString("$")
			}
		}
		message.WriteString("\n")
		explanations.WriteString("@\n")
	}

	findAnswer(log, explanations.String(), userState, userId)
	if len(userState.OutputAnswer) < 2 {
		return Number9to12(log, wordsForTask, howMuchWords, userState, userId)
	}
	return message.String()
}
