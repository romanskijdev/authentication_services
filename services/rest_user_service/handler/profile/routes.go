package userhandler

import (
	errm "authentication_service/core/errmodule"
	"authentication_service/rest_user_service/handler"
	typesm "authentication_service/rest_user_service/types"
	"github.com/go-chi/chi/v5"
	"net/http"
)

const (
	profileURI = "/profile"
)

type UsersReg struct {
	ipc *typesm.InternalProviderControl
}

// RegisterUsersRoutes регистрирует маршруты для группы user
func RegisterUsersRoutes(
	r chi.Router,
	ipc *typesm.InternalProviderControl,
) *errm.Error {

	s := &UsersReg{
		ipc: ipc,
	}

	r.Route("/api/users", func(r chi.Router) {
		r.Use(handler.JWTVerifier(ipc.Config))

		handler.RegisterRoute(r, http.MethodGet, profileURI, s.GetProfileHandler)
	})

	return nil
}
