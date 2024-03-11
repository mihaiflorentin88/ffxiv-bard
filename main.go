package main

import (
	"ffxvi-bard/cmd/cli"
	"ffxvi-bard/container"
)

func main() {
	cli.Execute()
	serviceContainer := container.NewServiceContainer()
	connection := serviceContainer.GetDatabaseDriver()
	defer connection.Close()

	//signals := make(chan os.Signal, 1)
	//signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	//<-signals
}
