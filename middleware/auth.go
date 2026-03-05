package middleware

import (
	"context"
	"net/http"
	"simple-product-api/utils"
	"strings"
)

// Happens first
func AuthenticateJWT(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization") //Bearer <Tokenstring>
		if authHeader == "" {
			http.Error(w, "Bearer String Empty", http.StatusUnauthorized)
			return
		}
		parts := strings.Split(authHeader, " ")
		if parts[0] != "Bearer" || len(parts) != 2 {
			http.Error(w, "Not JWT", http.StatusUnauthorized)
			return
		}
		token := parts[1] //<jwt Token>
		claims, err := utils.ParseToken(token)
		if err != nil {
			http.Error(w, "Failed Parsing", http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), utils.ClaimsKey, claims)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// happens after
func AuthenticateRole(role utils.Role) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			//sebenarnya agak redundant, claims invalid udh stopped di authJWT
			claims, ok := utils.GetClaimsFromContext(r.Context())
			if !ok {
				http.Error(w, "Claims Failed", http.StatusUnauthorized)
				return
			}

			//cek roles
			if claims.Role != string(role) {
				http.Error(w, "Incorrect Role Access", http.StatusForbidden)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
