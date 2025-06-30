package handlers

import (
	proto "Auth-service/pkg/proto/auth"
	"go.uber.org/fx"
)

var Module = fx.Module("handlers",
	fx.Provide(
		AuthHandlerConstructor,
		func(h *AuthHandler) proto.AuthServiceServer { return h },
	),
)
