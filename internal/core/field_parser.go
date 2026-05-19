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

type FieldSchema struct {
	Skip  bool
	Name  string
	Index []int
}

func NewFieldSchema(name, tag string, idx []int) FieldSchema {
	info := FieldSchema{Index: idx}

	for seq := range strings.SplitSeq(tag, ",") {
		seq = strings.TrimSpace(seq)

		if strings.ToLower(seq) == "skip" {
			info.Skip = true
			continue
		}

		if strings.ToLower(seq) == "col" {
			info.Name = name
			continue
		}

		keyVal := strings.Split(seq, ":")
		if len(keyVal) < 2 {
			continue
		}

		key, val := keyVal[0], keyVal[1]
		if strings.ToLower(key) == "col" {
			info.Name = val
			continue
		}
	}

	if info.Name == "" {
		info.Name = toSnake(name)
	}
	return info
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
