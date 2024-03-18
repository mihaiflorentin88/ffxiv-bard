package main

import (
	"ffxvi-bard/cmd/cli"
	"ffxvi-bard/container"
)

func main() {
	serviceContainer := container.NewServiceContainer()
	container.Load = serviceContainer
	cli.Execute()
	connection := container.Load.DatabaseDriver()
	defer connection.Close()

	//signals := make(chan os.Signal, 1)
	//signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	//<-signals
}
