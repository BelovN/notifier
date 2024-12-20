package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/BelovN/notifier/internal/models"
	_ "github.com/mattn/go-sqlite3"
	"strings"
)

type SqliteUserRepository struct {
	db *sql.DB
}

func NewSqliteUserRepository(databaseUrl string) (*SqliteUserRepository, error) {
	sqliteDatabase, err := sql.Open("sqlite3", databaseUrl)
	if err != nil {
		return nil, err
	}
	return &SqliteUserRepository{sqliteDatabase}, err
}

func (repo *SqliteUserRepository) FindByUsername(ctx context.Context, username string) (*models.User, error) {
	user := &models.User{}
	query := "SELECT id, username, is_subscribed FROM users WHERE username = ?"
	err := repo.db.QueryRowContext(ctx, query, username).Scan(&user.Id, &user.Username, &user.IsSubscribed)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (repo *SqliteUserRepository) Save(ctx context.Context, user *models.User) (*models.User, error) {
	query := "INSERT INTO users (username, is_subscribed, channel_id) VALUES (?, ?, ?)"
	result, err := repo.db.ExecContext(ctx, query, &user.Username, &user.IsSubscribed, &user.ChannelId)
	if err != nil {
		return nil, err
	}
	userId, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}
	user.Id = uint(userId)
	return user, nil
}

func (repo *SqliteUserRepository) Update(ctx context.Context, user *models.User) (*models.User, error) {
	query := "UPDATE users SET is_subscribed = ? WHERE username = ?"
	_, err := repo.db.ExecContext(ctx, query, &user.IsSubscribed, &user.Username)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (repo *SqliteUserRepository) GetOrCreateUser(ctx context.Context, username string, channelId int64) (*models.User, error) {
	user, err := repo.FindByUsername(ctx, username)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	if user != nil {
		return user, nil
	}
	newUser := &models.User{Username: username, ChannelId: channelId, IsSubscribed: true}
	return repo.Save(ctx, newUser)
}

func (repo *SqliteUserRepository) FilterUsers(ctx context.Context, filters map[string]interface{}) ([]*models.User, error) {
	query := "SELECT id, username, is_subscribed, channel_id FROM users"

	var whereClauses []string
	var args []interface{}

	for key, value := range filters {
		whereClauses = append(whereClauses, fmt.Sprintf("%s = ?", key))
		args = append(args, value)
	}

	if len(whereClauses) > 0 {
		query += " WHERE " + strings.Join(whereClauses, " AND ")
	}

	rows, err := repo.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		user := &models.User{}
		if err := rows.Scan(&user.Id, &user.Username, &user.IsSubscribed, &user.ChannelId); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}
