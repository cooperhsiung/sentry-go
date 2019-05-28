package alert

import (
	"errors"
	"fmt"
	"github.com/BurntSushi/toml"
	"gitlab.com/sentry-go/influx"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func Start() {
	var conf Conf
	var env = "dev"
	if os.Getenv("GO_ENV") == "prod" {
		env = "prod"
	}
	_, err := toml.DecodeFile("../config/"+env+".toml", &conf)
	if err != nil {
		fmt.Println(err)
	}

	for _, task := range conf.Tasks {
		Process(task, time.Now())
	}
}

func Process(task Task, fireTime time.Time) {

	var current float64

	for i, rule := range task.Rules {

		expr := ParseExpr(rule.Expr)

		if i == 0 { // 待优化，免得存两次
			current = AdaptQuery(task, expr)
			calcPoint := influx.Point{
				Project: task.Project,
				Field:   task.Field,
				Method:  strings.ToUpper(expr.Method),
				Tag:     expr.Value,
				Value:   current,
			}
			influx.SaveCalc(calcPoint)
		}

		if expr.Operator == ">" {
			if current-expr.Threshold > 0 {
				fmt.Println("triggered", task, rule)
				t, m := BuildMessage(task, rule, current)
				SendNotice(task, rule, t, m)
			}
		}

		if expr.Operator == "<" {
			if current-expr.Threshold < 0 {
				fmt.Println("triggered", task, rule)
				t, m := BuildMessage(task, rule, current)
				SendNotice(task, rule, t, m)
			}
		}
	}

}

func AdaptQuery(task Task, expr Expr) float64 {
	var qs = influx.Qs{
		Project: task.Project,
		Field:   task.Field,
		Period:  task.Period,
		Value:   expr.Value,
		Offset:  task.Offset,
	}
	op := strings.ToUpper(expr.Method)
	if op == "PERCENT" {
		return influx.QueryPercent(qs)
	}
	if op == "MEAN" {
		return influx.QueryMean(qs)
	}
	if op == "SUM" {
		return influx.QuerySum(qs)
	}
	if op == "COUNT" {
		return influx.QueryCount(qs)
	}
	// more op
	return 0
}

func BuildMessage(task Task, rule Rule, current float64) (title string, message string) {
	title = task.Title
	message = task.Message

	r, _ := regexp.Compile("{\\s*project\\s*}")
	title = r.ReplaceAllString(title, task.Project)
	message = r.ReplaceAllString(message, task.Project)

	r, _ = regexp.Compile("{\\s*field\\s*}")
	title = r.ReplaceAllString(title, task.Field)
	message = r.ReplaceAllString(message, task.Field)

	r, _ = regexp.Compile("{\\s*level\\s*}")
	title = r.ReplaceAllString(title, rule.Level)
	message = r.ReplaceAllString(message, rule.Level)

	r, _ = regexp.Compile("{\\s*value\\s*}")
	title = r.ReplaceAllString(title, fmt.Sprintf("%f", current))
	message = r.ReplaceAllString(message, fmt.Sprintf("%f", current))

	return

}

func ParseExpr(exprStr string) Expr {

	re, err := regexp.Compile(`(\w+)\((.*?)\) (.*?) (.*)`)
	if err != nil {
		log.Fatal(err)
	}

	m := re.FindStringSubmatch(exprStr)
	if len(m) != 5 {
		panic(errors.New("incorrect condition"))
	}

	v, _ := strconv.ParseFloat(m[4], 64)

	if len(m[3]) > 2 {
		log.Fatal("illegal operator")
	}

	return Expr{
		Method:    m[1],
		Value:     m[2],
		Operator:  m[3],
		Threshold: v,
	}

}

func SendNotice(task Task, rule Rule, title string, message string) {
	if len(task.Mail) > 0 {
		fmt.Println(task.Mail)
	}

	if len(task.Wechat) > 0 {
		fmt.Println(task.Wechat)
	}
}
