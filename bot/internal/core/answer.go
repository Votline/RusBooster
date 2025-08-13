package core

import (
	"fmt"
	"time"
	"strings"
	"strconv"

	"RusBooster/internal/stat"
	"RusBooster/internal/state"
	"RusBooster/internal/utils"
)

func findAnswer(explanations string, userState *state.UserState, userID int64) string {
	userState.Answer = 0
	userState.OutputTask = ""
	userState.OutputAnswer = ""
	userState.Explanations = explanations
	
	rows := strings.Split(explanations, "@")
	defer state.SetUserState(userState, userID)

	for i, row := range rows {
		cleanRow := strings.TrimSuffix(row, "$")
		lines := strings.Split(cleanRow, "$")

		bases := []rune{}
		for _, line := range lines {
			words := strings.Fields(line)
			if len(words) == 0 {continue}
			bases = append(bases, utils.FindUpperCase(words[0]))
		}
		if utils.AllEqual(bases) && bases[0] != 0 {
			userState.Answer += i+1
			userState.OutputTask += "\n" + rows[i]
			userState.OutputAnswer += strconv.Itoa(i+1)
		}
	}
	return userState.OutputAnswer + userState.OutputTask
}

func SendAnswer(userID int64, userAnswer int, userState *state.UserState) (string, error) {
	userState.IsActive = true
	defer state.SetUserState(userState, userID)

	emoji, err := checkAnswer(userID, userAnswer, userState)
	currentTask, _ := stat.GetSomething(userID, "current_task")
	msg := fmt.Sprintf("%sОтвет на задание №%d: %s\n%s",
		emoji, currentTask, userState.OutputAnswer, userState.OutputTask)
	return msg, err
}

func checkAnswer(userID int64, userAnswer int, userState *state.UserState) (string, error) {
	var emoji string
	srcFields, err := stat.GetExclude(userID, "id", "streak_freeze")
	if err != nil {return "", err}
	userFields := toIntMap(srcFields)

	now := time.Now().UTC().Add(time.Duration(userFields["time_zone"]) * time.Hour)
	currentDate := int(now.Unix()/86400)

	if utils.SumDigits(userAnswer) == userState.Answer {
		if utils.CheckEquals(strconv.Itoa(userAnswer), userState.OutputAnswer){
			userFields["current_score"] += 1
			emoji = "✅"
		} else {
			userFields["current_score"] -= 1
			emoji = "❌"
		}
	} else {
		userFields["current_score"] -= 1
		emoji = "❌"
	}

	if userFields["current_score"] < userFields["worst_task_score"] {
		userFields["worst_task_score"] = userFields["current_score"]
		userFields["worst_task_result"] = userFields["current_task"]
	} else if userFields["current_score"] > userFields["best_task_score"] {
		userFields["best_task_score"] = userFields["current_score"]
		userFields["best_task_result"] = userFields["current_task"]
	}

	if currentDate - userFields["last_active_date"] == 1 || userFields["last_active_date"] == 0 {
		userFields["streak"] += 1
		userFields["last_active_date"] = currentDate
		userState.OutputTask += fmt.Sprintf("\nВаш рекорд теперь: %d🎉", userFields["streak"])
	}

	srcFields = toInterfaceMap(userFields)
	stat.SetExclude(userID, srcFields, "id", "streak_freeze")
	return emoji, nil
}
