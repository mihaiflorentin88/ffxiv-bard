package container

import (
	"ffxvi-bard/config"
	"fmt"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"reflect"
	"sync"
)

var (
	Instances = make(map[string]interface{})
	mutex     sync.Mutex
	Container = ServiceContainer{}
)

type ServiceContainer struct{}

func (s *ServiceContainer) GetConfig() *config.Config {
	appConfig, err := config.NewConfig()
	if err != nil {
		panic(fmt.Sprintf("Cannot fetch the application config. Reason %s", err))
	}
	return appConfig
}

func GetConfig() *config.Config {
	appConfig, err := config.NewConfig()
	if err != nil {
		panic(fmt.Sprintf("Cannot fetch the application config. Reason %s", err))
	}
	return appConfig
}

//func GetInstance(name string, targetType interface{}) (interface{}, error) {
//	mutex.Lock()
//	defer mutex.Unlock()
//
//	if instance, ok := Instances[name]; ok {
//		target, ok := targetType.(*interface{})
//		if !ok {
//			return instance, nil
//		}
//		*target = instance
//		return target, nil
//	}
//	funcName := "Get" + cases.Title(language.English).String(name)
//	method := reflect.ValueOf(Container).MethodByName(funcName)
//	if !method.IsValid() {
//		return nil, fmt.Errorf("no method found for name %s", funcName)
//	}
//	results := method.Call(nil)
//	instance := results[0].Interface()
//	Instances[name] = instance
//	return instance, nil
//}

func GetInstance(name string) interface{} {
	mutex.Lock()
	defer mutex.Unlock()

	if instance, ok := Instances[name]; ok {
		return instance
	}

	funcName := "Get" + cases.Title(language.English).String(name)
	method := reflect.ValueOf(&Container).MethodByName(funcName)
	if !method.IsValid() {
		panic(fmt.Errorf("no method found for name %s", funcName))
	}

	results := method.Call(nil)
	if len(results) > 0 {
		instance := results[0].Interface()
		Instances[name] = instance
		return instance
	}
	panic(fmt.Errorf("method %s did not return an instance", funcName))
}
