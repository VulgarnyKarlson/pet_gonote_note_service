package http

import (
	"context"
	"net/http"

	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/adapters/auth"
	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/common/customerrors"
)

const UserCtxKey = "user"

func (w *Server) AuthMiddleware(authClient *auth.Wrapper) func(handler http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := r.Header.Get("Authorization")
			if token == "" || len(token) < 7 {
				http.Error(w, customerrors.ErrUnauthorized.Message, customerrors.ErrUnauthorized.Code)
				return
			}
			token = token[7:]

			resp, err := authClient.ValidateToken(r.Context(), token)
			if !resp.Valid || err != nil {
				http.Error(w, customerrors.ErrUnauthorized.Message, customerrors.ErrUnauthorized.Code)
				return
			}

			ctx := context.WithValue(r.Context(), UserCtxKey, resp.User)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
