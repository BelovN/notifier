package repositories

import (
	"context"
	"github.com/BelovN/notifier/internal/models"
)

type UserRepository interface {
	Save(ctx context.Context, user *models.User) (*models.User, error)
	Update(ctx context.Context, user *models.User) (*models.User, error)
	FindByUsername(ctx context.Context, username string) (*models.User, error)
	GetOrCreateUser(ctx context.Context, username string, channelId int64) (*models.User, error)
}
