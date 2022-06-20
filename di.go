package di

import (
	"sync"
)

var (
	once      sync.Once
	gInstance *diInstance
)

type diInstance struct {
	registry *Ioc
}

func newInstance() *diInstance {
	return &diInstance{
		registry: NewRegistry(),
	}
}

func init() {
	once.Do(func() {
		gInstance = newInstance()
	})
}

func Register(name string, value interface{}, opts ...RegistryOption) {
	gInstance.registry.Register(name, value)
}

func Bind(target interface{}) error {
	return gInstance.registry.Bind(target)
}

func MustBind(target interface{}) {
	if err := Bind(target); err != nil {
		panic(err)
	}
}
