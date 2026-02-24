package handlers

import (
	"fmt"
	"io"
	"kura/db"
	"kura/models"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func StartUploadSession(c *gin.Context) {
	userID := c.MustGet("userID").(string)

	var req models.SessionRequest

	// JSON受け取る
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "不正なデータ形式です"})
		return
	}

	// 保存場所を作成
	tempUUID := uuid.New().String()
	tempFileName := "temp_" + tempUUID
	tempFilePath := filepath.Join("uploads", tempFileName)

	// ファイルを作成
	file, err := os.Create(tempFilePath)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ファイルアップロードセッション開始に失敗しました"})
		return
	}
	file.Close()

	// DBに登録
	uploadID := uuid.New().String()

	session := models.UploadSession{
		UploadID:     uploadID,
		UserID:       userID,
		OriginalName: req.FileName,
		MimeType:     req.FileType,
		TotalSize:    req.FileSize,
		StoredSize:   0,
		TempFilePath: tempFilePath,
		ExpiresAt:    time.Now().Add(24 * time.Hour),
	}

	if result := db.Conn.Create(&session); result.Error != nil {
		os.Remove(tempFilePath)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ファイルアップロードセッション開始に失敗しました"})
		return
	}

	// チケットをクライアントに渡す
	c.JSON(http.StatusOK, gin.H{
		"message":    "アップロードセッションを開始しました",
		"upload_id":  uploadID,
		"chunk_size": 1024 * 1024 * 50,
	})
}

func AppendChunk(c *gin.Context) {
	userID := c.MustGet("userID").(string)

	const MaxChunkSize = 50 * 1024 * 1024 // 50MB
	incomingSize := c.Request.ContentLength

	if incomingSize == -1 || incomingSize > MaxChunkSize {
		c.JSON(http.StatusBadRequest, gin.H{"error": "アップロードに失敗しました"})
		return
	}

	// アップロードIDを取得
	uploadID := c.Param("upload_id")

	// セッションが有効か
	var session models.UploadSession
	if err := db.Conn.First(&session, "upload_id = ?", uploadID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "セッションが見つかりません"})
		return
	}

	// セッションを持っているか
	if session.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "セッションが見つかりません"})
		return
	}

	// 有効期限
	if time.Now().After(session.ExpiresAt) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "有効期限が切れています"})
		return
	}

	// サイズチェック
	if session.StoredSize >= session.TotalSize {
		c.JSON(http.StatusBadRequest, gin.H{"error": "すでにアップロードが完了しています"})
		return
	}

	// 送られたデータを読み込む
	body := c.Request.Body
	defer body.Close()

	// 一次ファイルを開く
	file, err := os.OpenFile(session.TempFilePath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "アップロードに失敗しました"})
		return
	}

	// 書き込む前のファイルサイズを記録
	tempFileSize := session.StoredSize

	// データを書き込む
	written, err := io.Copy(file, body)
	file.Close()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "アップロードに失敗しました"})
		return
	}

	// StoredSizeを更新
	session.StoredSize += int64(written)

	// DB更新
	if err := db.Conn.Save(&session).Error; err != nil {
		// DB更新に失敗するとDBと実際のファイルでギャップが生まれるので, 元のサイズに切り落とす
		os.Truncate(session.TempFilePath, tempFileSize)

		c.JSON(http.StatusInternalServerError, gin.H{"error": "アップロードに失敗しました"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "正常にチャンクをアップロードできました",
		"StoredSize": session.StoredSize,
		"TotalSize":  session.TotalSize,
	})
}

func CompleteUpload(c *gin.Context) {
	userID := c.MustGet("userID").(string)

	uploadID := c.Param("upload_id")
	// セッションを取得
	var session models.UploadSession
	if err := db.Conn.First(&session, "upload_id = ?", uploadID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "セッションが見つかりません"})
		return
	}

	// セッションを持っているか
	if session.UserID != userID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "セッションが見つかりません"})
		return
	}

	// アップロードが完了しているか
	if session.StoredSize != session.TotalSize {
		c.JSON(http.StatusBadRequest, gin.H{"error": "アップロードが完了していません"})
		return
	}

	// 正式なファイルを作成
	fileID := uuid.New().String()
	filePath := filepath.Join("uploads", fileID)

	// MimeTypeを検証
	mimeType, err := detectMimeType(session.TempFilePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ファイルの保存中にエラーが発生しました"})
		return
	}

	// 動画であればサムネイルを生成
	hasThumbnail := false
	if strings.HasPrefix(mimeType, "video/") {
		thumbPath := filePath + ".jpg"

		/*
		 * ffmpeg
		 * -y 		... 上書き許可
		 * -ss 1 	... 開始1秒
		 * -i 		... 入力ファイル
		 * -vframes 1... 1フレームだけ
		 * -f image2... 画像として出力
		 */
		cmd := exec.Command("ffmpeg", "-y", "-ss", "1", "-i", session.TempFilePath, "-vframes", "1", "-f", "image2", thumbPath)

		if err := cmd.Run(); err != nil {
			hasThumbnail = false
			fmt.Println("サムネイル作成に失敗", err)
		} else {
			hasThumbnail = true
		}
	}

	// application/octet-streamならユーザの宣言を信じる
	if mimeType == "application/octet-stream" {
		mimeType = session.MimeType
	}

	// ハッシュ計算
	hash, err := calcHash(session.TempFilePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ファイルの保存中にエラーが発生しました"})
		return
	}

	// リネーム
	if err := os.Rename(session.TempFilePath, filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "完了処理に失敗しました"})
		return
	}

	// DBに登録
	fileRecord := models.File{
		ID:           fileID,
		UserID:       userID,
		OriginalName: session.OriginalName,
		MimeType:     mimeType,
		HasThumbnail: hasThumbnail,
		Size:         session.TotalSize,
		FilePath:     filePath,
		HashSha256:   hash,
		DeleteAt:     time.Now().Add(30 * 24 * time.Hour),
	}

	if result := db.Conn.Create(&fileRecord); result.Error != nil {
		os.Rename(filePath, session.TempFilePath)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "完了処理に失敗しました"})
		return
	}

	// セッションを削除
	if err := db.Conn.Delete(&session).Error; err != nil {
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "正常にアップロード完了しました",
		"file_id": fileID,
	})
}
