package pkg

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

func Must[T any](value T, err error) T {
	if err != nil {
		log.Panic(err.Error())
	}
	return value
}

type DailyTime struct {
	Hour   int
	Minute int
}

func (d DailyTime) String() string {
	return fmt.Sprintf("%02d:%02d", d.Hour, d.Minute)
}

func NewDailyTime(hour int, minute int) DailyTime {
	if hour < 0 || hour > 24 {
		log.Panicf("Hour %v is nor valid", hour)
	}
	if minute < 0 || minute > 60 {
		log.Panicf("Hour %v is nor valid", minute)
	}

	return DailyTime{hour, minute}
}

type DailyPeriod struct {
	Start DailyTime
	End   DailyTime
}

func (d DailyPeriod) Minutes() int {
	return ((d.End.Hour - d.Start.Hour) * 60) + (d.End.Minute - d.Start.Minute)
}

func DailyPeriodFromString(start string, end string) DailyPeriod {
	startParts := strings.Split(start, ":")
	endParts := strings.Split(end, ":")

	return DailyPeriod{
		NewDailyTime(
			Must(strconv.Atoi(startParts[0])),
			Must(strconv.Atoi(startParts[1])),
		),
		NewDailyTime(
			Must(strconv.Atoi(endParts[0])),
			Must(strconv.Atoi(endParts[1])),
		),
	}
}

type DayType int

const (
	DayTypeWeekEnd = 1
	DayTypeWeekday = 2
	DayTypeFeast   = 3
	DayTypeHoliday = 4
)

type Day struct {
	Date     time.Time
	Type     DayType
	Duration DailyPeriod
	Breaks   []DailyPeriod
}

func (w Day) Hours() float64 {
	result := w.Duration.Minutes()

	for _, br := range w.Breaks {
		result -= br.Minutes()
	}

	return float64(result) / 60
}
