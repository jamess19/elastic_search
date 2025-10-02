package conf

import (
	"github.com/caarlos0/env/v6"
)

// AppConfig presents app conf
type AppConfig struct {
	AppEnv    string `env:"APP_ENV" envDefault:"dev"`
	Port      string `env:"PORT" envDefault:"8000"`
	LogFormat string `env:"LOG_FORMAT" envDefault:"text"`
	DBHost    string `env:"DB_HOST" envDefault:"dbmasternode.stg.int.finan.cc"`
	DBPort    string `env:"DB_PORT" envDefault:"5432"`
	DBUser    string `env:"DB_USER" envDefault:"finan_dev_user"`
	DBPass    string `env:"DB_PASS" envDefault:"Oo5Tah0re5eexoif"`
	DBName    string `env:"DB_NAME" envDefault:"finan_dev_ms_user_management"`
	EnableDB  string `env:"ENABLE_DB" envDefault:"true"`

	MSConsumer string `env:"MS_CONSUMER" envDefault:"http://ms-consumer"`
	FinanBusiness string `env:"FINAN_BUSINESS" envDefault:"http://finan-business"`
}

var config AppConfig

func SetEnv() {
	_ = env.Parse(&config)
}

func LoadEnv() AppConfig {
	return config
}

