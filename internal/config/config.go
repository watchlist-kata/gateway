package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
	"reflect"
)

type Config struct {
	ServerAddress string `env:"SERVER_ADDRESS"`

	KafkaAddress string `env:"KAFKA_ADDRESS"`
	KafkaTopic   string `env:"KAFKA_TOPIC"`

	AuthAddress          string `env:"AUTH_ADDRESS"`
	MediaAddress         string `env:"MEDIA_ADDRESS"`
	SubscriptionsAddress string `env:"SUBSCRIPTIONS_ADDRESS"`
	UserAddress          string `env:"USER_ADDRESS"`
	ReviewsAddress       string `env:"REVIEWS_ADDRESS"`
	WatchListsAddress    string `env:"WATCHLISTS_ADDRESS"`

	ImageBaseURL  string `env:"IMAGE_BASE_URL"`
	ImageCacheDir string `env:"IMAGE_CACHE_DIR"`
}

func Init() (*Config, error) {

	godotenv.Load()

	cfg := &Config{}
	v := reflect.ValueOf(cfg).Elem()
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		tag := t.Field(i).Tag.Get("env")

		if tag == "" {
			continue
		}

		value := os.Getenv(tag)
		if value == "" {
			return nil, fmt.Errorf("%s is not set", tag)
		}

		if field.CanSet() {
			field.SetString(value)
		}
	}

	return cfg, nil
}
