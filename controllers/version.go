package controllers

import (
	"github.com/labstack/echo/v4"
	"github.com/moody/config"
	"github.com/moody/helpers"
)

func Version(c echo.Context) error {
	res := helpers.Map{
		"version": config.Get("APP_VERSION").String(),
	}

	return helpers.Response(c, 201, res)
}
