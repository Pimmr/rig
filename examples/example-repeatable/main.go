package main

import (
	"fmt"
	"os"
	"time"

	"github.com/Pimmr/rig"
)

type date time.Time

func (d date) String() string {
	return time.Time(d).Format("2006-01-02")
}

func (d *date) Set(s string) error {
	t, err := time.Parse("2006-01-02", s)
	*d = date(t)
	return err
}

func main() {
	var (
		flagA []int
		flagB []string
		flagC []rig.URLValue
		flagD []date
	)

	err := rig.Parse(
		rig.Required(rig.TypeHint(rig.Repeatable(
			&flagA, rig.IntGenerator(), "flag-a", "FLAG_A", "flag A",
		), "repeatable integer")),
		rig.TypeHint(rig.Repeatable(&flagB, rig.StringGenerator(), "flag-b", "FLAG_B", "flag B"), "repeatable string"),
		rig.TypeHint(rig.Repeatable(&flagC, rig.URLGenerator(), "flag-c", "FLAG_C", "flag C"), "repeatable URL"),
		rig.TypeHint(rig.Repeatable(&flagD, rig.MakeGenerator(new(date)), "flag-d", "FLAG_D", "flag D (i.e 2006-01-02)"), "repeatable date"),
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
