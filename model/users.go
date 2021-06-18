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

type UserBm []ListBm

type ListBm []string

// CreateUser ユーザー作成
func CreateUser(ctx context.Context, user *User) error {
	_, err := db.ExecContext(ctx, "INSERT INTO users (id, name) VALUES (?, ?) ON DUPLICATE KEY UPDATE name = name", user.ID, user.Name)
	if err != nil {
		return fmt.Errorf("failed to insert user: %w", err)
	}

	return err
}

func GetBookmark(ctx context.Context, user *User) (UserBm, error) {

	// UserBm ⊂ ListBm := []string
	var userbookmarks UserBm
	//var listbookmark []ListBm

	//user_id=user.idとなる要素のentry_idをすべてrowsに代入
	rows, err := db.Query("SELECT entry_id FROM bookmarks where user_id=?", user.ID)
	if err != nil {
		return userbookmarks, fmt.Errorf("no bookmark exists: %w", err)
	}

	for rows.Next() {
		if err := rows.Scan(&userbookmarks); err != nil {
			return userbookmarks, fmt.Errorf("fatal error: %w", err)
		}
	}

	return userbookmarks, err
}
