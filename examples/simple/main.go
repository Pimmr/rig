package main

import (
	"fmt"
	"os"

	"github.com/Pimmr/config"
)

func main() {
	var (
		flagA int
		flagB = "foo"
	)

	err := config.Parse(
		os.Args[1:],
		config.Required(config.Int(&flagA, "flag-a", "FLAG_A", "flag A")),
		config.String(&flagB, "flag-b", "FLAG_B", "flag B"),
	)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(2)
	}

	fmt.Printf("flagA: %d\n", flagA)
	fmt.Printf("flagB: %q\n", flagB)
}
