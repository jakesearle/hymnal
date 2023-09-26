package main

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
)

func TrimWhitespaceAndDates(input string) string {
	re := regexp.MustCompile(`^\s+|[,\d\s–]+$`)
	return re.ReplaceAllString(input, "")
}

func TrimLeadingAsterisks(input string) string {
	re := regexp.MustCompile(`^\*+`)
	return re.ReplaceAllString(input, "")
}

func Map2[T, U any](data []T, f func(T) U) []U {
	res := make([]U, 0, len(data))
	for _, e := range data {
		res = append(res, f(e))
	}
	return res
}

func SaveHymns(hymns []*HymnEntry, filename string) {
	// Create or open the file for writing
	filename += ".ndjson"
	file, err := os.Create(filename)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer file.Close()

	for _, item := range hymns {
		hymnJson, err := json.Marshal(item)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		file.Write(hymnJson)
		file.WriteString("\n")
	}
}

func SaveAuthors(authors []*Author, filename string) {
	// Create or open the file for writing
	filename += ".ndjson"
	file, err := os.Create(filename)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer file.Close()

	for _, item := range authors {
		authorJson, err := json.Marshal(item)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		file.Write(authorJson)
		file.WriteString("\n")
	}
}
