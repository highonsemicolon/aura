package config

import "reflect"

const (
	DBUsers     = "users"
	DBOrders    = "orders"
	DBAnalytics = "analytics"
)

type MongoDB struct {
	URI         string `koanf:"uri" validate:"required,url"`
	UserDB      string `koanf:"user_db" validate:"required" dbalias:"users"`
	OrderDB     string `koanf:"order_db" validate:"required" dbalias:"orders"`
	AnalyticsDB string `koanf:"analytics_db" validate:"required" dbalias:"analytics"`
}

func (m MongoDB) Databases() map[string]string {
	dbs := make(map[string]string)
	t := reflect.TypeOf(m)
	v := reflect.ValueOf(m)

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		dbalias := field.Tag.Get("dbalias")
		if dbalias == "" {
			continue
		}

		val := v.Field(i).String()
		if val != "" {
			dbs[dbalias] = val
		}
	}
	return dbs
}
