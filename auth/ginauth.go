package auth

import (
	"time"

	"github.com/hsdfat/go-auth-middleware/core"
)

// Enhanced UserProvider interface
type UserProvider interface {
	GetUserByUsername(username string) (*core.User, error)
	GetUserByID(userID int) (*core.User, error)
	GetUserByEmail(email string) (*core.User, error)
	UpdateUserLastLogin(userID int, lastLogin time.Time) error
	IsUserActive(userID int) (bool, error)
}

