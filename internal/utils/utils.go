package utils

import "strings"

func MaxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func FirstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return value
		}
	}
	return ""
}
