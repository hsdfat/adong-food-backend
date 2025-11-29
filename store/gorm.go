package store

import (
	"adong-be/models"
	"time"

	"github.com/hsdfat/go-auth-middleware/core"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Store struct {
	GormClient *gorm.DB
}

var DB *Store = &Store{}

// Enhanced UserProvider interface
type UserProvider interface {
	GetUserByUsername(username string) (*core.User, error)
	GetUserByID(userID string) (*core.User, error)
	GetUserByEmail(email string) (*core.User, error)
	UpdateUserLastLogin(userID string, lastLogin time.Time) error
	IsUserActive(userID string) (bool, error)
}

// UserCreator interface for creating new users
type UserCreator interface {
	CreateUser(user *core.User) error
	UserExists(username string, email string) (bool, error)
	IsUsernameAvailable(username string) (bool, error)
	IsEmailAvailable(email string) (bool, error)
}

func (s *Store) GetUserByUsername(username string) (*core.User, error) {
	var dbUser models.User
	if err := s.GormClient.First(&dbUser, "user_name = ?", username).Error; err != nil {
		return nil, err
	}
	return convertToCoreUser(dbUser), nil
}

func (s *Store) GetUserByID(userID string) (*core.User, error) {
	var dbUser models.User
	if err := s.GormClient.First(&dbUser, "user_id = ?", userID).Error; err != nil {
		return nil, err
	}
	return convertToCoreUser(dbUser), nil
}

func (s *Store) GetUserByEmail(email string) (*core.User, error) {
	var user models.User
	if err := s.GormClient.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return convertToCoreUser(user), nil
}

func (s *Store) UpdateUserLastLogin(userID string, lastLogin time.Time) error {
	return nil
}

func (s *Store) IsUserActive(userID string) (bool, error) {
	return true, nil
}

func convertToCoreUser(dbUser models.User) *core.User {
	active := dbUser.Active != nil && *dbUser.Active
	return &core.User{
		ID:       dbUser.UserID,
		Username: dbUser.UserName,
		Email:    dbUser.Email,
		Password: dbUser.Password, // This is the hashed password
		Role:     dbUser.Role,
		IsActive: active,
		// Store plain password in a custom field if needed for verification
		// Note: The auth middleware will need to handle this
	}
}

// GetUserWithPlainPassword retrieves user with both hashed and plain password
func (s *Store) GetUserWithPlainPassword(username string) (*models.User, error) {
	var dbUser models.User
	if err := s.GormClient.First(&dbUser, "user_name = ?", username).Error; err != nil {
		return nil, err
	}
	return &dbUser, nil
}

// UserCreator implementation
func (s *Store) CreateUser(user *core.User) error {
	// Store plain password
	plainPassword := user.Password

	// Hash password if not already hashed
	hashedPassword := user.Password
	if len(user.Password) < 60 { // bcrypt hashes are 60 chars
		hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		hashedPassword = string(hash)
	}

	active := true
	dbUser := models.User{
		UserID:        user.ID,
		UserName:      user.Username,
		Password:      hashedPassword,
		PlainPassword: plainPassword,
		FullName:      user.Username, // Default to username if not provided
		Role:          user.Role,
		Email:         user.Email,
		Phone:         "",
		Active:        &active,
	}

	if err := s.GormClient.Create(&dbUser).Error; err != nil {
		return err
	}

	return nil
}

func (s *Store) UserExists(username string, email string) (bool, error) {
	var count int64
	err := s.GormClient.Model(&models.User{}).
		Where("user_name = ? OR email = ?", username, email).
		Count(&count).Error
	return count > 0, err
}

func (s *Store) IsUsernameAvailable(username string) (bool, error) {
	var count int64
	err := s.GormClient.Model(&models.User{}).
		Where("user_name = ?", username).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count == 0, nil
}

func (s *Store) IsEmailAvailable(email string) (bool, error) {
	if email == "" {
		return true, nil // Email is optional
	}
	var count int64
	err := s.GormClient.Model(&models.User{}).
		Where("email = ?", email).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count == 0, nil
}

// Helper function to hash passwords
func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// Helper function to verify passwords
// Supports both plain text and hashed password verification
func VerifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

// Helper function to verify password with fallback to plain text
// First tries bcrypt hash comparison, then falls back to plain text comparison
func VerifyPasswordWithPlainFallback(hashedPassword, plainPassword, inputPassword string) bool {
	// First try bcrypt verification
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(inputPassword)); err == nil {
		return true
	}

	// Fallback to plain text comparison if available
	if plainPassword != "" && plainPassword == inputPassword {
		return true
	}

	return false
}
