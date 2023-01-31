package helpers

import (
	"database/sql"

	"github.com/go-redis/redis/v8"
	"github.com/jinzhu/gorm"
)

var Rdb *redis.Client

func GetDB(ctx Context) *gorm.DB {
	return ctx.Get("db").(*gorm.DB)
}

func GetDBTx(ctx Context) *gorm.DB {
	return ctx.Get("tx").(*gorm.DB)
}

func GetResults(rows *sql.Rows) []map[string]interface{} {
	columns, err := rows.Columns()
	if err != nil {
		panic(err)
	}
	length := len(columns)
	result := make([]map[string]interface{}, 0)
	for rows.Next() {
		current := makeResultReceiver(length)
		if err := rows.Scan(current...); err != nil {
			panic(err)
		}
		value := make(map[string]interface{})
		for i := 0; i < length; i++ {
			value[columns[i]] = *(current[i]).(*interface{})
		}
		result = append(result, value)
	}
	return result
}

func makeResultReceiver(length int) []interface{} {
	result := make([]interface{}, 0, length)
	for i := 0; i < length; i++ {
		var current interface{}
		current = struct{}{}
		result = append(result, &current)
	}
	return result
}
