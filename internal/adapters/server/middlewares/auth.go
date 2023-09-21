package middlewares

import (
	"context"
	"net/http"

	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/adapters/auth"
	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/adapters/server"
	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/common/customerrors"
)

const UserCtxKey = "user"
const middlewareAuthID = "auth"

type AuthMiddleware struct {
	auth auth.Client
}

func RegisterAuthMiddleware(client auth.Client, s server.Server) {
	m := &AuthMiddleware{auth: client}
	s.Use(m.Middleware, middlewareAuthID)
}

func AuthID() string {
	return middlewareAuthID
}

func (w *AuthMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, reader *http.Request) {
		token := reader.Header.Get("Authorization")
		if token == "" || len(token) < 7 || token[:7] != "Bearer " {
			http.Error(writer, customerrors.ErrUnauthorized.Message, customerrors.ErrUnauthorized.Code)
			return
		}
		token = token[7:]
		resp, err := w.auth.ValidateToken(reader.Context(), token)
		if err != nil || resp == nil || !resp.Valid {
			http.Error(writer, customerrors.ErrUnauthorized.Message, customerrors.ErrUnauthorized.Code)
			return
		}

		ctx := context.WithValue(reader.Context(), UserCtxKey, resp.User)
		next.ServeHTTP(writer, reader.WithContext(ctx))
	})
}
