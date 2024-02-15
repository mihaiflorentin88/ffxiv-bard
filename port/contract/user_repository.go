package contract

import "ffxvi-bard/port/dto"

type UserRepositoryInterface interface {
	FindUserById(id int) (*dto.DatabaseUser, error)
	FindUserByUsername(username string) (*dto.DatabaseUser, error)
	FindUserByEmail(email string) (*dto.DatabaseUser, error)
	InsertUser(user dto.DatabaseUser) error
	UpdateUser(user dto.DatabaseUser) error
	FindUsersByProperties(user dto.DatabaseUser) ([]dto.DatabaseUser, error)
}
