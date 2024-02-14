package contract

type MenuItem interface {
	GetName() string
	GetLink() string
	GetIcon() string
	GetDescription() string
}
