package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/Pimmr/rig"
	"github.com/Pimmr/rig/validators"
)

type countTheDotsValue uint

func (d countTheDotsValue) String() string {
	if d == 0 {
		return "<none>"
	}

	s := ""
	for i := 0; i < int(d); i++ {
		s += "."
	}

	return s
}

func (d *countTheDotsValue) Set(s string) error {
	for _, c := range s {
		if c != '.' {
			return fmt.Errorf("%q is not a dot", c)
		}
	}

	*d = countTheDotsValue(len(s))
	return nil
}

func countTheDots(v *uint, name, env, usage string) *rig.Flag {
	return rig.TypeHint(
		rig.Var(
			(*countTheDotsValue)(v), name, env, usage,
			rarToIntValidator(rangeValidator(1, 8)),
		),
		"dotdotdot",
	)
}

func rangeValidator(min, max int) validators.Int {
	return func(i int) error {
		if i < min {
			return fmt.Errorf("integer should be greater than %d", min)
		}
		if i > max {
			return fmt.Errorf("integer should be less than %d", max)
		}

		return nil
	}
}

func rarToIntValidator(validator validators.Int) validators.Var {
	return func(v flag.Value) error {
		i, ok := v.(*countTheDotsValue)
		if !ok {
			return fmt.Errorf("expected type *int")
		}

		return validator(int(*i))
	}
}

func main() {
	var (
		flagA int
		flagB      = "foo"
		flagC uint = 7
		flagD []string
		flagE bool
	)

	err := rig.Parse(
		rig.Required(rig.Int(&flagA, "flag-a", "FLAG_A", "flag A", rangeValidator(1, 667))),
		rig.String(&flagB, "flag-b", "FLAG_B", ""),
		countTheDots(&flagC, "flag-c", "FLAG_C", "flag C"),
		rig.Repeatable(&flagD, rig.StringGenerator(), "flag-d", "FLAG_D", "flag D"),
		rig.Bool(&flagE, "flag-e-extra-long", "FLAG_E_LONG", "flag E"),
	)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(2)
	}

	fmt.Printf("flagA: %d\n", flagA)
	fmt.Printf("flagB: %q\n", flagB)
	fmt.Printf("flagC: %d\n", flagC)
	fmt.Printf("flagD: %q\n", flagD)
	fmt.Printf("flagE: %v\n", flagE)
}
