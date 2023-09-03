package http

import (
	"context"
	"net/http"

	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/common/customerrors"
)

func (s *Server) AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token == "" || len(token) < 7 || token[:7] != "Bearer " {
			http.Error(w, customerrors.ErrUnauthorized.Message, customerrors.ErrUnauthorized.Code)
			return
		}
		token = token[7:]
		resp, err := s.auth.ValidateToken(r.Context(), token)
		if !resp.Valid || err != nil {
			http.Error(w, customerrors.ErrUnauthorized.Message, customerrors.ErrUnauthorized.Code)
			return
		}

		ctx := context.WithValue(r.Context(), UserCtxKey, resp.User)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}
