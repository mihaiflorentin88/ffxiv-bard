package contract

type UserInterface interface {
	HasTokenExpired() bool
	Persist() error
	GetName() string
}
