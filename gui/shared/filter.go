package shared

import "strings"

func FilterByContains(haystack []string, needle string) []string {
	result := make([]string, 0)
	for _, entry := range haystack {
		if strings.Contains(entry, needle) {
			result = append(result, entry)
		}
	}
	return result
}
