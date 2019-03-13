package rig

import (
	"strings"
	"testing"
)

func TestFlagMissingError(t *testing.T) {
	for _, test := range []struct {
		flag               Flag
		errorShouldContain []string
	}{
		{
			flag:               Flag{Name: "", Env: ""},
			errorShouldContain: []string{},
		},
		{
			flag:               Flag{Name: "foo", Env: ""},
			errorShouldContain: []string{"-foo"},
		},
		{
			flag:               Flag{Name: "", Env: "BAR"},
			errorShouldContain: []string{"BAR"},
		},
		{
			flag:               Flag{Name: "foo", Env: "BAR"},
			errorShouldContain: []string{"-foo", "BAR"},
		},
	} {
		err := test.flag.missingError()
		if err == nil {
			t.Errorf("Flag(%+v).missingError(): expected error, got nil", test.flag)
			continue
		}

		errStr := err.Error()
		for _, s := range test.errorShouldContain {
			if !strings.Contains(errStr, s) {
				t.Errorf("Flag(%+v).missingError() = %q: expected to find %q in error string.", test.flag, errStr, s)
			}
		}
	}
}
