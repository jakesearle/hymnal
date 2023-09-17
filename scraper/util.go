package main

import (
	"regexp"
)

func TrimWhitespaceAndDates(input string) string {
	re := regexp.MustCompile(`^\s+|[,\d\s–]+$`)
	return re.ReplaceAllString(input, "")
}
