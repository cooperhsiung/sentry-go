package influx

import (
	"encoding/json"
	"fmt"
	"github.com/influxdata/influxdb/client/v2"
	"log"
	"time"
)

var db *client.Client

func InitInflux() *client.Client {

	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr: "http://localhost:8086",
	})

	if err != nil {
		log.Fatal(err)
	}

	db = &c
	return &c
}

func SaveNum(point Point) {

	bp, _ := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  "sentry_temp",
		Precision: "s",
	})

	measurement := point.Project + "." + point.Field
	var tags map[string]string
	fields := map[string]interface{}{
		"value": point.Value,
	}
	pt, err := client.NewPoint(measurement, tags, fields, time.Now())

	if err != nil {
		log.Fatal(err)
	}

	bp.AddPoint(pt)

	err = (*db).Write(bp)

	if err != nil {
		log.Fatal(err)
	}

}

func SaveCate(point Point) {

	bp, _ := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  "sentry_temp",
		Precision: "s",
	})

	measurement := point.Project + "." + point.Field
	tags := map[string]string{
		"value": point.Value.(string),
	}
	fields := map[string]interface{}{
		"_v": 1,
	}

	pt, err := client.NewPoint(measurement, tags, fields, time.Now())

	if err != nil {
		log.Fatal(err)
	}

	bp.AddPoint(pt)

	err = (*db).Write(bp)

	if err != nil {
		log.Fatal(err)
	}
}

func SaveCalc(point Point) {

	bp, _ := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  "sentry",
		Precision: "s",
	})

	var measurement string
	if point.Tag != "" {
		measurement = fmt.Sprintf("%s.%s.(%s).%s", point.Project, point.Field, point.Tag, point.Method)
	} else {
		measurement = fmt.Sprintf("%s.%s.%s", point.Project, point.Field, point.Method)
	}

	var tags map[string]string

	fields := map[string]interface{}{
		"value": point.Value,
	}

	pt, err := client.NewPoint(measurement, tags, fields, time.Now())

	if err != nil {
		log.Fatal(err)
	}

	bp.AddPoint(pt)

	err = (*db).Write(bp)

	if err != nil {
		log.Fatal(err)
	}

}

func QueryMean(qs Qs) float64 {
	q := client.Query{
		Command: fmt.Sprintf(`
    select mean(value) from "%s.%s"
    where time >= now() - %s - %s and time < now() - %s
    `, qs.Project, qs.Field, qs.Period, qs.Offset, qs.Offset),

		Database:        "sentry_temp",
		RetentionPolicy: "one_day",
	}

	resp, err := (*db).Query(q)
	if err != nil {
		log.Fatal(err)
	}

	result := resp.Results[0].Series[0].Values[0][1]
	value, err := result.(json.Number).Float64()
	if err != nil {
		log.Fatal(err)
	}
	return value

}

func QuerySum(qs Qs) float64 {
	q := client.Query{
		Command: fmt.Sprintf(`
    select sum(value) from "%s.%s"
    where time >= now() - %s - %s and time < now() - %s
    `, qs.Project, qs.Field, qs.Period, qs.Offset, qs.Offset),
		Database:        "sentry_temp",
		RetentionPolicy: "one_day",
	}

	resp, err := (*db).Query(q)
	if err != nil {
		log.Fatal(err)
	}

	result := resp.Results[0].Series[0].Values[0][1]
	value, err := result.(json.Number).Float64()
	if err != nil {
		log.Fatal(err)
	}
	return value
}

func QueryCount(qs Qs) float64 {
	q := client.Query{
		Command: fmt.Sprintf(`
    select count(_v) from "%s.%s"
    where value = '%s'
    and time >= now() - %s - %s and time < now() - %s
    `, qs.Project, qs.Field, qs.Value, qs.Period, qs.Offset, qs.Offset),
		Database:        "sentry_temp",
		RetentionPolicy: "one_day",
	}

	resp, err := (*db).Query(q)
	if err != nil {
		log.Fatal(err)
	}

	result := resp.Results[0].Series[0].Values[0][1]
	value, err := result.(json.Number).Float64()
	if err != nil {
		log.Fatal(err)
	}
	return value
}

func QueryCountAll(qs Qs) float64 {
	q := client.Query{
		Command: fmt.Sprintf(`
    select count(_v) from "%s.%s"
    where time >= now() - %s - %s and time < now() - %s
    `, qs.Project, qs.Field, qs.Period, qs.Offset, qs.Offset),
		Database:        "sentry_temp",
		RetentionPolicy: "one_day",
	}

	resp, err := (*db).Query(q)
	if err != nil {
		log.Fatal(err)
	}

	result := resp.Results[0].Series[0].Values[0][1]
	value, err := result.(json.Number).Float64()
	if err != nil {
		log.Fatal(err)
	}
	return value
}

func QueryPercent(qs Qs) float64 {
	partCount := QueryCount(qs)
	totalCount := QueryCountAll(qs)
	return partCount / totalCount
}
