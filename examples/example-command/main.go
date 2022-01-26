package main

import (
	"fmt"

	"github.com/Pimmr/rig"
)

var hello = "world"

type FooConfig struct {
	A int
	B string
}

func Foo(c FooConfig) error {
	fmt.Printf("%q: %+v\n", hello, c)

	return nil
}

type BarConfig struct {
	A []string
	B float64
}

func Bar(c BarConfig) error {
	fmt.Printf("%q: %+v\n", hello, c)

	return nil
}

type BazConfig struct {
	A bool
	B bool
}

func Baz(c BazConfig) error {
	fmt.Printf("%q: %+v\n", hello, c)

	return nil
}

func main() {
	err := rig.Commands(
		rig.SubCommands("test", "test commands",
			rig.StructCommand("foo", Foo, "this is foo ..."),
			rig.StructCommand("bar", Bar, "this is bar ..."),
		),
		rig.StructCommand("baz", Baz, "this is baz ..."),
	).AdditionalFlags(
		rig.String(&hello, "hello", "HELLO", "hello ..."),
	).ParseAndCall()
	if err != nil {
		panic(err)
	}
}
