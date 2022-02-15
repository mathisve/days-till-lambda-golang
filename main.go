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

var (
	date time.Time
)

//go:embed file.html
var f embed.FS

type Till struct {
	Days    int
	Hours   int
	Minutes int
	Seconds int
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

	date = Date(year, month, day)
}

func main() {
	lambda.Start(Handler)

}

func Handler(ctx context.Context) (output Output, err error) {
	buf := new(bytes.Buffer)

	tmpl, err := template.ParseFS(f, "file.html")
	tmpl.Execute(buf, calcUnitl(date))

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

func Date(year, month, day int) time.Time {
	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
}

func calcUnitl(d time.Time) Till {
	t := int(time.Until(d).Seconds())

	return Till{
		Days:    t / 86400,
		Hours:   t % 86400 / 3600,
		Minutes: t % 86400 % 3600 / 60,
		Seconds: t % 86400 % 3600 % 60,
	}
}
