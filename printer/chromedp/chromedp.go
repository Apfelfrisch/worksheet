package chromedp

import (
	"context"
	"embed"
	"html/template"
	"io"
	"log"
	"net/url"
	"strings"

	"github.com/apfelfrisch/timesheet/printer"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

//go:embed template.html
var tmplFS embed.FS

func Printer() printer.SheetPrinter {
	return printer.SheetPrinterFunc(PrintSheet)
}

func PrintSheet(sheet printer.Sheet, buf io.Writer) {
	tmpl, err := template.ParseFS(tmplFS, "template.html")
	if err != nil {
		log.Panic(err)
	}

	var htmlBuffer strings.Builder
	tmpl.Execute(&htmlBuffer, sheet)

	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	if err := chromedp.Run(ctx, printToPdfTask(htmlBuffer.String(), buf)); err != nil {
		log.Fatal(err)
	}
}

func printToPdfTask(html string, writer io.Writer) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate("data:text/html," + url.PathEscape(html)),
		chromedp.ActionFunc(func(ctx context.Context) error {
			buf, _, err := page.PrintToPDF().WithPrintBackground(true).Do(ctx)
			if err != nil {
				return err
			}
			writer.Write(buf)
			return nil
		}),
	}
}
