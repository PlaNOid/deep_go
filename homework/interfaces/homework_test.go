package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

type UserService struct {
	// not need to implement
	NotEmptyStruct bool
}
type MessageService struct {
	// not need to implement
	NotEmptyStruct bool
}

type Container struct {
	constructors map[string]func(*Container) (interface{}, error)
	instances    map[string]interface{}
}

func NewContainer() *Container {
	return &Container{
		constructors: make(map[string]func(*Container) (interface{}, error)),
		instances:    make(map[string]interface{}),
	}
}

func (c *Container) RegisterType(name string, constructor interface{}) {
	switch cons := constructor.(type) {
	case func() interface{}:
		c.constructors[name] = func(container *Container) (interface{}, error) {
			return cons(), nil
		}
	default:
		panic(fmt.Sprintf("Неподдерживаемая сигнатура для '%s'", name))
	}
}

func (c *Container) Resolve(name string) (interface{}, error) {
	if instance, ok := c.instances[name]; ok {
		return instance, nil
	}

	constructor, exists := c.constructors[name]
	if !exists {
		return nil, fmt.Errorf("тип '%s' не зарегистрирован", name)
	}

	instance, err := constructor(c)
	if err != nil {
		return nil, fmt.Errorf("ошибка создания '%s': %w", name, err)
	}

	c.instances[name] = instance
	return instance, nil
}

func TestDIContainer(t *testing.T) {
	container := NewContainer()
	container.RegisterType("UserService", func() interface{} {
		return &UserService{}
	})
	container.RegisterType("MessageService", func() interface{} {
		return &MessageService{}
	})

	userService1, err := container.Resolve("UserService")
	assert.NoError(t, err)
	userService2, err := container.Resolve("UserService")
	assert.NoError(t, err)

	u1 := userService1.(*UserService)
	u2 := userService2.(*UserService)
	assert.True(t, u1 == u2) // из-за singleton

	messageService, err := container.Resolve("MessageService")
	assert.NoError(t, err)
	assert.NotNil(t, messageService)

	paymentService, err := container.Resolve("PaymentService")
	assert.Error(t, err)
	assert.Nil(t, paymentService)
}
