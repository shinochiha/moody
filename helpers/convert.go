package helpers

import (
	"encoding/json"
	"math/big"
	"regexp"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/leekchan/accounting"
	"github.com/tidwall/sjson"
)

type Iconvert struct {
	Val interface{}
}

var link = regexp.MustCompile("(^[A-Za-z])|_([A-Za-z])")
var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

func ToCamelCase(str string) string {
	return link.ReplaceAllStringFunc(str, func(s string) string {
		return strings.ToUpper(strings.Replace(s, "_", "", -1))
	})
}

func ToSnakeCase(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

func Float64ToCurrency(n float64) string {
	ac := accounting.Accounting{Symbol: "Rp ", Thousand: ".", Decimal: ",", Precision: 0}
	return ac.FormatMoneyBigFloat(big.NewFloat(n))
}

func FlatResponse(data []map[string]interface{}, schema map[string]interface{}) []map[string]interface{} {
	res := []map[string]interface{}{}

	for _, d := range data {
		temp := map[string]interface{}{}
		for k, v := range d {
			for f, c := range schema["fields"].(map[string]map[string]string) {
				if k == c["as"] {
					if v != nil {
						if c["type"] == "float64" {
							val, _ := strconv.ParseFloat(string(v.([]byte)), 64)
							temp[f] = val
						} else if c["type"] == "int" || c["type"] == "int32" || c["type"] == "int64" {
							val, _ := v.(int64)
							temp[f] = val
						} else if c["type"] == "bool" {
							val, _ := v.(bool)
							temp[f] = val
						} else if c["type"] == "pq.StringArray" {
							val := string(v.([]byte))
							val = strings.ReplaceAll(val, "{", "")
							val = strings.ReplaceAll(val, "}", "")
							temp[f] = val
						} else {
							temp[f] = v
						}
					} else {
						temp[f] = v
					}
				}
			}
		}
		res = append(res, temp)
	}

	return res
}

func DotToInterface(data []map[string]interface{}, schema map[string]interface{}) []map[string]interface{} {
	temp := `[]`
	i := -1
	for _, d := range data {
		i++
		for k, v := range d {
			for f, c := range schema["fields"].(map[string]map[string]string) {
				if k == c["as"] && v != nil {
					if c["type"] == "int" || c["type"] == "float64" {
						val, _ := strconv.ParseFloat(string(v.([]byte)), 64)
						temp = DotNotationSet(temp, strconv.Itoa(i)+"."+f, val)
					} else if c["type"] == "int64" {
						temp = DotNotationSet(temp, strconv.Itoa(i)+"."+f, v.(int64))
					} else if c["type"] == "bool" {
						val, _ := strconv.ParseBool(string(v.([]byte)))
						temp = DotNotationSet(temp, strconv.Itoa(i)+"."+f, val)
					} else if c["type"] == "pq.StringArray" {
						val := string(v.([]byte))
						val = strings.ReplaceAll(val, "{", "")
						val = strings.ReplaceAll(val, "}", "")
						temp = DotNotationSet(temp, strconv.Itoa(i)+"."+f, strings.Split(val, ","))
					} else if v != "" {
						temp = DotNotationSet(temp, strconv.Itoa(i)+"."+f, v)
					}
				}
			}
		}
	}

	var res []map[string]interface{}
	json.Unmarshal([]byte(temp), &res)
	return res
}

func DotNotationSet(json, path string, value interface{}) string {
	res, _ := sjson.Set(json, path, value)
	return res
}

func BoolAddr(b bool) *bool {
	boolVar := b
	return &boolVar
}

func StringAddr(b string) *string {
	stringVar := b
	return &stringVar
}

func FloatAddr(b float64) *float64 {
	floatVar := b
	return &floatVar
}

func NewUUID() string {
	return uuid.New().String()
}

func NewToken() string {
	return strings.ReplaceAll(NewUUID(), "-", "")
}

func Convert(val interface{}) Iconvert {
	return Iconvert{Val: val}
}

func (v Iconvert) String() string {
	switch v.Val.(type) {
	case string:
		return v.Val.(string)
	case bool:
		return strconv.FormatBool(v.Val.(bool))
	case float32:
		return strconv.FormatFloat(float64(v.Val.(float32)), 'f', -1, 32)
	case float64:
		return strconv.FormatFloat(v.Val.(float64), 'f', -1, 64)
	case uint:
		return strconv.FormatUint(uint64(v.Val.(uint)), 10)
	case uint8:
		return strconv.FormatUint(uint64(v.Val.(uint8)), 10)
	case uint16:
		return strconv.FormatUint(uint64(v.Val.(uint16)), 10)
	case uint32:
		return strconv.FormatUint(uint64(v.Val.(uint32)), 10)
	case uint64:
		return strconv.FormatUint(v.Val.(uint64), 10)
	case uintptr:
		return strconv.FormatUint(uint64(v.Val.(uintptr)), 10)
	case int:
		return strconv.FormatInt(int64(v.Val.(int)), 10)
	case int8:
		return strconv.FormatInt(int64(v.Val.(int8)), 10)
	case int16:
		return strconv.FormatInt(int64(v.Val.(int16)), 10)
	case int32:
		return strconv.FormatInt(int64(v.Val.(int32)), 10)
	case int64:
		return strconv.FormatInt(v.Val.(int64), 10)
	default:
		return ""
	}
}

func (v Iconvert) Int() int {
	val, err := strconv.Atoi(v.String())
	if err != nil {
		return 0
	}
	return val
}

func InterfaceToMap(data interface{}) map[string]interface{} {
	res, ok := data.(map[string]interface{})
	if !ok {
		data_str, ok := data.(string)
		if ok {
			json.Unmarshal([]byte(data_str), &res)
		}
	}
	return res
}

func InterfaceToSlice(data interface{}) []interface{} {
	res, ok := data.([]interface{})
	if !ok {
		data_str, ok := data.(string)
		if ok {
			json.Unmarshal([]byte(data_str), &res)
		}
	}
	return res
}
