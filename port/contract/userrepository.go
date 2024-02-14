package contract

import "ffxvi-bard/port/dto"

type UserRepositoryInterface interface {
	GetUserById(id int) (dto.DatabaseUser, error)
}
