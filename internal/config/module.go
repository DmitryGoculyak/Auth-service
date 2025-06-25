package config

import (
	"Auth-service/internal/db"
	"Auth-service/pkg/client/billing"
	"Auth-service/pkg/client/currency"
	"go.uber.org/fx"
)

var Module = fx.Module("config",
	fx.Provide(
		LoadConfig,
		func(cfg *Config) *billing.BillingClientConfig { return cfg.BillingConfig },
		func(cfg *Config) *currency.CurrencyClientConfig { return cfg.CurrencyConfig },
		func(cfg *Config) *GrpcServiceConfig { return cfg.GrpcConfig },
		func(cfg *Config) *db.DBConfig { return cfg.DBConfig },
	),
)
