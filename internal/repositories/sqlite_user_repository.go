package repositories

import (
	"context"
	"database/sql"
	"errors"
	"github.com/BelovN/notifier/internal/models"
	_ "github.com/mattn/go-sqlite3"
)

type SqliteUserRepository struct {
	db *sql.DB
}

func NewSqliteUserRepository(databaseUrl string) (error, *SqliteUserRepository) {
	sqliteDatabase, err := sql.Open("sqlite3", databaseUrl)
	if err != nil {
		return err, nil
	}
	return nil, &SqliteUserRepository{sqliteDatabase}
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
