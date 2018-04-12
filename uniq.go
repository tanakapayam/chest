package main

func uniq(elements []string) []string {
	seen := map[string]bool{}
	results := []string{}

	for _, e := range elements {
		if seen[e] {
		} else {
			seen[e] = true
			results = append(results, e)
		}
	}

	return results
}
