package base

func Contains(arr []string, v string) bool {
	for _, k := range arr {
		if k == v {
			return true
		}
	}
	return false
}
