package handlers

import (
	"kura/db"
	"kura/models"
	"net/http"
	"net/url"
	"os"

	"github.com/gin-gonic/gin"
)

// func UploadFile(c *gin.Context) {
// 	// ファイルIDを生成
// 	fileID := uuid.New().String()

// 	// ファイル取得
// 	file, err := c.FormFile("file")
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "不正なデータ形式です"})
// 		return
// 	}

// 	// ユーザID取得
// 	userID := c.PostForm("user_id")
// 	if userID == "" {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "ユーザIDが必要です"})
// 		return
// 	}
// 	// ユーザIDが正しいか見る
// 	var user models.User
// 	if err := db.Conn.First(&user, "id = ?", userID).Error; err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "ユーザIDが不正です"})
// 		return
// 	}

// 	// ファイル名と保存先を決定
// 	filename := uuid.New().String()
// 	savePath := filepath.Join("uploads", filename)

// 	// ファイルを保存
// 	if err := c.SaveUploadedFile(file, savePath); err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "ファイルの保存中にエラーが発生しました"})
// 		return
// 	}

// 	// MimeTypeを取得
// 	mimeType, err := detectMimeType(savePath)
// 	if err != nil {
// 		os.Remove(savePath)
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "ファイルの保存中にエラーが発生しました"})
// 		return
// 	}

// 	// ハッシュ値計算
// 	hash, err := calcHash(savePath)
// 	if err != nil {
// 		os.Remove(savePath)
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "ファイルの保存中にエラーが発生しました"})
// 		return
// 	}

// 	// データベースに記録する内容
// 	fileRecord := models.File{
// 		ID:           fileID,
// 		UserID:       userID,
// 		OriginalName: file.Filename,
// 		MimeType:     mimeType,
// 		Size:         file.Size,
// 		FilePath:     savePath,
// 		HashSha256:   hash,
// 		DeleteAt:     time.Now().Add(30 * 24 * time.Hour),
// 	}

// 	// データベースに保存
// 	result := db.Conn.Create(&fileRecord)
// 	if result.Error != nil {
// 		os.Remove(savePath) // 保存していたファイルを削除する

// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "ファイルの保存中にエラーが発生しました"})
// 		return
// 	}

// 	// 成功レスポンス
// 	c.JSON(http.StatusOK, gin.H{"message": "ファイルのアップロードに成功しました", "file_id": fileID})
// }

func DownloadFile(c *gin.Context) {
	// idを取得
	id := c.Param("id")

	// DBからファイル情報を検索
	var fileData models.File
	if err := db.Conn.First(&fileData, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "ファイルが見つかりませんでした"})
		return
	}

	// 実際にディスクにあるか確認
	if _, err := os.Stat(fileData.FilePath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "ファイルが見つかりませんでした"})
		return
	}

	// レスポンスヘッダを作成
	c.Header("Content-Type", fileData.MimeType)

	fileName := url.QueryEscape(fileData.OriginalName)
	c.Header("Content-Disposition", "inline; filename*=UTF-8''"+fileName)

	// ファイルをレスポンス
	c.File(fileData.FilePath)
}

func DownloadThumbnail(c *gin.Context) {
	// idを取得
	id := c.Param("id")

	// DBからファイル情報を検索
	var fileData models.File
	if err := db.Conn.First(&fileData, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "ファイルが見つかりませんでした"})
		return
	}

	// 実際にディスクにあるか確認
	if _, err := os.Stat(fileData.FilePath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "ファイルが見つかりませんでした"})
		return
	}

	// サムネを持つかを確認
	if !fileData.HasThumbnail {
		c.JSON(http.StatusNotFound, gin.H{"error": "サムネイルが見つかりませんでした"})
		return
	}

	// ファイルパスを作成
	thumbPath := fileData.FilePath + ".jpg"

	// ファイルが存在したら返す
	if _, err := os.Stat(thumbPath); err == nil {
		c.File(thumbPath)
	} else {
		c.JSON(http.StatusNotFound, gin.H{"error": "サムネイルが見つかりませんでした"})
	}
}

func GetMeta(c *gin.Context) {
	id := c.Param("id")

	var file models.File
	// DBから検索
	if err := db.Conn.First(&file, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "ファイルが見つかりませんでした"})
		return
	}

	var res = models.FileMeta{
		ID:           file.ID,
		OriginalName: file.OriginalName,
		MimeType:     file.MimeType,
		HasThumbnail: file.HasThumbnail,
	}

	c.JSON(http.StatusOK, res)
}
