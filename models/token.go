package models

import (
	"time"
)

// TokenPair - Stores access and refresh tokens for user sessions
type TokenPair struct {
	SessionID        string    `gorm:"primaryKey;column:session_id" json:"sessionId"`
	AccessToken      string    `gorm:"column:access_token;not null" json:"accessToken,omitempty"`
	RefreshToken     string    `gorm:"column:refresh_token;not null" json:"refreshToken,omitempty"`
	AccessExpiresAt  time.Time `gorm:"column:access_expires_at;not null" json:"accessExpiresAt"`
	RefreshExpiresAt time.Time `gorm:"column:refresh_expires_at;not null" json:"refreshExpiresAt"`
	UserID           string    `gorm:"column:user_id;not null" json:"userId"`
	LastActivity     time.Time `gorm:"column:last_activity;autoUpdateTime" json:"lastActivity"`
	CreatedDate      time.Time `gorm:"column:created_date;autoCreateTime" json:"createdDate"`
	ModifiedDate     time.Time `gorm:"column:modified_date;autoUpdateTime" json:"modifiedDate"`

	// Relationships
	User *User `gorm:"foreignKey:UserID;references:UserID" json:"user,omitempty"`
}

func (TokenPair) TableName() string {
	return "auth_token_pairs"
}

// UserSession - Extended session information for tracking user activity
type UserSession struct {
	SessionID    string     `gorm:"primaryKey;column:session_id" json:"sessionId"`
	UserID       string     `gorm:"column:user_id;not null" json:"userId"`
	IPAddress    string     `gorm:"column:ip_address" json:"ipAddress"`
	UserAgent    string     `gorm:"column:user_agent" json:"userAgent"`
	LastActivity time.Time  `gorm:"column:last_activity;autoUpdateTime" json:"lastActivity"`
	IsActive     *bool      `gorm:"column:is_active;default:true" json:"isActive"`
	LoginTime    time.Time  `gorm:"column:login_time;autoCreateTime" json:"loginTime"`
	LogoutTime   *time.Time `gorm:"column:logout_time" json:"logoutTime,omitempty"`
	CreatedDate  time.Time  `gorm:"column:created_date;autoCreateTime" json:"createdDate"`
	ModifiedDate time.Time  `gorm:"column:modified_date;autoUpdateTime" json:"modifiedDate"`

	// Relationships
	User *User `gorm:"foreignKey:UserID;references:UserID" json:"user,omitempty"`
}

func (UserSession) TableName() string {
	return "auth_user_sessions"
}
