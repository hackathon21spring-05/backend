package model

import (
	"context"
	"fmt"
)

func SearchEntrys(ctx context.Context, tag string, userId string) ([]EntryDetail, error) {
	query := "SELECT q.*, COUNT(b.entry_id) AS number FROM bookmarks b " +
		"RIGHT OUTER JOIN (SELECT entrys.* FROM entrys JOIN tags WHERE entrys.id = tags.entry_id AND (tags.tag IN (?)) GROUP BY entrys.id) q " +
		"ON b.entry_id=q.id GROUP BY b.entry_id"
	var entrys []Entry
	err := db.SelectContext(ctx, &entrys, query, tag)
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
