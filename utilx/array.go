package utilx

func ArrayEqual(a, b []string) bool {

	if len(a) != len(b) {
		return false
	}

	for _, va := range a {

		eq := false

		for _, vb := range b {

			if va == vb {
				eq = true
				break
			}
		}

		if !eq {
			return false
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
