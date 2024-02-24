package str

import (
	"bytes"
	"strings"
	"unicode"
	"unicode/utf8"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// Under score to camel
func underscoreToCamel(s string) string {
	c := cases.Title(language.English, cases.NoLower)
	parts := strings.Split(s, "_")
	for i, part := range parts {
		parts[i] = c.String(part)
	}
	return strings.Join(parts, "")
}

// Under score to camel (abc_xyz to AbcXyz)
func UnderscoreToCapitalizeCamel(s string) string {
	if len(s) == 0 {
		return s
	}
	if len(s) == 1 {
		return string(unicode.ToUpper(rune(s[0])))
	}
	x := underscoreToCamel(s)
	return string(unicode.ToUpper(rune(x[0]))) + string(x[1:])
}

// Under score to lower camel (abc_xyz -> abcXyz)
// https://pkg.go.dev/strconv
func UnderscoreToLowerCamel(s string) string {
	if len(s) == 0 {
		return s
	}
	if len(s) == 1 {
		return string(unicode.ToLower(rune(s[0])))
	}
	x := underscoreToCamel(s)
	return string(unicode.ToLower(rune(x[0]))) + string(x[1:])
}

// First lower (Abcd -> abcd)
func FirstLower(s string) string {
	if len(s) == 0 {
		return s
	}
	if len(s) == 1 {
		return string(unicode.ToLower(rune(s[0])))
	}
	return string(unicode.ToLower(rune(s[0]))) + string(s[1:])
}

// First upper (abcd -> Abcd)
func FirstUpper(s string) string {
	if len(s) == 0 {
		return s
	}
	if len(s) == 1 {
		return string(unicode.ToLower(rune(s[0])))
	}
	return string(unicode.ToUpper(rune(s[0]))) + string(s[1:])
}

func IsASCII(s string) bool {
	return utf8.RuneCountInString(s) == len(s)
}

func RepeatString(count int, strs string) string {
	buf := &bytes.Buffer{}
	for i := 0; i < count; i++ {
		buf.WriteString(strs)
	}
	return buf.String()
}
