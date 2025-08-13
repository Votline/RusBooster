package core

import (
	"fmt"
	"strings"
	"math/rand"

	"RusBooster/internal/state"
	"RusBooster/internal/utils"
)

type WordItem struct {
	Word string
	Explanation string
	Base rune
}

func Number9to12(wordsForTask map[string]string, howMuchWords int, userState *state.UserState, userID int64) string {
	items := parseWords(wordsForTask)
	baseMap := groupByBase(items)

	chosenBases := chooseBases(baseMap, 5)
	selectedWords := selectWords(baseMap, chosenBases, howMuchWords, items)

	message, explanations := buildTask(selectedWords, howMuchWords)
	
	findAnswer(explanations, userState, userID)
	if len(userState.OutputAnswer) < 2 {
		return Number9to12(wordsForTask, howMuchWords, userState, userID)
	}
	
	return message
}

func parseWords(words map[string]string) []WordItem {
	result := make([]WordItem, 0, len(words))
	for word, explanation := range words {
		base := utils.FindUpperCase(explanation)
		result = append(result, WordItem{
			Word: word,
			Explanation: explanation,
			Base: base,
		})
	}
	return result
}

func groupByBase(items []WordItem) map[rune][]WordItem {
	baseMap := make(map[rune][]WordItem)
	for _, item := range items {
		if item.Base != 0 {
			baseMap[item.Base] = append(baseMap[item.Base], item)
		}
	}
	return baseMap
}

func chooseBases(baseMap map[rune][]WordItem, maxBases int) []rune {
	chosen := make([]rune, 0, maxBases)
	for base := range baseMap {
		chosen = append(chosen, base)
		if len(chosen) == maxBases {
			break
		}
	}
	for len(chosen) < maxBases {
		randBase := chosen[rand.Intn(len(chosen))]
		chosen = append(chosen, randBase)
	}

	rand.Shuffle(len(chosen), func(i, j int){
		chosen[i], chosen[j] = chosen[j], chosen[i]})
	return chosen
}

func selectWords(baseMap map[rune][]WordItem, bases []rune, countPerBase int, allItems []WordItem) []WordItem {
	selected := make([]WordItem, 0, len(bases)*countPerBase)

	for _, base := range bases {
		needed := append([]WordItem(nil), baseMap[base]...)
		rand.Shuffle(len(needed), func(i, j int){
			needed[i], needed[j] = needed[j], needed[i]})

		if len(needed) < countPerBase {
			for _, item := range allItems {
				if !containsWord(selected, item) && !containsWord(needed, item){
					needed = append(needed, item)
					if len(needed) >= countPerBase {
						break
					}
				}
			}
		}

		selected = append(selected, needed[:countPerBase]...)
	}

	rand.Shuffle(len(selected), func(i, j int){
		selected[i], selected[j] = selected[j], selected[i]})
	return selected
}
func containsWord(list []WordItem, word WordItem) bool {
	for _, item := range list {
		if item.Word == word.Word && item.Explanation == word.Explanation {
			return true
		}
	}
	return false
}

func buildTask(selected []WordItem, howMuchWords int) (string, string) {
	var msg strings.Builder
	var exp strings.Builder

	msg.WriteString("Укажите варианты ответов, в которых во всех словах одного ряда пропущена одна и та же буква. Запишите номера ответов.\n")

	index := 0
	for i := 0; i < 5; i++ {
		msg.WriteString(fmt.Sprintf("%d) ", i+1))
		for j := 0; j < howMuchWords; j++ {
			item := selected[index]
			index++
			msg.WriteString(item.Word + " ")
			exp.WriteString(item.Explanation + "\n")
			if j < howMuchWords-1 {
				exp.WriteString("$")
			}
		}
		msg.WriteString("\n")
		exp.WriteString("@\n")
	}

	return msg.String(), exp.String()
}
