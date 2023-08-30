package di

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/pkg/errors"
)

var (
	ErrNotSupport = fmt.Errorf("dp bind should be used on struct")
)

type bind struct {
	target   interface{}
	registry *Ioc
}

func NewBinder(target interface{}, registry *Ioc) *bind {
	return &bind{
		target:   target,
		registry: registry,
	}
}

func (b *bind) Bind() error {
	if b.target == nil {
		return fmt.Errorf("can not bind nil")
	}

	return b.bind(reflect.ValueOf(b.target), b.registry)
}

func (b *bind) bind(target reflect.Value, registry *Ioc) error {
	v := target

	for v.Kind() == reflect.Ptr || v.Kind() == reflect.Interface {
		if v.IsNil() {
			return fmt.Errorf("DipRegistry bind target nil")
		}
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return ErrNotSupport
	}

	if err := b.bindMember(target, v, registry); err != nil {
		return err
	}

	return nil
}

func (b *bind) bindMember(inputV reflect.Value, v reflect.Value, registry *Ioc) error {
	for i := 0; i != v.NumField(); i++ {
		fieldType := v.Type().Field(i)
		fieldValue := v.Field(i)
		tag := &Tag{}
		tag.ParseTag(fieldType, fieldType.Tag.Get("di"))

		if tag.IsSkip() {
			continue
		}

		if fieldType.Anonymous {
			if err := b.bind(fieldValue, b.registry); err != nil {
				if errors.Cause(err) == ErrNotSupport {
					continue
				} else {
					return errors.Wrapf(err, "embed field bind err, in type:%s, field:%s", v.Type(), fieldType.Name)
				}
			}
			continue
		}

		if tag.GetName() != "*" {
			value, err := registry.FetchByName(tag.GetName())
			if err != nil {
				return err
			}

			if value == nil {
				if tag.AllowEmpty() {
					continue
				} else {
					return fmt.Errorf("field value not found set in type:%s, field:%s", v.Type(), fieldType.Name)
				}
			}

			if err := b.SetField(inputV, fieldType, fieldValue, value); err != nil {
				return err
			}
		} else {
			value, err := registry.FetchByType(fieldType.Type)
			if err != nil {
				return err
			}

			if value == nil {
				if tag.AllowEmpty() {
					continue
				} else {
					return fmt.Errorf("field value not found set in type:%s, field:%s", v.Type(), fieldType.Name)
				}
			}

			if err := b.SetField(inputV, fieldType, fieldValue, value); err != nil {
				return err
			}
		}
	}

	return nil
}

func (b *bind) SetField(v reflect.Value, fieldType reflect.StructField, fieldValue reflect.Value, input interface{}) error {
	if fieldValue.CanSet() {
		if reflect.TypeOf(input).AssignableTo(fieldType.Type) {
			fieldValue.Set(reflect.ValueOf(input))
			return nil
		} else {
			return fmt.Errorf("field value not assignable in type:%s, field:%s, input:%s", v.Type(), fieldType.Name, reflect.TypeOf(input))
		}
	} else {
		setMethodName := fmt.Sprintf("Set%s", strings.Title(fieldType.Name))
		setMethod := v.MethodByName(setMethodName)
		if !setMethod.IsValid() {
			return fmt.Errorf("field value not assignable by method %s in type:%s, field:%s, input:%s", setMethodName, v.Type(), fieldType.Name, reflect.TypeOf(input))
		}

		setMethod.Call([]reflect.Value{reflect.ValueOf(input)})
		return nil
	}
}
