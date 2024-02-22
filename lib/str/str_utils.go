package str

import (
	"strings"
	"unicode"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// Under score to camel (abc_xyz to AbcXyz)
func UnderscoreToCamel(s string) string {
	c := cases.Title(language.English, cases.NoLower)
	parts := strings.Split(s, "_")
	for i, part := range parts {
		parts[i] = c.String(part)
	}
	return strings.Join(parts, "")
}

// Under score to lower camel (abc_xyz to abcXyz)
func UnderscoreToLowerCamel(s string) string {
	x := UnderscoreToCamel(s)
	return string(unicode.ToLower(rune(x[0]))) + string(x[1:])
}
