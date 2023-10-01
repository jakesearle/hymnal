package main

import (
	"math"
	"sort"
)

const NumEnglishHymns int = 341

type HymnFrequencyPair struct {
	hymnNumber int
	playCount  int
}

func calcHymnRatingMap(hymnNoToPlays map[int]int) map[int]int {
	sorted := sortByValues(hymnNoToPlays)
	return convert(sorted)
}

func convert(plays []HymnFrequencyPair) map[int]int {
	sextileSize := int(math.Ceil(float64(NumEnglishHymns) / 6.0))
	result := make(map[int]int, 0)
	for i := 0; i < 6; i++ {
		for j := i; j < i+sextileSize && j < NumEnglishHymns; j++ {
			result[plays[j].hymnNumber] = i
		}
	}
	return result
}

func sortByValues(inputMap map[int]int) []HymnFrequencyPair {
	var sortedSlice []HymnFrequencyPair
	// Convert the map to a slice of key-value pairs
	for k, v := range inputMap {
		sortedSlice = append(sortedSlice, HymnFrequencyPair{k, v})
	}
	// Sort the slice by values
	sort.Slice(sortedSlice, func(i, j int) bool {
		return sortedSlice[i].playCount < sortedSlice[j].playCount
	})
	return sortedSlice
}
