package main

import (
	"fmt"
	"os"

	"github.com/Pimmr/config"
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

func CountTheDots(v *uint, name, env, usage string) *config.Flag {
	return config.HintType(
		config.Var((*countTheDotsValue)(v), name, env, usage),
		"dotdotdot",
	)
}

func main() {
	var (
		flagA int
		flagB      = "foo"
		flagC uint = 7
	)

	err := config.Parse(
		config.Required(config.V(&flagA, "flag-a", "FLAG_A", "flag A")),
		config.V(&flagB, "flag-b", "FLAG_B", ""),
		CountTheDots(&flagC, "flag-c", "FLAG_C", "flag C"),
	)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(2)
	}

	fmt.Printf("flagA: %d\n", flagA)
	fmt.Printf("flagB: %q\n", flagB)
	fmt.Printf("flagC: %d\n", flagC)
}
