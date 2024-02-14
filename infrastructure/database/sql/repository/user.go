package database

import (
	"database/sql"
	"ffxvi-bard/port/contract"
	"ffxvi-bard/port/dto"
	"fmt"
)

type userRepository struct {
	driver contract.DatabaseDriverInterface
}

func NewUserRepository(driver contract.DatabaseDriverInterface) contract.UserRepositoryInterface {
	return &userRepository{driver: driver}
}

func (u *userRepository) GetUserById(id int) (dto.DatabaseUser, error) {
	query := `SELECT id, username, email, name, auth_provider, provider_id, avatar, discriminator, authorization_code,
			  access_token, refresh_token, token_type, expires_at, token_released_at, created_at, updated_at
			  FROM user WHERE id = ?`

	row, err := u.driver.FetchOne(query, id)
	if err != nil {
		return dto.DatabaseUser{}, fmt.Errorf("error fetching user by ID: %w", err)
	}

	var user dto.DatabaseUser
	err = row.Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Name,
		&user.AuthProvider,
		&user.ProviderID,
		&user.Avatar,
		&user.Discriminator,
		&user.AuthorizationCode,
		&user.AccessToken,
		&user.RefreshToken,
		&user.TokenType,
		&user.ExpiresAt,
		&user.TokenReleasedAt,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return dto.DatabaseUser{}, fmt.Errorf("no user found with ID %d", id)
		}
		return dto.DatabaseUser{}, fmt.Errorf("error scanning user: %w", err)
	}
	return user, nil
}
