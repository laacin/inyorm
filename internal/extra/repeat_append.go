package extra

func RepeatAppend[T any](slc []T, value T, times int) []T {
	if ln := len(slc); cap(slc) < ln+times {
		nslc := make([]T, ln, ln+times)
		copy(nslc, slc)

		for range times {
			nslc = append(nslc, value)
		}

		return nslc
	}

	for range times {
		slc = append(slc, value)
	}

	return slc
}
