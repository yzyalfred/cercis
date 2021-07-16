package utils

import "time"

func NowSecond() int64  {
	return time.Now().Unix()
}

func NowMillisecond() int64  {
	return time.Now().UnixNano() / 1000000
}