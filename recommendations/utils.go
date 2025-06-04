// recommendations/utils.go
package recommendations

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func mergeSlices(a, b []string) []string {
	result := append([]string{}, a...)
	for _, s := range b {
		if !contains(result, s) {
			result = append(result, s)
		}
	}
	return result
}
