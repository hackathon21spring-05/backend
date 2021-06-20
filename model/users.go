package model

import (
	"context"
	"fmt"
)

// User usersテーブルの構造体
type User struct {
	ID   string `db:"id" json:"id"`
	Name string `db:"name" json:"name"`
}

// CreateUser ユーザー作成
func CreateUser(ctx context.Context, user *User) error {
	_, err := db.ExecContext(ctx, "INSERT INTO users (id, name) VALUES (?, ?) ON DUPLICATE KEY UPDATE name = name", user.ID, user.Name)
	if err != nil {
		return fmt.Errorf("failed to insert user: %w", err)
	}

	return err
}

// GetBookmarks ブックマーク一覧を取得
func GetBookmarks(ctx context.Context, userId string) ([]EntryDetail, error) {
	// 2N+1を後で解決する！！！！！！TODO
	var entrys []Entry
	query := "SELECT q.* FROM bookmarks b " +
		"JOIN (SELECT bookmarks.created_at AS created_at, entrys.id AS id, url, caption, title, thumbnail, COUNT(bookmarks.entry_id) AS number FROM bookmarks JOIN entrys ON entrys.id=bookmarks.entry_id GROUP BY bookmarks.entry_id) q " +
		"ON b.entry_id=q.id WHERE b.user_id = ? ORDER BY q.created_at DESC"
	err := db.SelectContext(ctx, &entrys, query, userId)
	if err != nil {
		return []EntryDetail{}, fmt.Errorf("failed to get entry: %w", err)
	}

	// ここから完全にアウトじゃん
	entryDetails := make([]EntryDetail, len(entrys))
	for i, entry := range entrys {
		var tags []string
		err = db.SelectContext(ctx, &tags, "SELECT tag FROM tags WHERE entry_id=? ORDER BY tag", entry.ID)
		if err != nil {
			return []EntryDetail{}, fmt.Errorf("failed to get entry: %w", err)
		}
		isBookmark := true
		entryDetails[i] = EntryDetail{
			Entry:      entry,
			Tags:       tags,
			IsBookmark: isBookmark,
		}
	}
	return entryDetails, nil
}
