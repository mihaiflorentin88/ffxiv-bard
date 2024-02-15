package contract

import "ffxvi-bard/port/dto"

type UserInterface interface {
	HasTokenExpired() bool
	HydrateByEmail() error
	HydrateByUsername() error
	HydrateByID() error
	Persist() error
	GetName() string
	RenewToken() error
	ToDatabaseUserDTO() dto.DatabaseUser
}
