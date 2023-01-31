package db

import (
	"context"
	"fmt"
	"log"

	"github.com/go-redis/redis/v8"
	"github.com/jinzhu/gorm"

	_ "github.com/lib/pq"

	"github.com/moody/config"
	"github.com/moody/helpers"
)

var db *gorm.DB

func Connect() *gorm.DB {
	host := config.Get("DB_HOST").String()
	driver := config.Get("DB_DRIVER").String()
	port := config.Get("DB_PORT").String()
	name := config.Get("DB_NAME").String()
	user := config.Get("DB_USER").String()
	password := config.Get("DB_PASSWORD").String()

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable TimeZone=Asia/Shanghai", host, port, user, password, name)
	dbConnPool, err := gorm.Open(driver, psqlInfo)
	// untuk connect ke mysql
	// dbConnPool, err := gorm.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", user, password, host, port, name))
	if err != nil {
		fmt.Println(err.Error())
		fmt.Println("-------------------------------------------------------")
		fmt.Println("Connect Ke Database Gagal , Silahkan cek settingan nya pastikan pake posgresql")
		log.Fatal(err)
	}

	if err := dbConnPool.DB().Ping(); err != nil {
		log.Fatal(err)
	}

	if config.Get("DB_IS_DEBUG").Bool() {
		dbConnPool = dbConnPool.Debug()
	}

	maxOpenConns := config.Get("DB_MAX_OPEN_CONNS").Int()
	maxIdleConns := config.Get("DB_MAX_IDLE_CONNS").Int()
	connMaxLifetime := config.Get("DB_CONN_MAX_LIFETIME").Duration()

	dbConnPool.DB().SetMaxIdleConns(maxIdleConns)
	dbConnPool.DB().SetMaxOpenConns(maxOpenConns)
	dbConnPool.DB().SetConnMaxLifetime(connMaxLifetime)

	db = dbConnPool
	return db
}

func ConnectRedis() {
	helpers.Rdb = redis.NewClient(&redis.Options{
		Addr:     config.Get("REDIS_HOST").String() + ":" + config.Get("REDIS_PORT").String(),
		Password: config.Get("REDIS_PASSWORD").String(),
		DB:       config.Get("REDIS_DB").Int(),
	})

	_, err := helpers.Rdb.Ping(context.Background()).Result()
	if err != nil {
		log.Fatal(err)
	}
}
