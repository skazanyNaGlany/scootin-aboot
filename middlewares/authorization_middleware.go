package middlewares

import (
	"errors"
	"fmt"
	"net/http"
	"scootin-aboot/consts"
	"scootin-aboot/handlers"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// AuthorizationMiddleware represents a middleware for handling authorization in the API.
type AuthorizationMiddleware struct {
	api huma.API
}

func NewAuthorizationMiddleware(api huma.API) *AuthorizationMiddleware {
	return &AuthorizationMiddleware{api: api}
}

// Middleware is a function that performs authorization checks before executing the next handler.
// It checks if the request is for a non-authenticated path, and if so, authorizes the request.
// If the request is for an authenticated path, it parses the user ID from the Authorization header,
// retrieves the user from the UserRepository, and authorizes the request if the user exists.
// If the authorization fails, it returns an error response with a status code of 401 Unauthorized.
// If the authorization succeeds, it calls the next handler in the chain.
func (m *AuthorizationMiddleware) Middleware(ctx huma.Context, next func(huma.Context)) {
	var authorized bool

	if isNonAuthPath(ctx.Method(), ctx.URL().Path) {
		authorized = true
	} else {
		userId, err := uuid.Parse(ctx.Header("Authorization"))

		if err == nil {
			user, err := handlers.UserRepository.FindByID(userId)

			if err != nil {
				if !errors.Is(err, gorm.ErrRecordNotFound) {
					huma.WriteErr(
						m.api,
						ctx,
						http.StatusInternalServerError,
						http.StatusText(http.StatusInternalServerError),
					)
					return
				}
			}

			authorized = user != nil
		}
	}

	if !authorized {
		huma.WriteErr(m.api, ctx, http.StatusUnauthorized,
			"Proper static API key (user ID) is required", fmt.Errorf("invalid API key"),
		)
		return
	}

	next(ctx)
}

func isNonAuthPath(method string, path string) bool {
	return method == "POST" && path == consts.USERS
}
