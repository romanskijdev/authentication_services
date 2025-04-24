package authhandler

import (
	errm "authentication_service/core/errmodule"
	"authentication_service/rest_user_service/handler"
	typesm "authentication_service/rest_user_service/types"
	"github.com/go-chi/chi/v5"
	"net/http"
)

const (
	authURI    = "/issue"
	refreshURI = "/refresh"
)

type AuthReg struct {
	ipc *typesm.InternalProviderControl
}

// RegisterAuthRoutes регистрирует маршруты для группы auth
func RegisterAuthRoutes(
	r chi.Router,
	ipc *typesm.InternalProviderControl,
) *errm.Error {

	s := &AuthReg{
		ipc: ipc,
	}

	r.Route("/api/auth", func(r chi.Router) {
		handler.RegisterRoute(r, http.MethodPost, authURI, s.IssueTokensHandler)
		handler.RegisterRoute(r, http.MethodPost, refreshURI, s.RefreshTokensHandler)
	})

	return nil
}
