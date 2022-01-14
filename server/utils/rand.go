package utils

import (
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())

}
func RandString(PassLen int) string {
	letters := "AaBbCcDdEeFfGgHhIiJjKkLlMmNnOoPpQqRrSsTtUuVvWwXxYyZz0123456789*%#"

	length := len(letters)
	attr := make([]byte, PassLen)
	for i := 0; i < PassLen; i++ {

		//attr = append(attr, letters[rand.Int()%length])
		attr[i] = letters[rand.Int()%length]
	}

	return string(attr)
}
