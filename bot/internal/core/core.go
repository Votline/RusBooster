package core

import (
	"time"
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
	case 11, 12:
		return Number9to12(getWordsForTask(currentTask), 2, userState, userId)
	}
	return ""
}
