package shared

import "strings"

func SplitAndTrim(toSplit string) []string {
	theSlice := make([]string, 0)
	for _, tag := range strings.Split(toSplit, ",") {
		if len(tag) > 0 {
			theSlice = append(
				theSlice,
				strings.ReplaceAll(strings.TrimSpace(tag), " ", "-"),
			)
		}
	}
	return theSlice
}
