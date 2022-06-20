package simple

import (
	"fmt"
	di "github.com/douyacun/go-ioc"
	"testing"
)

type MessagePrinter interface {
	Print()
}

type UserProvider interface {
	GetUserName() string
}

type UserProviderImpl struct {
	UserName string `di:"userName"` // will match message
}

func (impl *UserProviderImpl) GetUserName() string {
	return impl.UserName
}

type MessagePrinterImpl struct {
	message string `di:"message"` // will match message
}

func (impl *MessagePrinterImpl) Print() {
	fmt.Println(impl.message)
}

func (impl *MessagePrinterImpl) SetMessage(m string) {
	impl.message = m
}

type GreeterService struct {
	UserProvider UserProvider   `di:"*"` // will match type
	Printer      MessagePrinter `di:"*"` // will match type
}

func (s *GreeterService) Print() {
	fmt.Printf("%s:", s.UserProvider.GetUserName())
	s.Printer.Print()
}

func Test(t *testing.T) {
	di.Register("message", "hello world")
	di.Register("userName", "douyacun")
	di.Register("printer", &MessagePrinterImpl{})
	di.Register("userProvider", &UserProviderImpl{})

	service := &GreeterService{}

	di.MustBind(service)
	service.Print()
}
