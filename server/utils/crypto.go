package utils

import (
	"crypto/md5"
	"fmt"
)

func Md5Salt(Passwd, Salt string) string {
	if Salt == "" {
		Salt = RandString(8)
	}

	value := md5.Sum([]byte(fmt.Sprintf("%s:%s", Salt, Passwd)))

	return fmt.Sprintf("%s:%x", Salt, value)

}
