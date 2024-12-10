package weather

import (
	"github.com/BelovN/notifier/internal/repositories"
)

type Controller struct {
	userRepo repositories.UserRepository
}
