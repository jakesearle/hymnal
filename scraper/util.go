package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
)

func GetCounterMap(input []int) map[int]int {
	m := make(map[int]int)
	for _, input := range input {
		m[input]++
	}
	return m
}

func DeleteTrailingNonNumerals(input string) string {
	re := regexp.MustCompile(`\D+$`)
	return re.ReplaceAllString(input, "")
}

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
	builder := strings.Builder{}
	for _, item := range hymns {
		hymnJson, err := json.Marshal(item)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		builder.Write(hymnJson)
		builder.WriteString("\n")
	}
	Write(builder.String(), filename)
}

func SaveAuthors(authors []*Author, filename string) {
	// Create or open the file for writing
	filename += ".ndjson"
	builder := strings.Builder{}
	for _, item := range authors {
		authorJson, err := json.Marshal(item)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		builder.Write(authorJson)
		builder.WriteString("\n")
	}
	Write(builder.String(), filename)
}

func Write(value, filename string) {
	if err := os.WriteFile(filename, []byte(value), 0666); err != nil {
		log.Fatal(err)
	}
}
