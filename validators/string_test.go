package validators

import "testing"

func TestStringNotEmpty(t *testing.T) {
	for _, test := range []struct {
		input       string
		expectError bool
	}{
		{
			input:       "",
			expectError: true,
		},
		{
			input:       " ",
			expectError: true,
		},
		{
			input:       "\t",
			expectError: true,
		},
		{
			input:       "foo",
			expectError: false,
		},
	} {
		err := StringNotEmpty()(test.input)
		if test.expectError && err == nil {
			t.Errorf("StringNotEmpty()(%q): expected error, got nil", test.input)
		}
		if !test.expectError && err != nil {
			t.Errorf("StringNotEmpty()(%q): unexpected error: %s", test.input, err)
		}
	}
}

func TestStringLengthRange(t *testing.T) {
	for _, test := range []struct {
		min, max    int
		input       string
		expectError bool
	}{
		{
			min:         2,
			max:         4,
			input:       "",
			expectError: true,
		},
		{
			min:         2,
			max:         4,
			input:       "f",
			expectError: true,
		},
		{
			min:         2,
			max:         4,
			input:       "fo",
			expectError: false,
		},
		{
			min:         2,
			max:         4,
			input:       "foo",
			expectError: false,
		},
		{
			min:         2,
			max:         4,
			input:       "foob",
			expectError: false,
		},
		{
			min:         2,
			max:         4,
			input:       "foobar",
			expectError: true,
		},
	} {
		err := StringLengthRange(test.min, test.max)(test.input)
		if test.expectError && err == nil {
			t.Errorf("StringLengthRange(%d, %d)(%q): expected error, got nil", test.min, test.max, test.input)
		}
		if !test.expectError && err != nil {
			t.Errorf("StringLengthRange(%d, %d)(%q): unexpected error: %s", test.min, test.max, test.input, err)
		}
	}
}

func TestStringLengthMin(t *testing.T) {
	for _, test := range []struct {
		min         int
		input       string
		expectError bool
	}{
		{
			min:         2,
			input:       "",
			expectError: true,
		},
		{
			min:         2,
			input:       "f",
			expectError: true,
		},
		{
			min:         2,
			input:       "fo",
			expectError: false,
		},
		{
			min:         2,
			input:       "foo",
			expectError: false,
		},
	} {
		err := StringLengthMin(test.min)(test.input)
		if test.expectError && err == nil {
			t.Errorf("StringLengthMin(%d)(%q): expected error, got nil", test.min, test.input)
		}
		if !test.expectError && err != nil {
			t.Errorf("StringLengthMin(%d)(%q): unexpected error: %s", test.min, test.input, err)
		}
	}
}

func TestStringLengthMax(t *testing.T) {
	for _, test := range []struct {
		max         int
		input       string
		expectError bool
	}{
		{
			max:         4,
			input:       "foo",
			expectError: false,
		},
		{
			max:         4,
			input:       "foob",
			expectError: false,
		},
		{
			max:         4,
			input:       "foobar",
			expectError: true,
		},
	} {
		err := StringLengthMax(test.max)(test.input)
		if test.expectError && err == nil {
			t.Errorf("StringLengthMax(%d)(%q): expected error, got nil", test.max, test.input)
		}
		if !test.expectError && err != nil {
			t.Errorf("StringLengthMax(%d)(%q): unexpected error: %s", test.max, test.input, err)
		}
	}
}

func TestStringExcludeChars(t *testing.T) {
	for _, test := range []struct {
		exclude     string
		input       string
		expectError bool
	}{
		{
			exclude:     "aeiou",
			input:       "foo",
			expectError: true,
		},
		{
			exclude:     "aeiou",
			input:       "bcdfg",
			expectError: false,
		},
		{
			exclude:     "",
			input:       "bcdfg",
			expectError: false,
		},
		{
			exclude:     "aeiou",
			input:       "",
			expectError: false,
		},
	} {
		err := StringExcludeChars(test.exclude)(test.input)
		if test.expectError && err == nil {
			t.Errorf("StringExcludeChars(%q)(%q): expected error, got nil", test.exclude, test.input)
		}
		if !test.expectError && err != nil {
			t.Errorf("StringExcludeChars(%q)(%q): unexpected error: %s", test.exclude, test.input, err)
		}
	}
}

func TestStringExcludePrefix(t *testing.T) {
	for _, test := range []struct {
		prefix      string
		input       string
		expectError bool
	}{
		{
			prefix:      "foo",
			input:       "foobar",
			expectError: true,
		},
		{
			prefix:      "foo",
			input:       "bar",
			expectError: false,
		},
		{
			prefix:      "foo",
			input:       "",
			expectError: false,
		},
		{
			prefix:      "",
			input:       "",
			expectError: true,
		},
		{
			prefix:      "",
			input:       "foo",
			expectError: true,
		},
	} {
		err := StringExcludePrefix(test.prefix)(test.input)
		if test.expectError && err == nil {
			t.Errorf("StringExcludePrefix(%q)(%q): expected error, got nil", test.prefix, test.input)
		}
		if !test.expectError && err != nil {
			t.Errorf("StringExcludePrefix(%q)(%q): unexpected error: %s", test.prefix, test.input, err)
		}
	}
}

func TestStringExcludeSuffix(t *testing.T) {
	for _, test := range []struct {
		suffix      string
		input       string
		expectError bool
	}{
		{
			suffix:      "bar",
			input:       "foobar",
			expectError: true,
		},
		{
			suffix:      "bar",
			input:       "foo",
			expectError: false,
		},
		{
			suffix:      "bar",
			input:       "",
			expectError: false,
		},
		{
			suffix:      "",
			input:       "",
			expectError: true,
		},
		{
			suffix:      "",
			input:       "bar",
			expectError: true,
		},
	} {
		err := StringExcludeSuffix(test.suffix)(test.input)
		if test.expectError && err == nil {
			t.Errorf("StringExcludeSuffix(%q)(%q): expected error, got nil", test.suffix, test.input)
		}
		if !test.expectError && err != nil {
			t.Errorf("StringExcludeSuffix(%q)(%q): unexpected error: %s", test.suffix, test.input, err)
		}
	}
}
