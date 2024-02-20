package str

import (
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// Under score to camel
func UnderscoreToCamel(s string) string {
	c := cases.Title(language.English, cases.NoLower)
	parts := strings.Split(s, "_")
	for i, part := range parts {
		parts[i] = c.String(part)
	}
	return strings.Join(parts, "")
}
