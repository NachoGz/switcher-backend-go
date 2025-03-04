package middleware

import (
	"context"
	"net/http"

	"github.com/NachoGz/switcher-backend-go/internal/database"
)

type contextKey string

const ConfigKey contextKey = "apiConfig"

type ApiConfig struct {
	DB *database.Queries
}

// WithConfig adds the API config to the request context
func WithConfig(cfg *ApiConfig) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), ConfigKey, cfg)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetConfig retrieves the API config from the request context
func GetConfig(r *http.Request) *ApiConfig {
	cfg, ok := r.Context().Value(ConfigKey).(*ApiConfig)
	if !ok {
		return nil
	}
	return cfg
}
