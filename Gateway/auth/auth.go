package auth

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(secret string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "havent token"})
			ctx.Abort()
			return
		}
		parts := strings.Split(authHeader, " ")
		//  authHeader у нас токен который мы взяли из запроса , а теперь мы делим
		// этот токен пробелом , потомучто он состоит из bearer dsclj2j2fj09v0jv03j03vr
		// разделим bearer и сам токен
		if len(parts) != 2 || parts[0] != "Bearer" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token format"})
			ctx.Abort()
			return
		}

		userID, role, err := ParseToken(parts[1],secret)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		ctx.Set("user_id", userID)
		ctx.Set("role", role)

		ctx.Next()
	}
}
func RBACMiddleware(permission string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		roleValue, exist := ctx.Get("role")
		if !exist {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		role, ok := roleValue.(string)
		if !ok {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		
		if HasPermission(role, permission) == false {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "have not permission"})
			return
		}
		ctx.Next()
	}
}
