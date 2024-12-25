package utils

import (
    "github.com/dgrijalva/jwt-go"
    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
    "net/http"
    "strings"
    "time"
    "ecommerce-api/models"
)

var jwtSecret = []byte("secret-key") 

func GenerateToken(userId uint, role string) (string, error) {
    claims := jwt.MapClaims{
        "user_id": userId,
        "role":    role,
        "exp":     time.Now().Add(time.Hour * 24).Unix(),
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(jwtSecret)
}

func ParseToken(tokenStr string) (jwt.MapClaims, error) {
    token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, jwt.NewValidationError("unexpected signing method", jwt.ValidationErrorSignatureInvalid)
        }
        return jwtSecret, nil
    })
    if err != nil || !token.Valid {
        return nil, err
    }
    claims, ok := token.Claims.(jwt.MapClaims)
    if !ok {
        return nil, jwt.NewValidationError("invalid claims format", jwt.ValidationErrorClaimsInvalid)
    }
    return claims, nil
}

func AuthMiddleware(db *gorm.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        token := c.GetHeader("Authorization")
        if token == "" {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token required"})
            c.Abort()
            return
        }

        tokenParts := strings.Split(token, " ")
        if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Authorization header format"})
            c.Abort()
            return
        }
        token = tokenParts[1]

        claims, err := ParseToken(token)
        if err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
            c.Abort()
            return
        }

        userId, ok := claims["user_id"].(float64) 
        if !ok {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token payload"})
            c.Abort()
            return
        }

        var user models.User
        if err := db.First(&user, uint(userId)).Error; err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
            c.Abort()
            return
        }

        c.Set("user", user)
        c.Next()
    }
}

func IsAdminMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        user, exists := c.Get("user")
        if !exists {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
            c.Abort()
            return
        }

        userModel, ok := user.(models.User)
        if !ok || userModel.Role != "admin" {
            c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
            c.Abort()
            return
        }

        c.Next()
    }
}