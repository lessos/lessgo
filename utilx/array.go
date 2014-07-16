package utilx

func ArrayEqual(a, b []string) bool {

	if len(a) != len(b) {
		return false
	}

	for _, va := range a {

		for _, vb := range b {

			if va != vb {
				return false
			}
		}
	}

	return true
}

func ArrayContain(s string, a []string) bool {

	for _, v := range a {

		if v == s {
			return true
		}
	}

	return false
}
