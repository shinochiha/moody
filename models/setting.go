package models

import (
	"time"

	"github.com/moody/helpers"
)

type Setting struct {
	Key       string     `json:"key,omitempty" gorm:"primary_key;type:varchar(100)"`
	Value     string     `json:"value,omitempty"`
	CreatedAt time.Time  `json:"created_at,omitempty" example:"2020-03-16T13:55:09.598136+07:00"`
	UpdatedAt time.Time  `json:"updated_at,omitempty" example:"2020-03-16T13:55:09.598136+07:00"`
	DeletedAt *time.Time `json:"-"`
}

type SettingShort struct {
	DBMigrationVersion string `json:"db.migration.version" example:"0036"`
	DBSeedVersion      string `json:"db.seed.version" example:"0011"`
	AddressID          string `json:"origin_address.subdistrict.id" example:"222a761e-7dfc-4586-8f39-0d0fb59bb050"`
}

func (o *Setting) Get(ctx helpers.Context) map[string]interface{} {
	fields := []Setting{}
	helpers.GetDB(ctx).Find(&fields)
	ret := map[string]interface{}{}
	for _, f := range fields {
		ret[f.Key] = f.Value
	}
	return ret
}

func (o *Setting) Update(ctx helpers.Context, data map[string]interface{}) map[string]interface{} {
	for k, v := range data {
		o := Setting{Key: k, Value: v.(string)}
		helpers.GetDB(ctx).Where(Setting{Key: k}).Assign(o).FirstOrCreate(&o)
	}
	return o.Get(ctx)
}

func GetSetting(ctx helpers.Context, key string) string {
	setting := Setting{}
	helpers.GetDB(ctx).Model(&Setting{}).Select("value").Where("`key` = ?", key).First(&setting)
	return setting.Value
}
