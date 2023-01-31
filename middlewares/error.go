package middlewares

import (
	"encoding/json"
	"fmt"
	"runtime"
	"sort"

	"github.com/labstack/echo/v4"

	"github.com/moody/config"
	"github.com/moody/helpers"
)

func ErrorHandler(err error, c echo.Context) {
	he, ok := err.(*echo.HTTPError)
	if ok {
		if he.Internal != nil {
			if herr, ok := he.Internal.(*echo.HTTPError); ok {
				he = herr
			}
		}
	} else {
		he = &echo.HTTPError{
			Code:    500,
			Message: "Failed to connect to the server, please try again later.",
		}
	}

	code := he.Code
	message := helpers.NewCtx(c).ErrorMessage()
	if len(message) == 0 {
		message = map[string]interface{}{
			"error": map[string]interface{}{
				"code":    code,
				"message": he.Message,
			},
		}
	}
	env := config.Get("APP_ENV").String()
	if code >= 500 && (env == "production" || env == "development") {
		temp := map[int]string{}
		trace := []string{}
		for i := 0; i <= 15; i++ {
			fun, file, no, _ := runtime.Caller(i)
			if file != "" {
				temp[i] = fmt.Sprintf("%s:%d on %s", file, no, runtime.FuncForPC(fun).Name())
			}
		}
		index := make([]int, 0)
		for i := range temp {
			index = append(index, i)
		}
		sort.Ints(index)
		for _, i := range index {
			trace = append(trace, fmt.Sprintf("#%d ", i)+temp[i])
		}

		log := map[string]interface{}{}
		log["env"] = env
		log["error"] = err.Error()
		log["request"] = c.Request().Method + " " + c.Path()
		log["body"] = BindBodyRequest(c)
		log["response"] = message
		log["trace"] = trace

		dataJson, _ := json.MarshalIndent(log, "", "  ")
		if env == "production" || env == "development" {
			b := map[string]string{"text": "```" + string(dataJson) + "```"}
			helpers.CallAPI("POST", config.Get("SLACK_URL").String(), b, map[string]string{})
		} else {
			fmt.Println(string(dataJson))
		}
	}
	if !c.Response().Committed {
		if c.Request().Method == "HEAD" {
			c.NoContent(he.Code)
		} else {
			c.JSON(code, message)
		}
	}
}

func BindBodyRequest(c echo.Context) echo.Map {
	body := echo.Map{}
	c.Bind(&body)
	return body
}
