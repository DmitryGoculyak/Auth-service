package config

import (
	"Auth-service/internal/db"
	"Auth-service/pkg/client/billing"
	"Auth-service/pkg/client/currency"
	"Auth-service/pkg/jwt"
	"fmt"
	"github.com/spf13/viper"
	"sync"
)

type GrpcServiceConfig struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

var (
	err    error
	config *Config
	s      sync.Once
)

type Config struct {
	BillingConfig  *billing.BillingClientConfig
	CurrencyConfig *currency.CurrencyClientConfig
	GrpcConfig     *GrpcServiceConfig
	DBConfig       *db.DBConfig
	JwtConfig      *jwt.JWTConfig
}

func LoadConfig() (*Config, error) {

	s.Do(func() {
		config = &Config{}

		viper.AddConfigPath(".")
		viper.SetConfigName("config")

		if err = viper.ReadInConfig(); err != nil {
			return
		}

		BillingConfig := viper.Sub("billing")
		CurrencyConfig := viper.Sub("currency")
		GrpcConfig := viper.Sub("service")
		DBConfig := viper.Sub("database")
		JwtConfig := viper.Sub("jwt")

		if err = parseSubConfig(BillingConfig, &config.BillingConfig); err != nil {
			return
		}
		if err = parseSubConfig(CurrencyConfig, &config.CurrencyConfig); err != nil {
			return
		}
		if err = parseSubConfig(GrpcConfig, &config.GrpcConfig); err != nil {
			return
		}
		if err = parseSubConfig(DBConfig, &config.DBConfig); err != nil {
			return
		}
		if err = parseSubConfig(JwtConfig, &config.JwtConfig); err != nil {
			return
		}
	})
	return config, err
}

func parseSubConfig[T any](subConfig *viper.Viper, parseTo *T) error {
	if subConfig == nil {
		return fmt.Errorf("can not read %T config: subconfig is nil", parseTo)
	}

	if err = subConfig.Unmarshal(parseTo); err != nil {
		return err
	}
	return nil
}
