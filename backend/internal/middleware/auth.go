package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/FallCatsinSeng/SWU_OSR/backend/internal/domain"
	"github.com/FallCatsinSeng/SWU_OSR/backend/internal/handler"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type contextKey string

const userClaimsKey contextKey = "user_claims"

// GetUserClaims extracts UserClaims from request context.
func GetUserClaims(ctx context.Context) (*domain.UserClaims, bool) {
	claims, ok := ctx.Value(userClaimsKey).(*domain.UserClaims)
	return claims, ok
}

// JWTAuth returns middleware that verifies JWT tokens from the Authorization
// header or the access_token cookie and injects UserClaims into context.
func JWTAuth(secret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenStr := extractToken(r)
			if tokenStr == "" {
				handler.RespondError(w, http.StatusUnauthorized, "missing or invalid token")
				return
			}

			token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
				if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, domain.ErrTokenInvalid
				}
				return []byte(secret), nil
			})
			if err != nil || !token.Valid {
				handler.RespondError(w, http.StatusUnauthorized, "invalid token")
				return
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				handler.RespondError(w, http.StatusUnauthorized, "invalid token claims")
				return
			}

			userIDStr, _ := claims["user_id"].(string)
			userID, err := uuid.Parse(userIDStr)
			if err != nil {
				handler.RespondError(w, http.StatusUnauthorized, "invalid user id in token")
				return
			}

			alias, _ := claims["alias"].(string)
			role, _ := claims["role"].(string)

			userClaims := &domain.UserClaims{
				UserID: userID,
				Alias:  alias,
				Role:   domain.Role(role),
			}

			ctx := context.WithValue(r.Context(), userClaimsKey, userClaims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func extractToken(r *http.Request) string {
	// Check Authorization header first
	authHeader := r.Header.Get("Authorization")
	if strings.HasPrefix(authHeader, "Bearer ") {
		return strings.TrimPrefix(authHeader, "Bearer ")
	}

	// Fallback to cookie
	cookie, err := r.Cookie("access_token")
	if err == nil {
		return cookie.Value
	}

	return ""
}
