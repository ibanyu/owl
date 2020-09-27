package util

func ErrsJoin(str string, err []error) string {
	if len(err) < 1 {
		return ""
	}
	result := ""
	for i, v := range err {
		if v == nil {
			continue
		}

		if i == 0 {
			result += v.Error()
			continue
		}
		result += v.Error() + str
	}
	return result
}
