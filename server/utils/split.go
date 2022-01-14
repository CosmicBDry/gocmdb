package utils

import "strings"

func SplitMd5Salt(text string) (string, string) {
	result := strings.SplitN(text, ":", 2)
	return result[0], result[1]
}
