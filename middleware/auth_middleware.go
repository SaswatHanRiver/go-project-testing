package middleware

import (
	"net/http"
	"strings"

	"go-project-testing/utils"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware - equivalent to JwtAuthenticationFilter in Spring Boot
func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			utils.Error(ctx, http.StatusUnauthorized, "Authorization header is required")
			ctx.Abort()
			return
		}

		var tokenString string

		// Accept both formats:
		// 1. "Bearer eyJhbG..."  (standard format)
		// 2. "eyJhbG..."         (token only - Swagger UI sends this sometimes)
		if strings.HasPrefix(authHeader, "Bearer ") {
			tokenString = strings.TrimPrefix(authHeader, "Bearer ")
		} else {
			tokenString = authHeader
		}

		tokenString = strings.TrimSpace(tokenString)
		if tokenString == "" {
			utils.Error(ctx, http.StatusUnauthorized, "Token is missing")
			ctx.Abort()
			return
		}

		claims, err := utils.ValidateToken(tokenString)
		if err != nil {
			utils.Error(ctx, http.StatusUnauthorized, "Invalid or expired token")
			ctx.Abort()
			return
		}

		ctx.Set("userID", claims.UserID)
		ctx.Set("email", claims.Email)
		ctx.Next()
	}
}
