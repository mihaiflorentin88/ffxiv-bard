package database

import (
	"database/sql"
	"ffxvi-bard/port/contract"
	"ffxvi-bard/port/dto"
	"fmt"
	"strings"
	"time"
)

type userRepository struct {
	driver contract.DatabaseDriverInterface
}

func NewUserRepository(driver contract.DatabaseDriverInterface) contract.UserRepositoryInterface {
	return &userRepository{driver: driver}
}

func (u *userRepository) FindById(id int) (*dto.DatabaseUser, error) {
	query := `
		SELECT 
    		id,
    		username,
    		email,
    		name,
    		auth_provider,
    		provider_id,
    		avatar,
    		discriminator,
    		authorization_code,
		    access_token,
		    refresh_token,
		    token_type,
		    expires_at,
		    token_released_at,
		    created_at,
		    updated_at
		  FROM user WHERE id = ?`

	row, err := u.driver.FetchOne(query, id)
	if err != nil {
		return nil, fmt.Errorf("error fetching user by ID: %w", err)
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
			return nil, fmt.Errorf("no user found with ID %d", id)
		}
		return nil, fmt.Errorf("error scanning user: %w", err)
	}
	return &user, nil
}

func (u *userRepository) FindByUsername(username string) (*dto.DatabaseUser, error) {
	query := "SELECT * FROM user WHERE username = ?"
	row, err := u.driver.FetchOne(query, username)
	if err != nil {
		return nil, err
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
		return nil, err
	}
	return &user, nil
}

func (u *userRepository) FindByEmail(email string) (*dto.DatabaseUser, error) {
	query := "SELECT * FROM user WHERE email = ?"
	row, err := u.driver.FetchOne(query, email)
	if err != nil {
		return nil, err
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
		return nil, err
	}
	return &user, nil
}

func (u *userRepository) Create(user dto.DatabaseUser) error {
	query := `
		INSERT INTO user (
			username, email, name, auth_provider, provider_id, avatar,
			discriminator, authorization_code, access_token, refresh_token,
			token_type, expires_at, token_released_at, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	_, err := u.driver.Execute(query,
		user.Username, user.Email, user.Name, user.AuthProvider, user.ProviderID, user.Avatar,
		user.Discriminator, user.AuthorizationCode, user.AccessToken, user.RefreshToken,
		user.TokenType, user.ExpiresAt, user.TokenReleasedAt, time.Now(), time.Now(),
	)
	return err
}

func (u *userRepository) Update(user dto.DatabaseUser) error {
	query := `
		UPDATE user SET
			auth_provider=?, provider_id=?, avatar=?,
			discriminator=?, authorization_code=?, access_token=?, refresh_token=?,
			token_type=?, expires_at=?, token_released_at=?, updated_at=?
		WHERE id=?
	`
	_, err := u.driver.Execute(query,
		user.AuthProvider, user.ProviderID, user.Avatar,
		user.Discriminator, user.AuthorizationCode, user.AccessToken, user.RefreshToken,
		user.TokenType, user.ExpiresAt, user.TokenReleasedAt, time.Now(), user.ID,
	)
	return err
}

func (u *userRepository) FindByProperties(user dto.DatabaseUser) ([]dto.DatabaseUser, error) {
	var users []dto.DatabaseUser
	baseQuery := "SELECT * FROM user WHERE "
	var conditions []string
	var args []interface{}

	if user.Username != "" {
		conditions = append(conditions, "username = ?")
		args = append(args, user.Username)
	}
	if user.Email != "" {
		conditions = append(conditions, "email = ?")
		args = append(args, user.Email)
	}
	if user.Name != nil {
		conditions = append(conditions, "name = ?")
		args = append(args, *user.Name)
	}
	if user.AuthProvider != nil {
		conditions = append(conditions, "auth_provider = ?")
		args = append(args, *user.AuthProvider)
	}

	if user.ProviderID != nil {
		conditions = append(conditions, "provider_id = ?")
		args = append(args, *user.ProviderID)
	}

	if user.AuthorizationCode != nil {
		conditions = append(conditions, "authorization_code = ?")
		args = append(args, *user.AuthorizationCode)
	}

	if user.AccessToken != nil {
		conditions = append(conditions, "access_token = ?")
		args = append(args, *user.AccessToken)
	}

	if user.RefreshToken != nil {
		conditions = append(conditions, "refresh_token = ?")
		args = append(args, *user.RefreshToken)
	}

	if user.TokenType != nil {
		conditions = append(conditions, "token_type = ?")
		args = append(args, *user.TokenType)
	}

	if len(conditions) == 0 {
		return nil, fmt.Errorf("no properties provided to search for users")
	}

	query := baseQuery + strings.Join(conditions, " AND ")
	rows, err := u.driver.FetchMany(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var dbu dto.DatabaseUser
		err := rows.Scan(&dbu.ID, &dbu.Username, &dbu.Email, &dbu.Name, &dbu.AuthProvider, &dbu.ProviderID,
			&dbu.Avatar, &dbu.Discriminator, &dbu.AuthorizationCode, &dbu.AccessToken, &dbu.RefreshToken,
			&dbu.TokenType, &dbu.ExpiresAt, &dbu.TokenReleasedAt, &dbu.CreatedAt, &dbu.UpdatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, dbu)
	}

	return users, nil
}
