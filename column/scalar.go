package column

import (
	"strconv"
	"strings"

	"github.com/laacin/inyorm/internal/stmt"
)

func (c *Column) Concat(values ...string)  { wrapConcat(c, values) }
func (c *Column) Substring(start, end int) { wrapSub(c, start, end) }

func wrapConcat(c *Column, values []string) {
	var sb strings.Builder

	sb.WriteString("CONCAT(")
	for i, v := range values {
		if i > 0 {
			sb.WriteString(", ")
		}

		sb.WriteString(stmt.FilterColumn(v))
	}

	sb.WriteString(")")
	*c = Column(sb.String())
}

func wrapSub(c *Column, start, end int) {
	var sb strings.Builder
	sb.WriteString("SUBSTRING(")
	sb.WriteString(string(*c))
	sb.WriteString(", ")
	sb.WriteString(strconv.Itoa(start))
	sb.WriteString(", ")
	sb.WriteString(strconv.Itoa(end))
	sb.WriteString(")")
	*c = Column(sb.String())
}
