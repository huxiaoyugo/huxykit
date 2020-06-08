package utils

import "time"

// 获取传入时间与当日凌晨的时间间隔
func GetIntervalFromDawn(tm time.Time) time.Duration {
	dawn := time.Date(tm.Year(),tm.Month(), tm.Day(),0,0,0,0,time.Local)
	return tm.Sub(dawn)
}