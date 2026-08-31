package middleware

import (
	"net/http"
	"strings"

	"go-project-testing/utils"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware - equivalent to JwtAuthenticationFilter extends OncePerRequestFilter
// In Spring Boot this is registered in SecurityConfig.filterChain()
// In Gin we apply it to a route group: router.Group("/api").Use(AuthMiddleware())
func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Read the Authorization header - "Bearer <token>"
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization header is required",
			})
			return
		}

		// Header must be "Bearer <token>" format
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization header format must be: Bearer <token>",
			})
			return
		}

		tokenString := parts[1]

		// Validate token - like jwtUtil.validateToken() in Spring Boot filter
		claims, err := utils.ValidateToken(tokenString)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid or expired token",
			})
			return
		}

		// Store user info in context - like SecurityContextHolder.getContext().setAuthentication()
		// Controllers can read this with: ctx.GetUint("userID")
		ctx.Set("userID", claims.UserID)
		ctx.Set("email", claims.Email)

		ctx.Next() // continue to the actual handler
	}
}
