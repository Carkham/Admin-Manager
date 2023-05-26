package utils

import (
	"admin/model"
	"fmt"
	"strconv"
)

func ParseCronExpr(cronConf *model.CronJobConfig) string {
	var minStr, hourStr, weekdayStr string
	if cronConf == nil {
		return "*  *  *  *  *"
	}
	if cronConf.TrigMin < 0 {
		minStr = "*/1"
	} else {
		hourStr = strconv.Itoa(cronConf.TrigMin)
	}

	if cronConf.TrigHour < 0 {
		hourStr = "*"
	} else {
		hourStr = strconv.Itoa(cronConf.TrigHour)
	}

	if cronConf.TrigWeekday < 0 {
		weekdayStr = "*"
	} else {
		weekdayStr = strconv.Itoa(cronConf.TrigWeekday)
	}

	return fmt.Sprintf(
		"%s  %s  *  *  %s",
		minStr,
		hourStr,
		weekdayStr,
	)
}
