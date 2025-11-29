package auth

import (
	"adong-be/store"
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/hsdfat/go-auth-middleware/core"
	"golang.org/x/crypto/bcrypt"
)

// LoginRequest represents the login credentials
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// CreateDualPasswordAuthenticator creates an authenticator that supports both hashed and plain text passwords
// It first tries bcrypt verification, then falls back to plain text comparison
func CreateDualPasswordAuthenticator(db *store.Store) func(c *gin.Context) (*core.User, error) {
	return func(c *gin.Context) (*core.User, error) {
		var loginReq LoginRequest
		if err := c.ShouldBindJSON(&loginReq); err != nil {
			return nil, errors.New("invalid request format")
		}

		// Get user with both passwords
		dbUser, err := db.GetUserWithPlainPassword(loginReq.Username)
		if err != nil {
			return nil, errors.New("invalid username or password")
		}

		// Check if user is active
		if dbUser.Active == nil || !*dbUser.Active {
			return nil, errors.New("user account is inactive")
		}

		// Try bcrypt verification first
		err = bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(loginReq.Password))
		if err == nil {
			// Bcrypt verification successful
			return &core.User{
				ID:       dbUser.UserID,
				Username: dbUser.UserName,
				Email:    dbUser.Email,
				Password: dbUser.Password,
				Role:     dbUser.Role,
				IsActive: *dbUser.Active,
			}, nil
		}

		// Fallback to plain text comparison if bcrypt fails
		if dbUser.PlainPassword != "" && dbUser.PlainPassword == loginReq.Password {
			// Plain text verification successful
			return &core.User{
				ID:       dbUser.UserID,
				Username: dbUser.UserName,
				Email:    dbUser.Email,
				Password: dbUser.Password,
				Role:     dbUser.Role,
				IsActive: *dbUser.Active,
			}, nil
		}

		// Both verifications failed
		return nil, errors.New("invalid username or password")
	}
}
