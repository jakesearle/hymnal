package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
)

// Returns a map of each time a hymn numbers was sung in my local ward from 2016-2023
// Maps from a hymn number to the number of times it was played
func GetWardHymnHistory() map[int]int {
	data := loadData()
	flattenedData := flatten(data)
	hymnNumberList := parseHymnNumber(flattenedData)
	hymnCounts := GetCounterMap(hymnNumberList)
	return fillEmptyHymns(hymnCounts)
}

func loadData() [][]string {
	// open file
	f, err := os.Open("./2016-2023-hymns.csv")
	if err != nil {
		log.Fatal(err)
	}

	// remember to close the file at the end of the program
	defer f.Close()

	// read csv values using csv.Reader
	csvReader := csv.NewReader(f)
	data, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}
	return data
}

func flatten(data [][]string) []string {
	ret := make([]string, 0)
	for _, row := range data {
		for j, cell := range row {
			if j > 0 {
				log.Fatal("There should only be one column of data")
				continue
			}
			ret = append(ret, cell)
		}
	}
	return ret
}

func parseHymnNumber(data []string) []int {
	numericalListings := make([]int, 0)
	for _, hymnListing := range data {
		hymnListing = DeleteTrailingNonNumerals(hymnListing)
		if num, err := strconv.Atoi(hymnListing); err == nil {
			numericalListings = append(numericalListings, num)
			continue
		}
		fmt.Printf("Not a number: %q\n", hymnListing)
	}
	return numericalListings
}

func fillEmptyHymns(numberToPlayCount map[int]int) map[int]int {
	for i := 0; i <= NumEnglishHymns; i++ {
		if _, exists := numberToPlayCount[i]; !exists {
			numberToPlayCount[i] = 0
		}
	}
	return numberToPlayCount
}
