package model

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"
)

// Entry entryテーブルの構造体
type Entry struct {
	ID        string    `db:"id" json:"-"`
	Url       string    `db:"url" json:"url"`
	Title     string    `db:"title" json:"title"`
	Caption   string    `db:"caption" json:"caption"`
	Thumbnail string    `db:"thumbnail" json:"thumbnail"`
	CreatedAt time.Time `db:"created_at" json:"-"`
}

// AddEntry 記事が存在しなければ記事を追加する
func AddEntry(ctx context.Context, entry *Entry) error {
	// すでに追加されているか確認
	// urlにuniqueが使えなかったので，とりあえずselect文で取ってくる方針にする
	var count int
	err := db.GetContext(ctx, &count, "SELECT COUNT(*) FROM entrys where url=?", entry.Url)
	if err != nil {
		return fmt.Errorf("failed to get entry: %w", err)
	}
	if count > 0 {
		return nil
	}

	// 記事の追加
	entryId := toHash(entry.Url)
	_, err = db.Exec("INSERT INTO entrys (id, url, title, caption, thumbnail) VALUES (?, ?, ?, ?, ?)", entryId, entry.Url, entry.Title, entry.Caption, entry.Thumbnail)
	if err != nil {
		return fmt.Errorf("failed to insert entry: %w", err)
	}
	return nil
}

// 文字列をハッシュ256化
func toHash(password string) string {
	converted := sha256.Sum256([]byte(password))
	return hex.EncodeToString(converted[:])
}
