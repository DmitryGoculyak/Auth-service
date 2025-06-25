package container

import (
	"Auth-service/internal/config"
	"Auth-service/internal/db"
	"Auth-service/internal/service"
	"Auth-service/pkg/client/billing"
	"Auth-service/pkg/client/currency"
	"Auth-service/validation"
	"go.uber.org/fx"
)

func Build() *fx.App {
	return fx.New(
		config.Module,
		db.Module,
		validation.Module,
		service.Module,
		billing.Module,
		currency.Module,
	)
}
