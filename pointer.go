package rig

import (
	"flag"
	"reflect"

	"github.com/pkg/errors"
)

type pointerFlag struct {
	Value PointerValue

	Var reflect.Value
}

func (p pointerFlag) String() string {
	if p.Value.IsNil() {
		return "<nil>"
	}

	return p.Value.String()
}

func (p *pointerFlag) Set(s string) error {
	if !p.Value.IsNil() {
		return p.Value.Set(s)
	}

	t := p.Var.Type().Elem().Elem()
	v := reflect.New(t)
	p.Var.Elem().Set(v)
	p.Value = noopInstanciator{
		Value: p.Value.New(v.Interface()),
	}

	return p.Set(s)
}

type noopInstanciator struct {
	flag.Value
}

func (noopInstanciator) New(interface{}) flag.Value {
	panic(errors.New("Not Implemented"))
}

func (noopInstanciator) IsNil() bool {
	return false
}

type PointerValue interface {
	flag.Value

	New(interface{}) flag.Value
	IsNil() bool
}

func Pointer(f *Flag, v interface{}) *Flag {
	iv, ok := f.Value.(PointerValue)
	if !ok {
		panic(errors.Errorf("%T does not implement the rig.PointerValue interface", f.Value))
	}

	return &Flag{
		Value: &pointerFlag{
			Value: iv,
			Var:   reflect.ValueOf(v),
		},
		Name:     f.Name,
		Env:      f.Env,
		Usage:    f.Usage,
		TypeHint: f.TypeHint,
		Required: f.Required,

		set:          f.set,
		defaultValue: f.defaultValue,
	}
}
