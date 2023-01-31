package helpers

import (
	"fmt"
	"reflect"
	"time"
)

type DateType time.Time

func (t *DateType) UnmarshalJSON(b []byte) error {
	src := string(b)
	fmt.Println(src)
	dt, err := time.Parse("2006-01-02", src)
	if err != nil {
		fmt.Println("masuk sini")
		return err
	}
	*t = DateType(dt)
	return nil
}

func GetStructName(s interface{}) string {
	if t := reflect.TypeOf(s); t.Kind() == reflect.Ptr {
		return t.Elem().Name()
	} else {
		return t.Name()
	}
}
