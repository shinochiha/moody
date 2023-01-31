package helpers

import (
	"github.com/labstack/echo/v4"
)

type HTTPList struct {
	Count       int `json:"count" example:"1"`
	PageContext struct {
		Page       int `json:"page" example:"1"`
		PerPage    int `json:"per_page" example:"10"`
		TotalPages int `json:"total_pages" example:"1"`
	} `json:"page_context"`
	Results interface{} `json:"results"`
}

type HTTPDeleted struct {
	Code    int    `json:"code" example:"200"`
	Message string `json:"message" example:"Data dengan id = '6e8ef30f-c443-48b7-89e4-964c207245d9' berhasil dihapus."`
}

type HTTPBadRequest struct {
	Error struct {
		Code    int    `json:"code" example:"400"`
		Message string `json:"message" example:"Field xxx wajib diisi."`
	} `json:"error"`
}

type HTTPUnauthorized struct {
	Error struct {
		Code    int    `json:"code" example:"401"`
		Message string `json:"message" example:"Token otentikasi tidak valid."`
	} `json:"error"`
}

type HTTPForbidden struct {
	Error struct {
		Code    int    `json:"code" example:"403"`
		Message string `json:"message" example:"Pengguna tidak memiliki cukup izin untuk mengakses sumber daya."`
	} `json:"error"`
}

type HTTPNotFound struct {
	Error struct {
		Code    int    `json:"code" example:"404"`
		Message string `json:"message" example:"Data dengan id = '6e8ef30f-c443-48b7-89e4-964c207245d9' tidak ditemukan."`
	} `json:"error"`
}

func SuccessResponse(c echo.Context, statusCode int, res interface{}) {
	c.JSON(statusCode, res)
}

func Response(c echo.Context, code int, res map[string]interface{}) error {
	if res["error"] != nil {
		e := res["error"].(map[string]interface{})
		if e["code"] != nil {
			code = e["code"].(int)
			NewCtx(c).SetErrorMessage(res)
			return echo.NewHTTPError(code, e["message"].(string))
		}
	}
	return c.JSON(code, res)
}

func ResponseInternalError(c echo.Context, e error) error {
	return echo.NewHTTPError(500, e.Error())
}

func NotFoundMessage(object, key, id string) map[string]interface{} {
	return map[string]interface{}{
		"error": map[string]interface{}{
			"code":    404,
			"message": object + " with " + key + " = '" + id + "' is not found.",
		},
	}
}

func InternalErrorMessage(message string) map[string]interface{} {
	if message == "" {
		message = "Internal Error"
	}
	return map[string]interface{}{
		"error": map[string]interface{}{
			"code":    500,
			"message": message,
		},
	}
}

func DeletedMessage(object, key, id string) map[string]interface{} {
	return map[string]interface{}{
		"code":    200,
		"message": object + " has been deleted.",
	}
}

func GeneralErrorMessage(code int, message string, detail map[string]interface{}) map[string]interface{} {
	err := map[string]interface{}{}
	err["code"] = code
	err["message"] = message
	if len(detail) > 0 {
		err["detail"] = detail
	}
	return map[string]interface{}{"error": err}
}

func ServiceErrorMessage(code int, service_name, message_key string, detail map[string]interface{}) map[string]interface{} {
	err := map[string]interface{}{}
	err["code"] = code
	err["message"] = "Failed to connect to " + service_name + ", please try again later."
	if detail[message_key] != nil {
		error_message, ok := detail[message_key].(string)
		if ok {
			err["message"] = error_message
		}
	}
	if len(detail) > 0 {
		err[service_name] = detail
	}
	return map[string]interface{}{"error": err}
}

func RequiredValue(fieldName string) (res Map) {
	res = Map{
		"error": Map{
			"error":   400,
			"message": "The field " + fieldName + " is required",
		},
	}
	return res
}
