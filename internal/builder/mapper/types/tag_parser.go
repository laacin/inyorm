package types

import (
	"strings"
	"unicode"

	"github.com/laacin/inyorm/internal/core"
)

type TagResult struct {
	Skip bool
	Name string
}

func ParseTag(fieldName, tag string) TagResult {
	result := TagResult{}

	for seq := range strings.SplitSeq(tag, core.TagKeySep) {
		seq = strings.TrimSpace(seq)

		if strings.ToLower(seq) == core.TagKeySkip {
			result.Skip = true
			continue
		}

		if strings.ToLower(seq) == core.TagKeyCol {
			result.Name = fieldName
			continue
		}

		keyVal := strings.Split(seq, core.TagValSep)
		if len(keyVal) < 2 {
			continue
		}

		key, val := keyVal[0], keyVal[1]
		if strings.ToLower(key) == core.TagKeyCol {
			result.Name = val
			continue
		}
	}

	if result.Name == "" {
		result.Name = toSnake(fieldName)
	}
	return result
}

// Helpers
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
