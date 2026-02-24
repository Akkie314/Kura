package middleware

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 署名用の鍵を取得
		secretStr := os.Getenv("JWT_SECRET")
		if secretStr == "" {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "認証に失敗しました"})
			return
		}
		var jwtSecret = []byte(secretStr)

		// Authorizationヘッダを取得
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "認証トークンが必要です"})
			return
		}

		// Bearerプレフィックスを取り除く
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "トークン形式が不正です. Bearer形式である必要があります"})
			return
		}

		// トークンの検証
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("予期しない書名アルゴリズムです : %v", token.Header["alg"])
			}

			return jwtSecret, nil
		})

		// 検証結果の判定
		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "無効なトークンです"})
			return
		}

		// ペイロードからuserIDを取り出す
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			if userID, ok := claims["userID"].(string); ok {

				c.Set("userID", userID)
				c.Next()
			} else {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "トークンにユーザIDが含まれていません"})
				return
			}
		} else {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "認証に失敗しました"})
			return
		}
	}
}
