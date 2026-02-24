package db

import (
	"fmt"
	"kura/models"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var Conn *gorm.DB

func Init() {
	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	port := "5432"
	timezone := os.Getenv("TZ")
	if timezone == "" {
		timezone = "Asia/Tokyo"
	}

	if (host == "") || (user == "") || (password == "") || (dbname == "") {
		panic("環境変数が設定されていません")
	}

	// DB接続
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=%s", host, user, password, dbname, port, timezone)

	var err error
	Conn, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("データベースに接続できませんでした" + err.Error())
	} else {
		fmt.Println("データベース接続成功")
	}

	// オートマイグレーション
	err = Conn.AutoMigrate(
		&models.AllowUser{},
		&models.User{},
		&models.File{},
		&models.UploadSession{},
	)
	if err != nil {
		panic("オートマイグレーションに失敗しました" + err.Error())
	} else {
		fmt.Println("オートマイグレーション成功")
	}
}
