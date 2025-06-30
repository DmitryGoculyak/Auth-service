package container

import (
	"Auth-service/internal/config"
	"Auth-service/internal/db"
	repo "Auth-service/internal/repository/pgsql"
	"Auth-service/internal/service"
	"Auth-service/internal/transport/rpc/handlers"
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
		repo.Module,
		handlers.Module,
	)
}
