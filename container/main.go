package container

import (
	"ffxvi-bard/config"
	"fmt"
)

var Load *ServiceContainer

type ServiceContainer struct {
	config         *config.Config
	domain         *DomainContainer
	infrastructure *InfrastructureContainer
	http           *HttpContainer
}

func NewServiceContainer() *ServiceContainer {
	return &ServiceContainer{
		domain:         &DomainContainer{},
		infrastructure: &InfrastructureContainer{},
		http:           &HttpContainer{},
	}
}

func (s *ServiceContainer) Config() *config.Config {
	if s.config != nil {
		return s.config
	}
	appConfig, err := config.NewConfig()
	s.config = appConfig
	if err != nil {
		panic(fmt.Sprintf("Cannot fetch the application config. Reason %s", err))
	}
	return appConfig
}
