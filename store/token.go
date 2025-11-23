package store

import (
	"adong-be/logger"
	"adong-be/models"
	"time"

	"github.com/hsdfat/go-auth-middleware/core"
)

var TokenInterface core.TokenStorage

// Enhanced TokenStorage interface for refresh token support
type TokenStorage interface {
	// Access token methods
	StoreTokenPair(sessionID string, accessToken, refreshToken string, accessExpiresAt, refreshExpiresAt time.Time, userID string) error
	GetAccessToken(sessionID string) (string, error)
	GetRefreshToken(sessionID string) (string, error)

	// Token validation
	IsAccessTokenValid(sessionID string) (bool, error)
	IsRefreshTokenValid(sessionID string) (bool, error)

	// Token management
	DeleteTokenPair(sessionID string) error
	RefreshTokenPair(sessionID string, newAccessToken, newRefreshToken string, accessExpiresAt, refreshExpiresAt time.Time) error

	// User session management
	RevokeAllUserTokens(userID string) error
	GetUserActiveSessions(userID string) ([]string, error)

	// Session tracking
	StoreUserSession(session core.UserSession) error
	GetUserSession(sessionID string) (*core.UserSession, error)
	UpdateSessionActivity(sessionID string, lastActivity time.Time) error
	DeleteUserSession(sessionID string) error

	// Cleanup expired tokens
	CleanupExpiredTokens() error
}

func NewTokenInterface() core.TokenStorage {
	return TokenInterface
}

func SetTokenInterface(tokenInterface core.TokenStorage) {
	TokenInterface = tokenInterface
}

func (s *Store) StoreTokenPair(sessionID string, accessToken, refreshToken string, accessExpiresAt, refreshExpiresAt time.Time, userID string) error {
	return s.GormClient.Create(&models.TokenPair{
		SessionID:        sessionID,
		AccessToken:      accessToken,
		RefreshToken:     refreshToken,
		AccessExpiresAt:  accessExpiresAt,
		RefreshExpiresAt: refreshExpiresAt,
		UserID:           userID,
	}).Error
}

func (s *Store) GetAccessToken(sessionID string) (string, error) {
	var tokenPair models.TokenPair
	if err := s.GormClient.First(&tokenPair, "session_id = ?", sessionID).Error; err != nil {
		return "", err
	}
	return tokenPair.AccessToken, nil
}

func (s *Store) GetRefreshToken(sessionID string) (string, error) {
	var tokenPair models.TokenPair
	if err := s.GormClient.First(&tokenPair, "session_id = ?", sessionID).Error; err != nil {
		return "", err
	}
	return tokenPair.RefreshToken, nil
}

func (s *Store) IsAccessTokenValid(sessionID string) (bool, error) {
	var tokenPair models.TokenPair
	if err := s.GormClient.First(&tokenPair, "session_id = ?", sessionID).Error; err != nil {
		return false, err
	}
	return tokenPair.AccessToken != "" && tokenPair.AccessExpiresAt.After(time.Now()), nil
}

func (s *Store) IsRefreshTokenValid(sessionID string) (bool, error) {
	var tokenPair models.TokenPair
	if err := s.GormClient.First(&tokenPair, "session_id = ?", sessionID).Error; err != nil {
		return false, err
	}
	return tokenPair.RefreshToken != "" && tokenPair.RefreshExpiresAt.After(time.Now()), nil
}

func (s *Store) DeleteTokenPair(sessionID string) error {
	return s.GormClient.Delete(&models.TokenPair{}, "session_id = ?", sessionID).Error
}

func (s *Store) RefreshTokenPair(sessionID string, newAccessToken, newRefreshToken string, accessExpiresAt, refreshExpiresAt time.Time) error {
	return s.GormClient.Model(&models.TokenPair{}).Where("session_id = ?", sessionID).Updates(models.TokenPair{
		AccessToken:      newAccessToken,
		RefreshToken:     newRefreshToken,
		AccessExpiresAt:  accessExpiresAt,
		RefreshExpiresAt: refreshExpiresAt,
	}).Error
}

func (s *Store) RevokeAllUserTokens(userID string) error {
	return s.GormClient.Delete(&models.TokenPair{}, "user_id = ?", userID).Error
}

func (s *Store) GetUserActiveSessions(userID string) ([]string, error) {
	var tokenPairs []models.TokenPair
	if err := s.GormClient.Where("user_id = ?", userID).Find(&tokenPairs).Error; err != nil {
		return nil, err
	}
	var sessionIDs []string
	for _, tokenPair := range tokenPairs {
		sessionIDs = append(sessionIDs, tokenPair.SessionID)
	}
	return sessionIDs, nil
}

func (s *Store) UpdateSessionActivity(sessionID string, lastActivity time.Time) error {
	return s.GormClient.Model(&models.TokenPair{}).Where("session_id = ?", sessionID).Update("last_activity", lastActivity).Error
}

func (s *Store) StoreUserSession(session core.UserSession) error {
	// Convert core.UserSession to models.UserSession
	userSession := models.UserSession{
		SessionID:    session.SessionID,
		UserID:       session.UserID,
		IPAddress:    session.IPAddress,
		UserAgent:    session.UserAgent,
		LastActivity: session.LastActivity,
		// IsActive:     session.IsActive,
		// LoginTime:    session.LoginTime,
		// LogoutTime:   session.LogoutTime,
	}
	return s.GormClient.Create(&userSession).Error
}

func (s *Store) GetUserSession(sessionID string) (*core.UserSession, error) {
	var userSession models.UserSession
	if err := s.GormClient.First(&userSession, "session_id = ?", sessionID).Error; err != nil {
		return nil, err
	}

	// Convert models.UserSession to core.UserSession
	return &core.UserSession{
		SessionID:    userSession.SessionID,
		UserID:       userSession.UserID,
		IPAddress:    userSession.IPAddress,
		UserAgent:    userSession.UserAgent,
		LastActivity: userSession.LastActivity,
		// IsActive:     userSession.IsActive,
		// LoginTime:    userSession.LoginTime,
		// LogoutTime:   userSession.LogoutTime,
	}, nil
}

func (s *Store) DeleteUserSession(sessionID string) error {
	// Update logout time before deleting
	now := time.Now()
	if err := s.GormClient.Model(&models.UserSession{}).Where("session_id = ?", sessionID).Update("logout_time", now).Error; err != nil {
		// If session doesn't exist, continue with token deletion
		logger.Log.Error("Failed to update session logout time", "error", err)
	}

	// Delete from both token pairs and user sessions
	if err := s.GormClient.Delete(&models.TokenPair{}, "session_id = ?", sessionID).Error; err != nil {
		logger.Log.Error("Fail to delete token pair", "error", err)
		return err
	}
	err := s.GormClient.Delete(&models.UserSession{}, "session_id = ?", sessionID).Error
	if err != nil {
		logger.Log.Error("Fail to delete user session", "error", err)
		return err
	}
	return nil
}

func (s *Store) CleanupExpiredTokens() error {
	return s.GormClient.Where("access_expires_at < ?", time.Now()).Delete(&models.TokenPair{}).Error
}

func (s *Store) GetTokenPair(sessionID string) (*models.TokenPair, error) {
	var tokenPair models.TokenPair
	if err := s.GormClient.First(&tokenPair, "session_id = ?", sessionID).Error; err != nil {
		return nil, err
	}
	return &tokenPair, nil
}
