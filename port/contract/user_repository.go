package contract

import "ffxvi-bard/port/dto"

type UserRepositoryInterface interface {
	FindById(id int64) (*dto.DatabaseUser, error)
	FindByUsername(username string) (*dto.DatabaseUser, error)
	FindByEmail(email string) (*dto.DatabaseUser, error)
	Create(user dto.DatabaseUser) error
	Update(user dto.DatabaseUser) error
	FindByProperties(user dto.DatabaseUser) ([]dto.DatabaseUser, error)
}
