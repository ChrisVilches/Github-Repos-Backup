package util

import "strings"

func FormatDate(input string) string {
	input = strings.TrimSuffix(input, "Z")
	input = strings.ReplaceAll(input, "-", "")
	input = strings.ReplaceAll(input, ":", "")
	input = strings.Replace(input, "T", "-", 1)
	return input
}
