package handlers

import (
	"kura/db"
	"kura/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateUser(c *gin.Context) {
	var input models.User

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "不正なデータ形式です"})
		return
	}

	result := db.Conn.Create(&input)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ユーザーの作成に失敗しました"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "ユーザ登録完了"})
}
