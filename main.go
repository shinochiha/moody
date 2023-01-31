package main

import (
	"sync"

	"github.com/moody/config"
	"github.com/moody/db"
	"github.com/moody/routes"

	"github.com/jinzhu/gorm"
	echoSwagger "github.com/swaggo/echo-swagger"
)

var (
	syncOnce sync.Once
	dbh      *gorm.DB
)

func main() {
	syncOnce.Do(DBConnection)
	defer dbh.Close()

	r := routes.Routes(dbh)
	r.GET("/swagger/*", echoSwagger.WrapHandler)

	err := r.Start(":" + config.Get("APP_PORT").String())
	if err != nil {
		r.Logger.Fatal(err)
	}
}

func DBConnection() {
	dbh = db.Connect()
	// db.ConnectRedis()
	db.Migrate()
	db.Seed()
}

// @securityDefinitions.apikey ApiKeyAuth
// @in cookies
// @name Token

// @description
// @description ## Token
// @description
// @description ```
// @description Token in cookies
// @description ```
// @description
// @description ## Query params
// @description
// @description By default, we support a common way for selecting fields, filtering, searching, sorting, and pagination in URL query params on `GET` method:
// @description
// @description ### Field
// @description
// @description Get selected fields in GET result, example:
// @description ```
// @description GET /api/resources?fields=field_a,field_b,field_c
// @description ```
// @description equivalent to sql:
// @description ```sql
// @description SELECT field_a, field_b, field_c FROM resources
// @description ```
// @description
// @description ### Filter
// @description
// @description Adds fields request condition (multiple conditions) to the request, example:
// @description ```
// @description GET /api/resources?field_a=value_a&field_b.$gte=value_b&field_c.$like=value_c&field_d.$ilike=value_d%
// @description ```
// @description equivalent to sql:
// @description ```sql
// @description SELECT * FROM resources WHERE (field_a = 'value_a') AND (field_b >= value_b) AND (field_c LIKE '%value_c%') AND (LOWER(field_d) LIKE LOWER('value_d%'))
// @description ```
// @description
// @description #### Available filter conditions
// @description
// @description * `$eq`: equal (`=`)
// @description * `$ne`: not equal (`!=`)
// @description * `$gt`: greater than (`>`)
// @description * `$gte`: greater than or equal (`>=`)
// @description * `$lt`: lower than (`<`)
// @description * `$lte`: lower than or equal (`<=`)
// @description * `$like`: contains (`LIKE '%value%'`)
// @description * `$ilike`: contains case insensitive (`LOWER(field) LIKE LOWER('%value%')`)
// @description * `$nlike`: not contains (`NOT LIKE '%value%'`)
// @description * `$nilike`: not contains case insensitive (`LOWER(field) NOT LIKE LOWER('%value%')`)
// @description * `$in`: in range, accepts multiple values (`IN ('value_a', 'value_b')`)
// @description * `$nin`: not in range, accepts multiple values (`NOT IN ('value_a', 'value_b')`)
// @description * `$regexp`: regex (`REGEXP '%value%'`)
// @description * `$nregexp`: not regex (`NOT REGEXP '%value%'`)
// @description
// @description ### Or
// @description
// @description Adds `OR` conditions to the request, example:
// @description ```
// @description GET /api/resources?or=field_a:val_a|field_b.$gte:val_b;field_c.$lte:val_c|field_d.$like:val_d
// @description ```
// @description equivalent to sql:
// @description ```sql
// @description SELECT * FROM resources WHERE (field_a=val_a OR field_b <= val_b) AND (field_c <= val_c OR field_d LIKE '%val_d%')
// @description ```
// @description
// @description ### Search
// @description
// @description Adds a search conditions to the request, example:
// @description ```
// @description GET /api/resources?search=field_a,field_b:term_1;field_c,field_d:term_2
// @description ```
// @description equivalent to sql:
// @description ```sql
// @description SELECT * FROM resources WHERE (LOWER(field_a) LIKE LOWER('%term_1%') OR LOWER(field_b) LIKE LOWER('%term_1%')) AND (LOWER(field_c) LIKE LOWER('%term_2%') OR LOWER(field_d) LIKE LOWER('%term_2%'))
// @description ```
// @description
// @description ### Sort
// @description
// @description Adds sort by field (by multiple fields) and order to query result, example:
// @description ```
// @description GET /api/resources?sorts=field_a,-field_b,field_c:i,-field_d:i
// @description ```
// @description equivalent to sql:
// @description ```sql
// @description SELECT * FROM resources ORDER BY field_a ASC, field_b DESC, LOWER(field_c) ASC, LOWER(field_d) DESC
// @description ```
// @description
// @description ### Page
// @description
// @description Specify the page of results to return, example:
// @description ```
// @description GET /api/resources?page=3&per_page=10
// @description ```
// @description equivalent to sql:
// @description ```sql
// @description SELECT * FROM resources LIMIT 10 OFFSET 20
// @description ```
// @description
// @description ### Per Page
// @description
// @description Specify the number of records to return in one request, example:
// @description ```
// @description GET /api/resources?per_page=10
// @description ```
// @description equivalent to sql:
// @description ```sql
// @description SELECT * FROM resources LIMIT 10
// @description ```
