package main

import (
	"fmt"
	"os"
	"time"

	"github.com/Pimmr/rig"
	"github.com/Pimmr/rig/validators"
)

func palindrome(s string) error {
	for i, j := 0, len(s)-1; i < j; {
		if s[i] != s[j] {
			return fmt.Errorf("string should be a palindrome")
		}

		i++
		j--
	}

	return nil
}

func main() {
	var (
		flagA = 12
		flagB = 4.2
		flagC = "madam"
		flagD = 1 * time.Hour
	)

	err := rig.Parse(
		rig.Int(&flagA, "flag-a", "FLAG_A", "flag A", validators.IntRange(0, 54)),
		rig.Float64(&flagB, "flag-b", "FLAG_B", "flag B", validators.Float64Range(0.4, 12.5)),
		rig.String(
			&flagC, "flag-c", "FLAG_C", "flag C",
			validators.StringExcludeChars("bB"), validators.StringLengthMin(5), palindrome),
		rig.Duration(&flagD, "flag-d", "FLAG_D", "flag D", validators.DurationRounded(10*time.Minute)),
	)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(2)
	}

	fmt.Printf("flagA: %d\n", flagA)
	fmt.Printf("flagB: %f\n", flagB)
	fmt.Printf("flagC: %q\n", flagC)
	fmt.Printf("flagD: %s\n", flagD)
}
