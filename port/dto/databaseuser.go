package dto

import "time"

// DatabaseUser represents a user record from the database.
type DatabaseUser struct {
	ID                int        `db:"id"`
	Username          string     `db:"username"`
	Email             string     `db:"email"`
	Name              *string    `db:"name"` // Use pointer for nullable fields
	AuthProvider      *string    `db:"auth_provider"`
	ProviderID        *string    `db:"provider_id"`
	Avatar            *string    `db:"avatar"`
	Discriminator     *string    `db:"discriminator"`
	AuthorizationCode *string    `db:"authorization_code"`
	AccessToken       *string    `db:"access_token"`
	RefreshToken      *string    `db:"refresh_token"`
	TokenType         *string    `db:"token_type"`
	ExpiresAt         *time.Time `db:"expires_at"` // Use pointer for nullable fields
	TokenReleasedAt   *time.Time `db:"token_released_at"`
	CreatedAt         time.Time  `db:"created_at"`
	UpdatedAt         time.Time  `db:"updated_at"`
}

func NewDatabaseUser(
	id int,
	username string,
	email string,
	name,
	authProvider,
	providerID,
	avatar,
	discriminator,
	authorizationCode,
	accessToken,
	refreshToken,
	tokenType *string,
	expiresAt,
	tokenReleasedAt *time.Time,
	createdAt,
	updatedAt time.Time,
) *DatabaseUser {
	return &DatabaseUser{
		ID:                id,
		Username:          username,
		Email:             email,
		Name:              name,
		AuthProvider:      authProvider,
		ProviderID:        providerID,
		Avatar:            avatar,
		Discriminator:     discriminator,
		AuthorizationCode: authorizationCode,
		AccessToken:       accessToken,
		RefreshToken:      refreshToken,
		TokenType:         tokenType,
		ExpiresAt:         expiresAt,
		TokenReleasedAt:   tokenReleasedAt,
		CreatedAt:         createdAt,
		UpdatedAt:         updatedAt,
	}
}