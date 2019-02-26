package main

import (
	"fmt"
	"os"
	"time"

	"github.com/Pimmr/config"
)

type Date time.Time

func (d Date) String() string {
	return time.Time(d).Format("2006-01-02")
}

func (d *Date) Set(s string) error {
	t, err := time.Parse("2006-01-02", s)
	*d = Date(t)
	return err
}

func main() {
	var (
		flagA []int
		flagB []string
		flagC []config.URLValue
		flagD []Date
	)

	err := config.Parse(
		config.Required(config.TypeHint(config.Repeatable(
			&flagA, config.IntGenerator(), "flag-a", "FLAG_A", "flag A",
		), "repeatable integer")),
		config.TypeHint(config.Repeatable(&flagB, config.StringGenerator(), "flag-b", "FLAG_B", "flag B"), "repeatable string"),
		config.TypeHint(config.Repeatable(&flagC, config.URLGenerator(), "flag-c", "FLAG_C", "flag C"), "repeatable URL"),
		config.TypeHint(config.Repeatable(&flagD, config.MakeGenerator(new(Date)), "flag-d", "FLAG_D", "flag D (i.e 2006-01-02)"), "repeatable date"),
	)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(2)
	}

	fmt.Printf("flagA: %v\n", flagA)
	fmt.Printf("flagB: %q\n", flagB)
	fmt.Printf("flagC: %q\n", flagC)
	fmt.Printf("flagD: %q\n", flagD)
}
