# Golang 依赖注入

为什么我们需要依赖注入：

```
Most important, for me, is making it easy to follow the Single Responsibility Principle.

DI/IoC makes it simple for me to manage dependencies between objects. 
In turn, that makes it easier for me to break coherent functionality off into it's 
own contract (interface). 

As a result, my code has been far more modularized since I learned of DI/IoC.
Another result of this is that I can much more easily see my way through to a 
design that supports the Open-Closed Principle. This is one of the most confidence 
inspiring techniques (second only to automated testing). I doubt I could espouse the 
virtues of Open-Closed Principle enough.

DI/IoC is one of the few things in my programming career that has been a "game changer." 
There is a huge gap in quality between code I wrote before & after learning DI/IoC. 
Let me emphasize that some more. HUGE improvement in code quality.
```

# 安装

``

# 用法

构造一个问候服务需要知道

1. 向谁问候
2. 问候什么

```go
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
```

通常来说的初始化过程

```
message := "say hello"
printer := &MessagePrinterImpl{}
printer.SetMessage(message)

userName := "word"
userProvider := &UserProviderImpl{UserName: userName}

service := &MessagePrinterService{Printer: printer, UserProvider:userProvider}
service.Print()
```

使用依赖注入后

```
	di.Register("message", "hello world")
	di.Register("userName", "douyacun")
	di.Register("printer", &MessagePrinterImpl{})
	di.Register("userProvider", &UserProviderImpl{})

	service := &GreeterService{}

	di.MustBind(service)
	service.Print()
```


