package core

import (
	"strings"
	"unicode"
)

const TAG = "inyorm"

// ---- Tag keys
const (
	SKIP = "skip"
	COL  = "c"
	TBL  = "t"
)

type FieldResult struct {
	Name  string
	Index []int
}

func ParseField(fieldName, tag string, idx []int) FieldResult {
	name := toSnake(fieldName)
	return FieldResult{Name: name, Index: idx}
}

func toSnake(v string) string {
	if v == "" {
		return ""
	}

	var b strings.Builder
	runes := []rune(v)

	for i, r := range runes {
		if unicode.IsUpper(r) {
			if i > 0 {
				prev := runes[i-1]

				if unicode.IsLower(prev) ||
					(unicode.IsUpper(prev) &&
						i+1 < len(runes) &&
						unicode.IsLower(runes[i+1])) {
					b.WriteByte('_')
				}
			}

			b.WriteRune(unicode.ToLower(r))
			continue
		}

		b.WriteRune(r)
	}

	return b.String()
}
