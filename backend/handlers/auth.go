package handlers

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"kura/db"
	"kura/models"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"golang.org/x/oauth2"
)

func getDiscordConfig() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     os.Getenv("DISCORD_CLIENT_ID"),
		ClientSecret: os.Getenv("DISCORD_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("DISCORD_REDIRECT_URL"),
		Scopes:       []string{"identify"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://discord.com/api/oauth2/authorize",
			TokenURL: "https://discord.com/api/oauth2/token",
		},
	}
}

// ログインページへリダイレクト
func DiscordLogin(c *gin.Context) {
	state, err := generateRandomString(32)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ログイン中に失敗しました"})
		return
	}

	// stateをクッキーに保存
	/* Name 		: oauth_state
	 * Value		: state
	 * Expiration 	: 10分
	 * Path 		: /
	 * Domain		: localhost
	 * Secure		: true
	 * HttpOnly		: true(JavaScriptからアクセスできなくなるそう)
	 */
	c.SetCookie("oauth_state", state, 600, "/", "", true, true)

	config := getDiscordConfig()
	url := config.AuthCodeURL(state)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func DiscordCallback(c *gin.Context) {
	// Cookieにあるoauth_stateを取得
	oauthState, err := c.Cookie("oauth_state")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "認証に失敗しました"})
		return
	}

	// URLパラメータのstateと照合
	if c.Query("state") != oauthState {
		c.JSON(http.StatusBadRequest, gin.H{"error": "認証に失敗しました"})
		return
	}

	// Cookieを削除
	c.SetCookie("oauth_state", "", -1, "/", "", true, true)

	// 認証コードを取得
	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "認証に失敗しました"})
		return
	}

	config := getDiscordConfig()

	// コードをアクセストークンに変換
	token, err := config.Exchange(context.Background(), code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "認証に失敗しました"})
		return
	}

	// Discordのユーザ情報を取得
	client := config.Client(context.Background(), token)
	resp, err := client.Get("https://discord.com/api/users/@me")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "認証に失敗しました"})
		return
	}
	defer resp.Body.Close()

	var discordUser models.DiscordUser
	if err := json.NewDecoder(resp.Body).Decode(&discordUser); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "認証中にエラーが発生しました"})
		return
	}

	// DB処理
	var user models.User
	if err := db.Conn.First(&user, "discord_id = ?", discordUser.ID).Error; err != nil {
		user = models.User{
			DiscordID: discordUser.ID,
			Username:  discordUser.Username,
		}

		if err := db.Conn.Create(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "認証中にエラーが発生しました"})
			return
		}
	} else {
		user.Username = discordUser.Username
		db.Conn.Save(&user)
	}

	// JWT生成
	jwtSecretStr := os.Getenv("JWT_SECRET")
	if jwtSecretStr == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "認証中にエラーが発生しました"})
		return
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID":   user.ID,
		"username": user.Username,
		"exp":      time.Now().Add(72 * time.Hour).Unix(),
	})

	tokenString, err := jwtToken.SignedString([]byte(jwtSecretStr))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "認証中にエラーが発生しました"})
		return
	}

	// フロントエンドにリダイレクト
	frontendURL := os.Getenv("FRONTEND_URL")
	c.SetCookie("token", tokenString, 3600*72, "/", "", false, false)
	c.Redirect(http.StatusTemporaryRedirect, frontendURL)
}

// ランダムな文字列を生成する
func generateRandomString(length int) (string, error) {
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}
