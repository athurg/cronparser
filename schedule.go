package main

import (
	"fmt"
	"time"
)

var allHours = []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23}
var allMinutes = []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 58, 59}

type Schedule struct {
	Command string

	Minutes  []int  //限制可用的分，为空时表示无限制
	Hours    []int  //限制可用的时，为空时表示无限制
	Days     []int  //限制可用的日，为空时表示无限制
	Months   []int  //限制可用的月，为空时表示无限制
	Weekdays []int  //限制可用的星期，为空时表示无限制
}

func (schedule Schedule) ScheduleMap(from, to time.Time) map[int]map[int]string {
	var hours []int
	var minutes []int

	if from.IsZero() || to.IsZero() {
		//不考虑日期
		hours = schedule.Hours
		minutes = schedule.Minutes
	} else {
		//指定日期范围
		toDate := time.Date(to.Year(), to.Month(), to.Day(), 0, 0, 0, 0, time.UTC)
		fromDate := time.Date(from.Year(), from.Month(), from.Day(), 0, 0, 0, 0, time.UTC)

		for date := fromDate; !date.After(toDate); date = date.AddDate(0, 0, 1) {
			if len(schedule.Weekdays) != 0 && !findInSlices(schedule.Weekdays, int(date.Weekday())) {
				continue
			}
			if len(schedule.Months) != 0 && !findInSlices(schedule.Months, int(date.Month())) {
				continue
			}

			if len(schedule.Days) != 0 && !findInSlices(schedule.Days, date.Day()) {
				continue
			}

			hours = schedule.Hours
			minutes = schedule.Minutes
			break
		}
	}

	//没有符合的计划任务
	if hours == nil || minutes == nil {
		return nil
	}

	if len(hours) == 0 {
		hours = allHours
	}

	if len(minutes) == 0 {
		minutes = allMinutes
	}

	data := map[int]map[int]string{}
	for _, hour := range hours {
		if data[hour] == nil {
			data[hour] = map[int]string{}
		}
		for _, minute := range minutes {
			data[hour][minute] = schedule.Command
		}
	}
	return data
}

func findInSlices(slice []int, value int) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}

	return false
}

//转换[]Schedule为指定时间段的小时分钟的调度表
func schedulesToMinutesCommandMap(schedules []Schedule, from, to time.Time) (map[int]map[int][]string, []string) {
	//初始化调度命令清单
	commands := make([]string, 0)

	//初始化一个24hourx60min的空白调度表
	scheduleMap := make(map[int]map[int][]string)
	for hour := 0; hour < 24; hour++ {
		hourMap := map[int][]string{}
		for minute := 0; minute < 60; minute++ {
			hourMap[minute] = nil
		}
		scheduleMap[hour] = hourMap
	}

	//将调度任务填入调度表、命令清单
	for _, schedule := range schedules {
		hourSchedules := schedule.ScheduleMap(from, to)
		for hour, minuteSchedules := range hourSchedules {
			for minute, command := range minuteSchedules {
				commands = append(commands, fmt.Sprintf("%02d:%02d: %s", hour, minute, command))
				scheduleMap[hour][minute] = append(scheduleMap[hour][minute], wheneverCommandTrim(command))
			}
		}
	}

	return scheduleMap, commands
}
