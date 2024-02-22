package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"

	"golang.org/x/net/html"
)

const cacheDir = "cache"
const baseUrl = "https://www.churchofjesuschrist.org"
const authorPath = "/study/manual/using-the-hymnbook/authors-and-composers?lang=eng"
const hymnListingPath = "/study/manual/hymns?lang=eng"

type HymnEntry struct {
	Number         int
	Language       string
	Name           string
	Group          string
	Descriptor     string
	NumVerses      int
	NumExtraVerses int
	TextAuthor     string
	MusicAuthor    string
	Notes          string
	Url            string
	LocalRating    int
	GlobalRating   int
}

type Author struct {
	Name    string
	Credits []int
}

type AllRatings struct {
	LocalRatings  map[int]int `json:"LocalRatings"`
	GlobalRatings map[int]int `json:"GlobalRatings"`
}

func main() {
	ratings := LoadOrCacheRatings()
	printData(ratings)
}

func printData(ratings AllRatings) {
	combined := make(map[int]int, 0)
	for k, v := range ratings.GlobalRatings {
		combined[k] = v
	}
	for k, v := range ratings.LocalRatings {
		combined[k] += v
	}
	printValues(combined)
}

func getAuthors() []*Author {
	requestUrl := baseUrl + authorPath
	doc := GetSoup(requestUrl)
	authorNodes := QueryAll(doc, ".body ul li")
	authors := make([]*Author, 0)
	for _, authorNode := range authorNodes {
		pTags := QueryAll(authorNode, "p")
		if len(pTags) < 1 {
			fmt.Println("There's no author name here... Hm.")
			continue
		}
		if len(pTags) < 2 {
			fmt.Println("This author has no citations")
		}
		authorName := TrimLeadingAsterisks(GetText(pTags[0]))
		credits := Map2(pTags[1:], GetInt)
		author := &Author{
			Name:    authorName,
			Credits: credits,
		}
		authors = append(authors, author)
	}
	return authors
}

func getHymnInfo() []*HymnEntry {
	noteMap := importNotes()
	localRatings := importLocalRatings()
	requestUrl := baseUrl + hymnListingPath
	doc := GetSoup(requestUrl)
	hymnList := make([]*HymnEntry, 0, 0)
	currIndex := 1
	for _, group := range QueryAll(doc, "nav.manifest > ul.doc-map > li") {
		groupStr := GetText(group)
		for _, hymnTitle := range QueryAll(group, "a") {
			path := AttrOr(hymnTitle, "href", "")
			name := GetText(hymnTitle)
			url := baseUrl + path
			note := ""
			if val, ok := noteMap[currIndex]; ok {
				note = val
			}
			hymn := &HymnEntry{
				Number:         currIndex,
				Group:          groupStr,
				Language:       "en",
				Name:           name,
				Url:            url,
				Descriptor:     getDescriptor(url),
				NumVerses:      getNumVerses(url),
				NumExtraVerses: getExtraVerses(url),
				Notes:          note,
				LocalRating:    localRatings[currIndex],
			}
			hymnList = append(hymnList, hymn)
			currIndex++
		}
	}
	return hymnList
}

func getDescriptor(url string) string {
	doc := GetSoup(url)
	descriptorNode := Query(doc, "figure p.label")
	if descriptorNode == nil {
		return ""
	}
	return GetText(descriptorNode)
}

func getNumVerses(url string) int {
	doc := GetSoup(url)
	return len(QueryAll(doc, "div.stanza"))
}

func getExtraVerses(url string) int {
	doc := GetSoup(url)
	return len(QueryAll(doc, ".verses-below-the-music div.stanza"))
}

func getTextAuthor(url string) string {
	doc := GetSoup(url)
	citations := QueryAll(doc, "div.citation-info p")
	return citations[0].LastChild.Data
}

func getMusicAuthor(url string) string {
	doc := GetSoup(url)
	citations := QueryAll(doc, "div.citation-info p")
	return citations[len(citations)-1].LastChild.Data
}

func GetSoup(url string) *html.Node {
	htmlString := LoadOrCacheHtml(url)
	soup, err := html.Parse(strings.NewReader(htmlString))
	if err != nil {
		log.Fatal(err)
	}
	return soup
}

func importNotes() map[int]string {
	// Open the JSON file
	file, err := os.Open("notes.json")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return nil
	}
	defer file.Close()

	// Read the content of the file
	data, err := io.ReadAll(file)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return nil
	}

	// Create a variable of the struct type to hold the decoded data
	var notesMap map[string]interface{}

	// Unmarshal the JSON data into the Go data structure
	err = json.Unmarshal(data, &notesMap)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return nil
	}

	interpretedMap := make(map[int]string)
	for key, value := range notesMap {
		keyInt, err := strconv.Atoi(key)
		if err != nil {
			fmt.Println("Error parsing string to int:", err)
			continue
		}
		interpretedMap[int(keyInt)] = fmt.Sprint(value)
	}
	return interpretedMap
}

func importLocalRatings() map[int]int {
	playMap := GetWardHymnHistory()
	return calcHymnRatingMap(playMap)
}

func importGlobalRatings(hymnInfo []*HymnEntry) map[int]int {
	confNames := GetConferenceHymns()
	numbers := GetHymnNumbers(confNames, hymnInfo)
	numberCounts := fillEmptyHymns(GetCounterMap(numbers))
	return numberCounts
}

func printValues(popularity map[int]int) {
	// Convert the map to a slice of key-value pairs
	var pairs []Pair
	for k, v := range popularity {
		pairs = append(pairs, Pair{k, v})
	}

	// Sort the slice by value in descending order
	sort.Slice(pairs, func(i, j int) bool {
		return pairs[i].Value > pairs[j].Value
	})

	// Print the sorted map
	for _, pair := range pairs {
		fmt.Printf("%d: %d\n", pair.Key, pair.Value)
	}
}

type Pair struct {
	Key   int
	Value int
}
