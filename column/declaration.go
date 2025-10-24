package column

import "github.com/laacin/inyorm/internal/stmt"

type Column string

func (c Column) Name() string  { return wrapName(c) }
func (c *Column) As(as string) { wrapAs(c, as) }

// ----- Helpers
func wrapName(c Column) string {
	return stmt.SetColumn(string(c))
}

func wrapAs(c *Column, as string) {
	*c = *c + " AS " + Column(as)
}
