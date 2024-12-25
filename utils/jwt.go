package utils

import (
    "github.com/dgrijalva/jwt-go"
    "time"
    "ecommerce-api/models"
    "gorm.io/gorm"
    "net/http"
    "github.com/gin-gonic/gin"
)

var jwtSecret = []byte("your-secret-key")

// GenerateToken generates a JWT token for a given user ID and role
func GenerateToken(userId uint, role string) (string, error) {
    claims := jwt.MapClaims{
        "user_id": userId,
        "role":    role,
        "exp":     time.Now().Add(time.Hour * 24).Unix(),
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(jwtSecret)
}

// ParseToken parses the JWT token and returns the claims
func ParseToken(tokenStr string) (jwt.MapClaims, error) {
    token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
        return jwtSecret, nil
    })
    if err != nil || !token.Valid {
        return nil, err
    }
    return token.Claims.(jwt.MapClaims), nil
}

// AuthMiddleware is a Gin middleware for validating JWT tokens
func AuthMiddleware(db *gorm.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        token := c.GetHeader("Authorization")
        if token == "" {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token required"})
            c.Abort()
            return
        }

        claims, err := ParseToken(token)
        if err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
            c.Abort()
            return
        }

        userId := uint(claims["user_id"].(float64)) // type assertion from float64 to uint
        var user models.User
        if err := db.First(&user, userId).Error; err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
            c.Abort()
            return
        }

        c.Set("user", user) // Attach the user to the context
        c.Next()
    }
}
