package di

import (
	"encoding/json"
	"testing"
)

func TestStructBindName(t *testing.T) {
	type BarService interface{}

	type JojoService interface{}

	type FooBase struct {
		JojoService JojoService `di:"JojoService"`
	}

	type Foo struct {
		FooBase    `di:"embed"`
		BarService BarService `di:"BarService"`
	}

	target := &Foo{}
	target.FooBase.JojoService = 1

	r := NewRegistry()
	r.Register("BarService", struct{}{})
	r.Register("JojoService", struct{}{})

	b := NewBinder(target, r)
	if err := b.Bind(); err != nil {
		t.Fatal(err)
	}

	t.Log(target)
}

func TestStructBindType(t *testing.T) {
	type BarService interface{}

	type JojoService interface{}

	type FooBase struct {
		JojoService int `di:"*"`
	}

	type Foo struct {
		FooBase    `di:"embed"`
		BarService string `di:"*"`
	}

	target := &Foo{}
	target.FooBase.JojoService = 1

	r := NewRegistry()
	r.Register("whateverInt", 1)
	r.Register("whateverString", "test")

	b := NewBinder(target, r)
	if err := b.Bind(); err != nil {
		t.Fatal(err)
	}

	t.Log(target)
}

func TestStructBindFunc(t *testing.T) {
	type BarService interface{}

	type JojoService interface{}

	type FooBase struct {
		JojoService int `di:"*"`
	}

	type Foo struct {
		FooBase    `di:"embed"`
		BarService string `di:"*"`
	}

	target := &Foo{}
	target.FooBase.JojoService = 1

	r := NewRegistry()
	r.Register("whateverInt", func(string) int { return 1 })
	r.Register("whateverString", func() string { return "test" })

	b := NewBinder(target, r)
	if err := b.Bind(); err != nil {
		t.Fatal(err)
	}

	t.Log(target)
}

func TestBind(t *testing.T) {
	instance := newInstance()

	gInstance = instance

	type FooBase struct {
		JojoService int `di:"*"`
	}

	type Foo struct {
		FooBase    `di:"embed"`
		BarService string `di:"*"`
	}

	Register("string", "string")
	Register("int", 1)
	Register("foo", &Foo{})

	var jbm = &struct {
		Test    string `di:"*"`
		TestInt int    `di:"*"`
		Foo     *Foo   `di:"*"`
	}{}
	MustBind(&jbm)

	data, _ := json.Marshal(jbm)
	t.Log(string(data))
}
