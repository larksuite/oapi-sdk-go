package tools

import (
	"strconv"
	"time"
)

func TsToTime(ts string) (time.Time, error) {
	_ts, err := strconv.ParseFloat(ts, 64)
	if err != nil {
		return time.Time{}, err
	}
	ns := int64(_ts * float64(time.Second))
	return time.Unix(0, ns), nil
}
