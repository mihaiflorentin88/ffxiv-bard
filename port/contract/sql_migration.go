package contract

type SqlMigrationDriverInterface interface {
	Execute(commandType string)
}
