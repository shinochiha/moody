package helpers

import (
	"reflect"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
)

const FormatDateOnly = "2006-01-02"

func IsCodeExists(c echo.Context, tableName, fieldName, code string) (res Map, err error) {
	var isExists int64
	tx := GetDBTx(c)
	tx.Table(tableName).Where(fieldName+" = ?", code).Count(&isExists)
	if isExists >= 1 {
		res = Map{"error": Map{
			"code":    400,
			"message": "The " + fieldName + " with " + "'" + code + "'" + " already used",
		}}
		return res, nil
	}
	return res, nil
}

func IsExists(c echo.Context, tableName, fieldName, param string) (res Map, err error) {
	var isExists int64
	tx := GetDBTx(c)
	tx.Table(tableName).Where(fieldName+" = ?", param).Count(&isExists)
	if isExists >= 1 {
		res = Map{
			"error": Map{
				"code":    400,
				"message": "The " + fieldName + " with " + "'" + param + "'" + " already used",
			},
		}
		return res, nil
	}
	return res, nil
}

func NewCode(c echo.Context, tableName, fieldName, baseName string) string {
	replacer := strings.NewReplacer(",", "", ".", "", ";", "")
	l := ""
	baseName = replacer.Replace(baseName)
	words := strings.Fields(baseName)
	for _, word := range words {
		l += word[0:1]
	}
	l += "-"

	//get next code
	var nextCode int64
	tx := GetDBTx(c)
	tx.Table(tableName).Where(fieldName+" LIKE ?", l+"%").Count(&nextCode)
	nextCode += 1
	for i := 0; i < (5 - DigitLen(int(nextCode))); i++ {
		l += "0"
	}
	return strings.ToUpper(l) + strconv.FormatInt(nextCode, 10)
}

func DigitLen(number int) int {
	count := 0
	for number != 0 {
		number /= 10
		count += 1
	}
	return count
}

func InArray(needle interface{}, haystack interface{}) (exists bool, index int) {
	exists = false
	index = -1

	switch reflect.TypeOf(haystack).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(haystack)

		for i := 0; i < s.Len(); i++ {
			if reflect.DeepEqual(needle, s.Index(i).Interface()) {
				index = i
				exists = true
				return
			}
		}
	}

	return
}
