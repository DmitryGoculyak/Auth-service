package service

import (
	"Auth-service/internal/server"
	"go.uber.org/fx"
)

var Module = fx.Module("service",
	fx.Provide(
		AuthServiceConstructor,
		func(s *AuthService) AuthServiceServer { return s },
	),
	fx.Invoke(
		server.RunServer,
	),
)
