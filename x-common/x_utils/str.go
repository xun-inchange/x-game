package x_utils

import "strings"

func SplitSpace(str string) []string {
	return strings.Fields(strings.TrimSpace(str))
}