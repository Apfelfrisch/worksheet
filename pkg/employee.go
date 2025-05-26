package pkg

import (
	"io"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/apfelfrisch/timesheet/printer"
	"github.com/apfelfrisch/timesheet/printer/chromedp"

	"github.com/wlbr/feiertage"
)

type WorkMonth struct {
	WorkDays []Day
}

func (w WorkMonth) TotalHours() float64 {
	var totalHours float64
	for _, day := range w.WorkDays {
		totalHours += day.Hours()
	}
	return totalHours
}

type Employee struct {
	FirstName  string
	LastName   string
	CompanyNo  string
	PersonalNo string
	Workday    DailyPeriod
	Breaks     []DailyPeriod
	Workdays   []time.Weekday
	Holidays   map[time.Time]struct{}
	Printer    printer.SheetPrinter
}

func (e Employee) PrintSheet(year int, month time.Month, buf io.Writer) {
	rows := make([]printer.Row, 0, 28)

	workMonth := e.DaysInMonth(year, month)
	for _, day := range workMonth.WorkDays {
		row := make([]string, 12)

		row[0] = day.Date.Format("02")
		row[1] = printer.MapWeekDay(day.Date.Weekday())

		if day.Hours() > 0.0 {
			row[2] = day.Duration.Start.String()
			row[3] = day.Duration.End.String()
			row[4] = strconv.FormatFloat(day.Hours(), 'f', 2, 64)
		}

		if day.Type == DayTypeWeekday {
			for i, b := range day.Breaks {
				row[i+5] = b.Start.String()
				row[i+6] = b.End.String()
			}
		} else if day.Type == DayTypeHoliday {
			row[11] = "Urlaub"
		} else if day.Type == DayTypeFeast {
			row[11] = holidayMap(day.Date.Year())[day.Date]
		}

		rows = append(rows, printer.RowFromSlice(row))
	}

	e.Printer.PrintSheet(printer.Sheet{
		CompanyNo:  e.CompanyNo,
		PersonalNo: e.PersonalNo,
		Employee:   e.LastName + ", " + e.FirstName,
		Month:      printer.MapMonth(month) + " " + strconv.Itoa(year),
		TotalHours: strconv.FormatFloat(workMonth.TotalHours(), 'f', 2, 64),
		Rows:       rows,
	}, buf)
}

func (e Employee) DaysInMonth(year int, month time.Month) WorkMonth {
	feastdays := holidayMap(year)
	start := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
	end := start.AddDate(0, 1, 0)

	var days []Day
	for d := start; d.Before(end); d = d.AddDate(0, 0, 1) {
		if !slices.Contains(e.Workdays, d.Weekday()) {
			days = append(days, Day{Date: d, Type: DayTypeWeekEnd})
			continue
		}

		midnight := d.Truncate(24 * time.Hour)

		var dType DayType
		if _, isFeast := feastdays[midnight]; isFeast {
			dType = DayTypeFeast
		} else if _, isHoliday := e.Holidays[midnight]; isHoliday {
			dType = DayTypeHoliday
		} else {
			dType = DayTypeWeekday
		}

		for _, workBreak := range e.Breaks {
			days = append(days, Day{
				Date:     d,
				Type:     dType,
				Duration: e.Workday,
				Breaks:   []DailyPeriod{workBreak},
			})
		}
	}

	return WorkMonth{days}
}

func holidayMap(year int) map[time.Time]string {
	holidayMap := map[time.Time]string{}
	for _, h := range feiertage.Niedersachsen(year).Feiertage {
		date := h.Time.Truncate(24 * time.Hour)
		holidayMap[date] = h.Text
	}
	return holidayMap
}

func NewEmployee(workday string, lunchbreak string, printer printer.SheetPrinter) Employee {
	workdayParts := strings.Split(workday, "-")
	lunchbreakParts := strings.Split(lunchbreak, "-")

	return Employee{
		Workday: DailyPeriodFromString(workdayParts[0], workdayParts[1]),
		Breaks: []DailyPeriod{
			DailyPeriodFromString(lunchbreakParts[0], lunchbreakParts[1]),
		},
		Workdays: []time.Weekday{time.Monday, time.Tuesday, time.Wednesday, time.Thursday, time.Friday},
		Printer:  printer,
	}
}

func NewEmployeeFromConfig(config Config) Employee {
	var dayBreaks []DailyPeriod
	for _, b := range config.Workday.Breaks {
		dayBreaks = append(dayBreaks, DailyPeriodFromString(b.Start, b.End))
	}

	holidaysMap := make(map[time.Time]struct{})
	for _, holiday := range config.Holidays {
		fromDate := Must(time.Parse("2006-01-02", holiday.From))
		holidaysMap[fromDate] = struct{}{}
		if holiday.Until != "" {
			untilDate := Must(time.Parse("2006-01-02", holiday.Until))
			for d := fromDate; d.Before(untilDate) || d.Equal(untilDate); d = d.AddDate(0, 0, 1) {
				holidaysMap[d] = struct{}{}
			}
		}
	}

	return Employee{
		FirstName:  config.FirstName,
		LastName:   config.LastName,
		CompanyNo:  config.CompanyNo,
		PersonalNo: config.PersonalNo,
		Workday:    DailyPeriodFromString(config.Workday.Start, config.Workday.End),
		Breaks:     dayBreaks,
		Workdays:   []time.Weekday{time.Monday, time.Tuesday, time.Wednesday, time.Thursday, time.Friday},
		Holidays:   holidaysMap,
		Printer:    chromedp.Printer(),
	}
}
