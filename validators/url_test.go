package validators

import (
	"fmt"
	"net/url"
	"testing"
)

func urlMustParse(s string) *url.URL {
	u, err := url.Parse(s)
	if err != nil {
		panic(fmt.Errorf("urlMustParse parsing %q: %w", s, err))
	}

	return u
}

func TestURLScheme(t *testing.T) {
	for _, test := range []struct {
		scheme      string
		input       *url.URL
		expectError bool
	}{
		{
			scheme:      "",
			input:       urlMustParse("https://example.com/foo"),
			expectError: false,
		},
		{
			scheme:      "",
			input:       urlMustParse("//example.com/foo"),
			expectError: false,
		},
		{
			scheme:      "https",
			input:       urlMustParse("https://example.com/foo"),
			expectError: false,
		},
		{
			scheme:      "https",
			input:       urlMustParse("http://example.com/foo"),
			expectError: true,
		},
		{
			scheme:      "https",
			input:       urlMustParse("//example.com/foo"),
			expectError: true,
		},
	} {
		err := URLScheme(test.scheme)(test.input)
		if test.expectError && err == nil {
			t.Errorf("URLScheme(%q)(%q): expected error, got nil", test.scheme, test.input)
		}
		if !test.expectError && err != nil {
			t.Errorf("URLScheme(%q)(%q): unexpected error: %s", test.scheme, test.input, err)
		}
	}
}

func TestURLExcludeScheme(t *testing.T) {
	for _, test := range []struct {
		scheme      string
		input       *url.URL
		expectError bool
	}{
		{
			scheme:      "",
			input:       urlMustParse("https://example.com/foo"),
			expectError: false,
		},
		{
			scheme:      "",
			input:       urlMustParse("//example.com/foo"),
			expectError: false,
		},
		{
			scheme:      "http",
			input:       urlMustParse("https://example.com/foo"),
			expectError: false,
		},
		{
			scheme:      "http",
			input:       urlMustParse("http://example.com/foo"),
			expectError: true,
		},
		{
			scheme:      "",
			input:       urlMustParse("http://example.com/foo"),
			expectError: false,
		},
		{
			scheme:      "http",
			input:       urlMustParse("//example.com/foo"),
			expectError: false,
		},
	} {
		err := URLExcludeScheme(test.scheme)(test.input)
		if test.expectError && err == nil {
			t.Errorf("URLExcludeScheme(%q)(%q): expected error, got nil", test.scheme, test.input)
		}
		if !test.expectError && err != nil {
			t.Errorf("URLExcludeScheme(%q)(%q): unexpected error: %s", test.scheme, test.input, err)
		}
	}
}
