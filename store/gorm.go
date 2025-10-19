package store

import (
	"adong-be/models"
	"time"

	"github.com/hsdfat/go-auth-middleware/core"
	"gorm.io/gorm"
)

type Store struct {
	GormClient *gorm.DB
}

var DB *Store = &Store{}

// Enhanced UserProvider interface
type UserProvider interface {
	GetUserByUsername(username string) (*core.User, error)
	GetUserByID(userID int) (*core.User, error)
	GetUserByEmail(email string) (*core.User, error)
	UpdateUserLastLogin(userID int, lastLogin time.Time) error
	IsUserActive(userID int) (bool, error)
}

func (s *Store) GetUserByUsername(username string) (*core.User, error) {
	var dbUser models.User
	if err := s.GormClient.First(&dbUser, "userid = ?", username).Error; err != nil {
		return nil, err
	}
	return convertToCoreUser(dbUser), nil
}

func (s *Store) GetUserByID(userID int) (*core.User, error) {
	return nil, nil
}

func (s *Store) GetUserByEmail(email string) (*core.User, error) {
	var user models.User
	if err := s.GormClient.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return convertToCoreUser(user), nil
}

func (s *Store) UpdateUserLastLogin(userID int, lastLogin time.Time) error {
	return nil
}

func (s *Store) IsUserActive(userID int) (bool, error) {
	return true, nil
}

func convertToCoreUser(dbUser models.User) *core.User {
	return &core.User{
		Username: dbUser.UserID,
		Email:    dbUser.Email,
		Password: dbUser.Password,
		Role:     dbUser.Role,
		IsActive: true,
	}
}
