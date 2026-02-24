package main

import (
	"kura/db"
	"kura/handlers"
	"kura/middleware"
	"kura/scheduler"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

var DB *gorm.DB

func main() {
	// envを読み込む
	if err := godotenv.Load(); err != nil {
		panic(err)
	}

	// DB初期化
	db.Init()

	// 期限切れセッション / 期限切れファイルを削除
	scheduler.StartCleanupJob()

	// ルーティング
	r := gin.Default()

	r.SetTrustedProxies([]string{
		"127.0.0.1",      // localhost
		"10.0.0.0/8",     // プライベートIP帯A
		"172.16.0.0/12",  // プライベートIP帯B (Dockerは主にここ)
		"192.168.0.0/16", // プライベートIP帯C
	})

	// CORS設定
	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"http://localhost:3000",
			"https://kura.akkie.me",
		},
		AllowMethods: []string{
			"GET", "POST", "PATCH",
		},
		AllowHeaders: []string{
			"Origin",
			"Content-Type",
			"Authorization",
		},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	auth := r.Group("/auth")
	{
		auth.GET("/login", handlers.DiscordLogin)
		auth.GET("/callback", handlers.DiscordCallback)
	}

	// 認証が不要
	public := r.Group("/")
	{
		public.GET("/files/:id/meta", handlers.GetMeta)
		public.GET("/files/:id/thumb", handlers.DownloadThumbnail)
	}

	// 認証が必要
	api := r.Group("/")
	api.Use(middleware.AuthMiddleware())
	{
		api.POST("/users", handlers.CreateUser)
		api.POST("/upload/session", handlers.StartUploadSession)
		api.PATCH("/upload/:upload_id", handlers.AppendChunk)
		api.POST("/upload/:upload_id/complete", handlers.CompleteUpload)
	}
	r.GET("/files/:id", handlers.DownloadFile)

	r.Run()
}
