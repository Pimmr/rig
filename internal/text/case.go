package text

import "unicode"

func ToSnakeCase(s, sep string) string {
	ret := ""
	prev := '\000'

	rr := []rune(s)
	for i, r := range rr {
		if i != 0 && unicode.IsUpper(r) && unicode.IsLower(prev) {
			ret += sep
		} else if i != 0 && i != len(rr)-1 && unicode.IsUpper(r) && unicode.IsUpper(prev) && unicode.IsLower(rr[i+1]) {
			ret += sep
		}
		prev = r
		ret += string(unicode.ToLower(r))
	}

	return ret
}

func ToUpperSnakeCase(s, sep string) string {
	ret := ""
	prev := '\000'

	rr := []rune(s)
	for i, r := range rr {
		if i != 0 && unicode.IsUpper(r) && unicode.IsLower(prev) {
			ret += sep
		} else if i != 0 && i != len(rr)-1 && unicode.IsUpper(r) && unicode.IsUpper(prev) && unicode.IsLower(rr[i+1]) {
			ret += sep
		}
		prev = r
		ret += string(unicode.ToUpper(r))
	}

	return ret
}
