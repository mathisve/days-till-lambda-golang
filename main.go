package main

import (
	"bytes"
	"context"
	"embed"
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
	Days       int
	DaysString string

	Hours       int
	HoursString string

	Minutes       int
	MinutesString string

	Seconds       int
	SecondsString string
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
	lambda.Start(Handler)
}

func Handler(ctx context.Context) (output Output, err error) {
	buf := new(bytes.Buffer)

	tmpl, err := template.ParseFS(f, "file.html")
	tmpl.Execute(buf, calcUntil(d))

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

	st := Till{
		Days:    t / 86400,
		Hours:   t % 86400 / 3600,
		Minutes: t % 86400 % 3600 / 60,
		Seconds: t % 86400 % 3600 % 60,
	}

	// check for correct spelling
	if st.Days == 1 {
		st.DaysString = "day"
	} else {
		st.DaysString = "days"
	}

	if st.Hours == 1 {
		st.HoursString = "hour"
	} else {
		st.HoursString = "hours"
	}

	if st.Minutes == 1 {
		st.MinutesString = "minute"
	} else {
		st.MinutesString = "minutes"
	}

	if st.Seconds == 1 {
		st.SecondsString = "second"
	} else {
		st.SecondsString = "seconds"
	}

	return st
}
