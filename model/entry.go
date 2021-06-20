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
	Count     int       `db:"number" json:"count"`
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
	err := db.GetContext(ctx, &entry, "SELECT entrys.* , COUNT(bookmarks.entry_id) AS number FROM entrys LEFT OUTER JOIN bookmarks ON entrys.id=bookmarks.entry_id WHERE entrys.id=? GROUP BY bookmarks.entry_id ", entryId)
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

// 週間記事ブックマーク件数が熱いやつを取得（10件）
func GetHotEntrys(ctx context.Context, userId string) ([]EntryDetail, error) {
	var entrys []Entry
	query := "SELECT entrys.* , COUNT(bookmarks.entry_id) AS number FROM entrys LEFT OUTER JOIN bookmarks ON entrys.id=bookmarks.entry_id GROUP BY bookmarks.entry_id ORDER BY number DESC LIMIT 6"
	t := time.Now()
	t.AddDate(0, 0, -7)
	t.Format("2006-01-02T15:04:05Z07:00")
	err := db.SelectContext(ctx, &entrys, query)
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
		entryDetails[i] = EntryDetail{
			Entry:      entry,
			Tags:       tags,
			IsBookmark: isBookmark,
		}
	}
	return entryDetails, nil
}

// 新着記事を50件取得
func GetNewEntrys(ctx context.Context, userId string) ([]EntryDetail, error) {
	// 2N+1を後で解決する！！！！！！TODO
	var entrys []Entry
	err := db.SelectContext(ctx, &entrys, "SELECT entrys.* , COUNT(bookmarks.entry_id) AS number FROM entrys LEFT OUTER JOIN bookmarks ON entrys.id=bookmarks.entry_id GROUP BY bookmarks.entry_id ORDER BY entrys.created_at DESC LIMIT 50")
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
		entryDetails[i] = EntryDetail{
			Entry:      entry,
			Tags:       tags,
			IsBookmark: isBookmark,
		}
	}
	return entryDetails, nil
}

// AddEntry 記事を追加する
func AddEntry(ctx context.Context, entry Entry) error {
	// 記事の追加
	if entry.ID == "" {
		return nil
	}
	_, err := db.Exec("INSERT INTO entrys (id, url, title, caption, thumbnail) SELECT ?, ?, ?, ?, ? FROM dual WHERE NOT EXISTS (SELECT * FROM entrys WHERE id=?)", entry.ID, entry.Url, entry.Title, entry.Caption, entry.Thumbnail, entry.ID)
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

// タグの削除
func DeleteTags(ctx context.Context, entryId string, tags []string) error {
	query := "DELETE FROM tags WHERE tag=? and entry_id=?"
	for _, tag := range tags {
		_, err := db.Exec(query, tag, entryId)
		if err != nil {
			return fmt.Errorf("failed to delete tags: %w", err)
		}
	}
	return nil
}

// Bookmarkの削除
func DeleteBookmark(ctx context.Context, entryId string, userId string) error {
	query := "DELETE FROM bookmarks WHERE entry_id=? and user_id=?"
	_, err := db.Exec(query, entryId, userId)
	if err != nil {
		return fmt.Errorf("failed to delete tags: %w", err)
	}
	return nil
}

func FindEntry(ctx context.Context, entryId string) (int, error) {
	var numEntrys int
	err := db.GetContext(ctx, &numEntrys, "SELECT count(*) FROM entrys WHERE id=?", entryId)
	if err != nil {
		return -1, fmt.Errorf("failed to get entry: %w", err)
	}

	return numEntrys, err
}

// 文字列をハッシュ256化
func ToHash(password string) string {
	converted := sha256.Sum256([]byte(password))
	return hex.EncodeToString(converted[:])
}
