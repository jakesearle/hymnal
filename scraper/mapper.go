package main

func GetHymnNumbers(names []string, entries []*HymnEntry) (hymnNumbers []int) {
	officialNames := make([]string, 0, len(entries))
	for _, entry := range entries {
		officialNames = append(officialNames, CaselessPunctionless(entry.Name))
	}

	for _, name := range names {
		cleanName := CaselessPunctionless(name)
		no := containsNumber(entries, cleanName)
		if no != -1 {
			hymnNumbers = append(hymnNumbers, no)
		}
	}
	return
}

func MergeLocalRatings(entries []*HymnEntry, ratings map[int]int) []*HymnEntry {
	for _, entry := range entries {
		if rating, ok := ratings[entry.Number]; ok {
			entry.LocalRating = rating
		}
	}
	return entries
}

func MergeGlobalRatings(entries []*HymnEntry, ratings map[int]int) []*HymnEntry {
	for _, entry := range entries {
		if rating, ok := ratings[entry.Number]; ok {
			entry.GlobalRating = rating
		}
	}
	return entries
}

func containsNumber(entries []*HymnEntry, s string) int {
	for _, val := range entries {
		if CaselessPunctionless(s) == CaselessPunctionless(val.Name) {
			return val.Number
		}
	}
	return -1
}
