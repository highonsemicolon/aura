package config

const (
	DBUsers     = "users"
	DBOrders    = "orders"
	DBAnalytics = "analytics"
)

type MongoDB struct {
	URI         string `koanf:"uri" validate:"required,url"`
	UserDB      string `koanf:"user_db" validate:"required"`
	OrderDB     string `koanf:"order_db" validate:"required"`
	AnalyticsDB string `koanf:"analytics_db" validate:"required"`
}

func (m MongoDB) Databases() map[string]string {
	return map[string]string{
		DBUsers:     m.UserDB,
		DBOrders:    m.OrderDB,
		DBAnalytics: m.AnalyticsDB,
	}
}
