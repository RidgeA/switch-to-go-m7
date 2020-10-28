package crawler

func unique(strs []string) []string {
	exists := make(map[string]struct{})
	result := make([]string, 0)

	for _, v := range strs {
		_, has := exists[v]
		if !has {
			exists[v] = struct{}{}
			result = append(result, v)
		}
	}

	return result
}
