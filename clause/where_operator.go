package clause

import "strings"

func opEqual(ph *PlaceholderGen, negated bool) string {
	op := "= "
	if negated {
		op = "!= "
	}

	return op + ph.Next()
}

func opLike(ph *PlaceholderGen, negated bool) string {
	op := "LIKE "
	if negated {
		op = "NOT LIKE "
	}

	return op + ph.Next()
}

func opIn(ph *PlaceholderGen, negated bool, count int) string {
	placeholders := make([]string, count)
	for i := range placeholders {
		placeholders[i] = ph.Next()
	}

	inClause := "(" + strings.Join(placeholders, ", ") + ")"
	if negated {
		return "NOT IN " + inClause
	}
	return "IN " + inClause
}

func opBetween(ph *PlaceholderGen, negated bool) string {
	op := "BETWEEN "
	if negated {
		op = "NOT BETWEEN "
	}

	return op + ph.Next() + " AND " + ph.Next()
}

func opGreater(ph *PlaceholderGen, negated bool) string {
	op := "> "
	if negated {
		op = "<= "
	}

	return op + ph.Next()
}

func opLess(ph *PlaceholderGen, negated bool) string {
	op := "< "
	if negated {
		op = ">= "
	}
	return op + ph.Next()
}

func opIsNull(negated bool) string {
	if negated {
		return "IS NOT NULL"
	}
	return "IS NULL"
}
