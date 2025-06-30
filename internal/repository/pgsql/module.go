package pgsql

import (
	repo "Auth-service/internal/repository"
	"go.uber.org/fx"
)

var Module = fx.Module("pgsql",
	fx.Provide(
		AuthRepoConstructor,
		func(r *AuthRepo) repo.AuthRepository { return r },
	),
)
