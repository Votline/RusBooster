package core

import "log"

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
