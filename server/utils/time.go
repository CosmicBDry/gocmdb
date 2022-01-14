package utils

import (
	"strings"
	"time"
)

func TimeParse(date string) *time.Time {

	Time, _ := time.Parse("2006-01-02", date)

	return &Time

}

func TimeStringSplit(date string) string {

	return strings.SplitN(date, " ", 2)[0]

}
