package core

import (
	"log"

	"RusBooster/internal/utils"
)

func toIntMap(src map[string]interface{}) map[string]int {
	dst := make(map[string]int, len(src))
	
	for k, v := range src {
		switch val := v.(type) {
		case int:
			dst[k] = val
		case int64:
			dst[k] = int(val)
		case float64:
			dst[k] = int(val)
		default:
			log.Printf("Non-numeric value: %v|%v", v, val)
			dst[k] = 0
		}
	}

	return dst
}

func toInterfaceMap(src map[string]int) map[string]interface{} {
	dst := make(map[string]interface{}, len(src))

	for k, v := range src {
		dst[k] = v
	}

	return dst
}

func allEqual(letters []rune) bool {
	if len(letters) == 0 {
		return false
	}
	for _, r := range letters {
		if r != letters[0] {
			return false
		}
	}
	return true
}

func checkEquals(correct, verifiable string) bool {
	correct = utils.GetDigitHash(correct)
	verifiable = utils.GetDigitHash(verifiable)
	return correct == verifiable
}

func sumDigits(number int) int {
	sum := 0
	for number > 0 {
		digit := number % 10
		sum += digit
		number /= 10
	}
	return sum
}
