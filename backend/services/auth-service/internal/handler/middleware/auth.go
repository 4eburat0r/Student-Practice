package middleware

import (
    "net/http"
    "strings"

    "auth-service/pkg/jwt"

    "github.com/gin-gonic/gin"
)

func Auth(jwtManager *jwt.JWTManager) gin.HandlerFunc {
    return func(c *gin.Context) {
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            c.JSON(http.StatusUnauthorized, gin.H{
                "error": "Authorization header required",
            })
            c.Abort()
            return
        }

        parts := strings.SplitN(authHeader, " ", 2)
        if len(parts) != 2 || parts[0] != "Bearer" {
            c.JSON(http.StatusUnauthorized, gin.H{
                "error": "Invalid authorization header format",
            })
            c.Abort()
            return
        }

        token := parts[1]
        claims, err := jwtManager.ValidateAccessToken(token)
        if err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{
                "error": "Invalid or expired token",
            })
            c.Abort()
            return
        }

        c.Set("user_id", claims.UserID)
        c.Set("email", claims.Email)
        c.Set("role", claims.Role)

        c.Next()
    }
}

func GetUserID(c *gin.Context) (int, bool) {
    userID, exists := c.Get("user_id")
    if !exists {
        return 0, false
    }
    return userID.(int), true
}

func GetEmail(c *gin.Context) (string, bool) {
    email, exists := c.Get("email")
    if !exists {
        return "", false
    }
    return email.(string), true
}

func GetRole(c *gin.Context) (string, bool) {
    role, exists := c.Get("role")
    if !exists {
        return "", false
    }
    return role.(string), true
}
