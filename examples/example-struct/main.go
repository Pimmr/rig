package main

import (
	"fmt"

	"github.com/Pimmr/rig"
	"github.com/Pimmr/rig/structToFlags"
)

type foo struct {
	FlagA string         `rig-flag:"flag-a" rig-env:"FLAG_A" rig-usage:"Flag A"`
	FlagB int            `rig-flag:"flag-b" rig-env:"FLAG_B" rig-usage:"Flag B" rig-required:"true"`
	FlagC []int          `rig-flag:"flag-c" rig-env:"FLAG_C" rig-usage:"Flag C" rig-typehint:"many ints"`
	FlagD []rig.URLValue `rig-flag:"flag-d" rig-usage:"Flag D"`
}

func main() {
	var f foo

	flags, err := structToFlags.StructToFlags(&f)
	if err != nil {
		panic(err)
	}

	err = rig.Parse(flags...)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", f)
}
