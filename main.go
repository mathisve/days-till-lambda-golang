package main

import (
	"bytes"
	"context"
	"embed"
	"fmt"
	"html/template"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
)

var d time.Time

//go:embed file.html
var f embed.FS

type Till struct {
	HasPassed bool
	Days      int
	Hours     int
	Minutes   int
	Seconds   int
}

type Template struct {
	Header string
	Time   string
	Footer string
}

type Output struct {
	StatusCode int               `json:"statusCode"`
	Headers    map[string]string `json:"headers"`
	Body       string            `json:"body"`
}

func init() {
	year, _ := strconv.Atoi(os.Getenv("year"))
	month, _ := strconv.Atoi(os.Getenv("month"))
	day, _ := strconv.Atoi(os.Getenv("day"))
	hour, _ := strconv.Atoi(os.Getenv("hour"))
	minute, _ := strconv.Atoi(os.Getenv("minute"))

	d = date(year, month, day, hour, minute)
}

func main() {
	// ctx := context.Background()
	// output, err := Handler(ctx)
	// if err != nil {
	// 	log.Println(err)
	// }

	// fmt.Println(output)
	lambda.Start(Handler)
}

func Handler(ctx context.Context) (output Output, err error) {
	buf := new(bytes.Buffer)

	text := makeTemplate(calcUntil(d))

	tmpl, err := template.ParseFS(f, "file.html")
	tmpl.Execute(buf, text)

	if err != nil {
		log.Println(err)

		output = Output{
			StatusCode: 400,
		}
	} else {
		output = Output{
			StatusCode: 200,
			Headers: map[string]string{
				"Content-Type": "text/html",
			},
			Body: buf.String(),
		}
	}

	return output, nil
}

func date(year, month, day, hour, minute int) time.Time {
	return time.Date(year, time.Month(month), day, hour, minute, 0, 0, time.Local)
}

func calcUntil(d time.Time) Till {
	t := int(time.Until(d).Seconds())

	return Till{
		HasPassed: t < 0,
		Days:      pos(t / 86400),
		Hours:     pos(t % 86400 / 3600),
		Minutes:   pos(t % 86400 % 3600 / 60),
		Seconds:   pos(t % 86400 % 3600 % 60),
	}
}

func makeTemplate(t Till) Template {
	var tmpl Template

	if t.HasPassed {
		tmpl.Header = "Time since bliss:"
		tmpl.Footer = "Time since my plane landed!"
	} else {
		tmpl.Header = "Time till bliss:"
		tmpl.Footer = "Time until my plane lands!"
	}

	if t.Days == 1 {
		tmpl.Time += fmt.Sprintf("%d day ", t.Days)
	} else {
		tmpl.Time += fmt.Sprintf("%d days ", t.Days)
	}

	if t.Hours == 1 {
		tmpl.Time += fmt.Sprintf("%d hour ", t.Hours)
	} else {
		tmpl.Time += fmt.Sprintf("%d hours ", t.Hours)
	}

	if t.Minutes == 1 {
		tmpl.Time += fmt.Sprintf("%d minute ", t.Minutes)
	} else {
		tmpl.Time += fmt.Sprintf("%d minutes ", t.Minutes)
	}

	if t.Seconds == 1 {
		tmpl.Time += fmt.Sprintf("%d second ", t.Seconds)
	} else {
		tmpl.Time += fmt.Sprintf("%d seconds ", t.Seconds)
	}

	return tmpl
}

func pos(i int) int {
	if i < 0 {
		return i * -1
	} else {
		return i
	}
}
