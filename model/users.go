package model

import (
	"context"
	"fmt"
)

// User usersテーブルの構造体
type User struct {
	ID   string `db:"user_id" json:"user_id"`
	Name string `db:"name" json:"name"`
}

// CreateUser ユーザー作成
func CreateUser(ctx context.Context, user *User) error {
	_, err := db.ExecContext(ctx, "INSERT INTO users (user_id, name) VALUES (?, ?) ON DUPLICATE KEY UPDATE name = name", user.ID, user.Name)
	if err != nil {
		return fmt.Errorf("failed to insert user: %w", err)
	}

	return err
}
