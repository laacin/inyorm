package standard

import "strconv"

func (dial *DialectStandard) String(v string) string {
	return "'" + v + "'"
}

func (dial *DialectStandard) Number(v int) string {
	return strconv.Itoa(v)
}

func (dial *DialectStandard) Float(v float64) string {
	return strconv.FormatFloat(float64(v), 'f', -1, 32)
}

func (dial *DialectStandard) Bool(v bool) string {
	if v {
		return "1"
	}
	return "0"
}

func (dial *DialectStandard) Null() string {
	return "NULL"
}

func (dial *DialectStandard) Placeholder(number int) string {
	return "?"
}

func (dial *DialectStandard) Wildcard() string {
	return "*"
}
