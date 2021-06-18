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
	ID        string    `db:"id" json:"id"`
	Url       string    `db:"url" json:"url"`
	Title     string    `db:"title" json:"title"`
	Caption   string    `db:"caption" json:"caption"`
	Thumbnail string    `db:"thumbnail" json:"thumbnail"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

// EntryDetail entry詳細を返す
type EntryDetail struct {
	Entry
	Tags       []string `db:"tags" json:"tags"`
	IsBookmark bool     `json:"isBookmark"`
}

// getEntryDetail 特定の記事詳細を取得する
func GetEntryDetail(ctx context.Context, userId string, entryId string) (EntryDetail, error) {
	// 記事の本体情報を取得
	entry := Entry{}
	err := db.GetContext(ctx, &entry, "SELECT * FROM entrys WHERE id=?", entryId)
	if err != nil {
		return EntryDetail{}, fmt.Errorf("failed to get entry: %w", err)
	}

	// タグリストを取得
	var tags []string
	err = db.SelectContext(ctx, &tags, "SELECT tag FROM tags WHERE entry_id=?", entryId)
	if err != nil {
		return EntryDetail{}, fmt.Errorf("failed to get tags: %w", err)
	}

	// ブックマークをしているのか取得
	var count int
	isBookmark := false
	err = db.GetContext(ctx, &count, "SELECT COUNT(*) FROM bookmarks WHERE user_id=? and entry_id=?", userId, entryId)
	if err != nil {
		return EntryDetail{}, fmt.Errorf("failed to get bookmarks: %w", err)
	}
	if count > 0 {
		isBookmark = true
	}

	return EntryDetail{
		Entry:      entry,
		Tags:       tags,
		IsBookmark: isBookmark,
	}, nil
}

// 新着記事を50件取得
func GetNewEntrys(ctx context.Context, userId string) ([]EntryDetail, error) {
	// 2N+1を後で解決する！！！！！！TODO
	var entrys []Entry
	err := db.SelectContext(ctx, &entrys, "SELECT * FROM entrys ORDER BY created_at DESC LIMIT 50")
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
		var count int
		err = db.GetContext(ctx, &count, "SELECT COUNT(*) FROM bookmarks WHERE user_id=? and entry_id=?", userId, entry.ID)
		if err != nil {
			return []EntryDetail{}, fmt.Errorf("failed to get entry: %w", err)
		}
		isBookmark := false
		if count > 0 {
			isBookmark = true
		}
		entryDetails[i].ID = entry.ID
		entryDetails[i].Url = entry.Url
		entryDetails[i].Title = entry.Title
		entryDetails[i].Caption = entry.Caption
		entryDetails[i].Thumbnail = entry.Thumbnail
		entryDetails[i].CreatedAt = entry.CreatedAt

		entryDetails[i].Tags = tags
		entryDetails[i].IsBookmark = isBookmark
	}
	return entryDetails, nil
}

// AddEntry 記事が存在しなければ記事を追加する
func AddEntry(ctx context.Context, entry *Entry) error {
	// すでに追加されているか確認
	// urlにuniqueが使えなかったので，とりあえずselect文で取ってくる方針にする
	var count int
	entryId := ToHash(entry.Url)
	err := db.GetContext(ctx, &count, "SELECT COUNT(*) FROM entrys where id=?", entryId)
	if err != nil {
		return fmt.Errorf("failed to get entry: %w", err)
	}
	if count > 0 {
		return nil
	}
	// URLからデータをとってくる
	entryData, err := getEntryContent(entry.Url)
	if err != nil {
		return fmt.Errorf("failed to get entry Data: %w", err)
	}

	// 記事の追加
	_, err = db.Exec("INSERT INTO entrys (id, url, title, caption, thumbnail) VALUES (?, ?, ?, ?, ?)", entryId, entry.Url, entryData.Title, entryData.Caption, entryData.Thumbnail)
	if err != nil {
		return fmt.Errorf("failed to insert entry: %w", err)
	}
	return nil
}

// ブックマークに追加
func AddBookMark(ctx context.Context, userId string, entryId string) error {
	query := "INSERT INTO bookmarks (user_id, entry_id) SELECT ?, ? FROM dual WHERE NOT EXISTS (SELECT * FROM bookmarks WHERE user_id=? and entry_id=? )"
	_, err := db.ExecContext(ctx, query, userId, entryId, userId, entryId)
	if err != nil {
		return fmt.Errorf("failed to insert bookmarks: %w", err)
	}
	return nil
}

// タグの追加
func AddTags(ctx context.Context, entryId string, tags []string) error {
	query := "INSERT INTO tags (tag, entry_id) SELECT ?, ? FROM dual WHERE NOT EXISTS (SELECT * FROM tags WHERE tag=? and entry_id=?)"
	// Bulk Insertできなかった……助けて！！！
	for _, tag := range tags {
		_, err := db.Exec(query, tag, entryId, tag, entryId)
		if err != nil {
			return fmt.Errorf("failed to insert tags: %w", err)
		}
	}
	return nil
}

func FindEntry(ctx context.Context, entryId string) (numentrys int, err error) {

	err = db.GetContext(ctx, &numentrys, "SELECT count(*) FROM entrys WHERE entryId=?", entryId)
	if err != nil {
		return -1, fmt.Errorf("failed to get entry: %w", err)
	}

	return numentrys, err
}

// 文字列をハッシュ256化
func ToHash(password string) string {
	converted := sha256.Sum256([]byte(password))
	return hex.EncodeToString(converted[:])
}
