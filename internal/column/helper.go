package column

func tbl(dflt string, provided []string) string {
	if len(provided) > 0 {
		return provided[0]
	}
	return dflt
}
