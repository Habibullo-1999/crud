package middleware

import (
	"context"
	"errors"
	"net/http"

)

const (
	MANEGER  = "MANAGER"
	ADMIN	= "ADMIN"
)

var ErrNoAuthentication = errors.New("No authentication")

var authenticationContextKey = &contextKey{"authentication context"}

type contextKey struct {
	name string
}

func (c *contextKey) String() string {
	return c.name
}

type HasAnyRoleFunc func (ctx context.Context, roles ...string)  bool

type IDFunc func(ctx context.Context, token string) (int64, error)



func Authenticate(idFunc IDFunc) func(http.Handler) http.Handler {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			token := r.Header.Get("Authorization")

			id, err := idFunc(r.Context(), token)
			if err != nil {
				http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}

			ctx := context.WithValue(r.Context(), authenticationContextKey, id)
			r = r.WithContext(ctx)

			handler.ServeHTTP(rw, r)

		})
	}

}

func Authentication(ctx context.Context) (int64, error) {
	if value, ok := ctx.Value(authenticationContextKey).(int64); ok {
		return value, nil 
	}
	return 0, ErrNoAuthentication
}