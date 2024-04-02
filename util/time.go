package util

import (
	"time"
)

func StringToTime(timeStr string) (time.Time, error) {
	layout := "2006-01-02 15:04:05 -0700 MST"
	t, err := time.Parse(layout, timeStr)

	if err != nil {
		return time.Time{}, err
	}
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		return time.Time{}, err
	}

	return t.In(loc), nil
}
