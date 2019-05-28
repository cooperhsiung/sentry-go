package handler

import (
	"github.com/labstack/echo"
	"gitlab.com/sentry-go/influx"
)

func SaveHandler(c echo.Context) error {

	p := new(influx.Point)
	if err := c.Bind(p); err != nil {
		return c.JSON(400, map[string]interface{}{"msg": "format error"})
	}

	if p.Type == "numeric" {
		// save numric
		influx.SaveNum(*p)
	} else if p.Type == "categorical" {
		influx.SaveCate(*p)
	} else {
		return c.JSON(400, map[string]interface{}{"msg": "unsupported type"})

	}

	return c.JSON(200, map[string]interface{}{"msg": "ok"})
}
