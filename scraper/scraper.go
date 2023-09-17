package main

import (
	"container/list"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

	"golang.org/x/net/html"
)

const cacheDir = "cache"
const baseUrl = "https://www.churchofjesuschrist.org"

type HymnEntry struct {
	Index          int
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
}

func main() {
	hymns := getHymnInfo()
	saveHymns(hymns)
}

func getHymnInfo() *list.List {
	noteMap := importNotes()
	requestUrl := baseUrl + "/study/manual/hymns?lang=eng"
	doc := getSoup(requestUrl)
	hymnList := list.New()
	currIndex := 1
	for _, group := range QueryAll(doc, "nav.manifest > ul.doc-map > li") {
		groupStr := GetText(group)
		fmt.Println(groupStr)
		for _, hymnTitle := range QueryAll(group, "a") {
			path := AttrOr(hymnTitle, "href", "")
			name := GetText(hymnTitle)
			url := baseUrl + path
			note := ""
			if val, ok := noteMap[currIndex]; ok {
				note = val
			}

			hymn := &HymnEntry{
				Index:          currIndex,
				Group:          groupStr,
				Language:       "en",
				Name:           name,
				Url:            url,
				Descriptor:     getDescriptor(url),
				NumVerses:      getNumVerses(url),
				NumExtraVerses: getExtraVerses(url),
				TextAuthor:     getTextAuthor(url),
				MusicAuthor:    getMusicAuthor(url),
				Notes:          note,
			}
			fmt.Println(hymn)
			hymnList.PushBack(hymn)
			currIndex++
			// break
		}
		// break
	}
	return hymnList
}

func saveHymns(hymns *list.List) {
	// Create or open the file for writing
	filename := "hymns.ndjson"
	file, err := os.Create(filename)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer file.Close()

	for hymnIter := hymns.Front(); hymnIter != nil; hymnIter = hymnIter.Next() {
		hymn := hymnIter.Value
		hymnJson, err := json.Marshal(hymn)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		file.Write(hymnJson)
		file.WriteString("\n")
	}
	fmt.Printf("Hymns saved to %s\n", filename)
}

func getDescriptor(url string) string {
	doc := getSoup(url)
	descriptorNode := Query(doc, "figure p.label")
	if descriptorNode == nil {
		return ""
	}
	return GetText(descriptorNode)
}

func getNumVerses(url string) int {
	doc := getSoup(url)
	return len(QueryAll(doc, "div.stanza"))
}

func getExtraVerses(url string) int {
	doc := getSoup(url)
	return len(QueryAll(doc, ".verses-below-the-music div.stanza"))
}

func getTextAuthor(url string) string {
	doc := getSoup(url)
	citations := QueryAll(doc, "div.citation-info p")
	return citations[0].LastChild.Data
}

func getMusicAuthor(url string) string {
	doc := getSoup(url)
	citations := QueryAll(doc, "div.citation-info p")
	return citations[len(citations)-1].LastChild.Data
}

func getSoup(url string) *html.Node {
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
