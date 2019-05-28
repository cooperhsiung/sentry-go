package main

import (
	"errors"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/robfig/cron"
	"gitlab.com/sentry-go/alert"
	"log"
	"os"
	"regexp"
	"time"
)

func main() {
	done := make(chan bool)

	log.Println("schedule start..")

	var conf alert.Conf
	var env = "dev"
	if os.Getenv("GO_ENV") == "prod" {
		env = "prod"
	}
	_, err := toml.DecodeFile("../config/"+env+".toml", &conf)
	if err != nil {
		fmt.Println(err)
	}

	c := cron.New()

	for _, task := range conf.Tasks {
		s := ParseCron(task.Every)

		fmt.Println("cron string:", s)

		err = c.AddFunc(s, func() {
			fmt.Println(time.Now().Format(time.RFC3339))
			alert.Process(task, time.Now())
		})
		if err != nil {
			log.Fatal(err)
		}
	}

	c.Start()

	<-done
}

func ParseCron(every string) string {
	re, _ := regexp.Compile(`(\d+)([mhs])`)

	m := re.FindStringSubmatch(every)

	if m[2] == "m" {
		return fmt.Sprintf("0 */%s * * * *", m[1])
	}
	if m[2] == "h" {
		return fmt.Sprintf("0 0 */%s * * *", m[1])
	}
	if m[2] == "s" {
		return fmt.Sprintf("*/%s * * * * *", m[1])
	} else {
		panic(errors.New("incorrect every"))
	}

	return ""
}
