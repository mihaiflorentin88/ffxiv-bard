package container

import "ffxvi-bard/config"

func GetConfig() *config.Config {
	appConfig, err := config.NewConfig()
	if err != nil {
		panic("Cannot load the application config.")
	}
	return appConfig
}
