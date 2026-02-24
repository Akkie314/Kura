package scheduler

import (
	"fmt"
	"kura/db"
	"kura/models"
	"os"
	"time"
)

func StartCleanupJob() {
	ticker := time.NewTicker(1 * time.Hour)

	go func() {
		runCleanup()

		for {
			<-ticker.C
			runCleanup()
		}
	}()
}

func runCleanup() {
	fmt.Println("[Scheduler] 期限切れデータのクリーンアップを開始します...")
	now := time.Now()

	cleanupExpiredFiles(now)
	cleanupExpiredSessions(now)

	fmt.Println("[Scheduler] クリーンアップが完了しました")
}

// 期限切れファイルの削除
func cleanupExpiredFiles(now time.Time) {
	var expiredFiles []models.File

	if err := db.Conn.Where("delete_at < ?", now).Find(&expiredFiles).Error; err != nil {
		fmt.Printf("[Scheduler] 期限切れファイルの取得に失敗しました: %v\n", err)
		return
	}

	for _, file := range expiredFiles {
		// ファイルを削除
		if err := os.Remove(file.FilePath); err != nil && !os.IsNotExist(err) {
			fmt.Printf("[Scheduler] 物理ファイルの削除に失敗しました(ID: %s): %v\n", file.ID, err)
			continue
		}

		// サムネイルも削除
		if file.HasThumbnail {
			thumbPath := file.FilePath + ".jpg"
			if err := os.Remove(thumbPath); err != nil && !os.IsNotExist(err) {
				fmt.Printf("[Scheduler] サムネイルの削除に失敗しました(ID: %s): %v\n", file.ID, err)
			}
		}

		// DBレコードの削除
		if err := db.Conn.Delete(&file).Error; err != nil {
			fmt.Printf("[Scheduler] DBレコードの削除に失敗しました(ID: %s): %v\n", file.ID, err)
		} else {
			fmt.Printf("[Scheduler] 期限切れファイルを削除しました(ID: %s)\n", file.ID)
		}
	}
}

// 期限切れセッションを削除
func cleanupExpiredSessions(now time.Time) {
	var expiredSessions []models.UploadSession

	// 期限切れのセッションを取得
	if err := db.Conn.Where("expires_at < ?", now).Find(&expiredSessions).Error; err != nil {
		fmt.Printf("[Scheduler] 期限切れセッションの取得に失敗しました: %v\n", err)
		return
	}

	for _, session := range expiredSessions {
		// 一時ファイルの削除
		if err := os.Remove(session.TempFilePath); err != nil && !os.IsNotExist(err) {
			fmt.Printf("[Scheduler] 一時ファイルの削除に失敗しました(UploadID: %s): %v\n", session.UploadID, err)
			continue
		}

		// DBレコードの削除
		if err := db.Conn.Delete(&session).Error; err != nil {
			fmt.Printf("[Scheduler] セッションレコードの削除に失敗しました(UploadID: %s): %v\n", session.UploadID, err)
		} else {
			fmt.Printf("[Scheduler] 期限切れセッションを削除しました(UploadID: %s)\n", session.UploadID)
		}
	}
}
