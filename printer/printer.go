package printer

import (
	"fmt"
	"io"
	"time"
)

type SheetPrinter interface {
	PrintSheet(sheet Sheet, buf io.Writer)
}

type SheetPrinterFunc func(sheet Sheet, buf io.Writer)

func (f SheetPrinterFunc) PrintSheet(sheet Sheet, buf io.Writer) {
	f(sheet, buf)
}

type Sheet struct {
	CompanyNo  string
	PersonalNo string
	Employee   string
	Month      string
	TotalHours string
	Rows       []Row
}

type Row struct {
	Day             string
	WeekDay         string
	Start           string
	End             string
	Hours           string
	BreakOneStart   string
	BreakOneEnd     string
	BreakTwoStart   string
	BreakTwoEnd     string
	BreakThreeStart string
	BreakThreeEnd   string
	Comment         string
}

func RowFromSlice(args []string) Row {
	if len(args) != 12 {
		panic("Nee excacly 12 agruments")
	}

	return Row{args[0], args[1], args[2], args[3], args[4], args[5], args[6], args[7], args[8], args[9], args[10], args[11]}
}

func MapWeekDay(day time.Weekday) string {
	switch day {
	case time.Friday:
		return "Fr."
	case time.Monday:
		return "Mo."
	case time.Saturday:
		return "Sa."
	case time.Sunday:
		return "So."
	case time.Thursday:
		return "Do."
	case time.Tuesday:
		return "Di."
	case time.Wednesday:
		return "Mi."
	default:
		panic(fmt.Sprintf("unexpected time.Weekday: %#v", day))
	}
}
