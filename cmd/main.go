package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/apfelfrisch/timesheet/pkg"
)

func main() {
	year := flag.Int("year", time.Now().Year(), "Year (e.g. 2025)")
	month := flag.Int("month", int(time.Now().Month()), "Month (1-12)")
	flag.Parse()

	var pdf bytes.Buffer
	nils := pkg.NewEmployeeFromConfig(pkg.NewConfigFromFile())
	nils.PrintSheet(*year, time.Month(*month), &pdf)

	os.WriteFile(fmt.Sprintf("arbeitszeiten-%v-%v.pdf", *year, *month), pdf.Bytes(), 0644)
}
