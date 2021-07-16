package utils

import "time"

// time stamp sec
func NowSecond() int64  {
	return time.Now().Unix()
}

// time stamp ms
func NowMillisecond() int64  {
	return time.Now().UnixNano() / 1000000
}