package stmt

import "strings"

const PrefixCol = "__COL__"

func IsColumn(value string) bool {
	return strings.HasPrefix(value, PrefixCol)
}

func ReadColumn(value string) string {
	return value[len(PrefixCol):]
}

func FilterColumn(value string) string {
	if IsColumn(value) {
		return ReadColumn(value)
	}

	return "'" + value + "'"
}

func SetColumn(value string) string {
	return PrefixCol + value
}
