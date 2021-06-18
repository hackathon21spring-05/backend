package model

import (
	"context"
	"fmt"
)

func SearchEntrys(ctx context.Context, tag string, userId string) ([]EntryDetail, error) {
	query := "SELECT entrys.* FROM entrys, tags WHERE entrys.id = tags.entry_id AND (tags.tag IN (?)) GROUP BY entrys.Id ORDER BY created_at DESC"
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
