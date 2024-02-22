package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

func LoadOrCacheHtml(url string) string {
	filename := urlToFilename(url)
	if !fileExists(filename) {
		getAndSaveHtml(url, filename)
	}
	return getCacheFileContents(filename)
}

func LoadOrCacheRatings() AllRatings {
	filepath := getCacheFilepath("ratings.json")
	if !fileExists(filepath) {
		hymnInfo := getHymnInfo()
		globalRatings := importGlobalRatings(hymnInfo)
		localRatings := importLocalRatings()
		ratings := AllRatings{
			LocalRatings:  localRatings,
			GlobalRatings: globalRatings,
		}
		WriteObjectToFile(ratings, filepath)
	}

	var tmp AllRatings
	if err := ReadObjectFromFile(filepath, &tmp); err != nil {
		fmt.Println("Error reading object from file:", err)
	}
	return tmp
}

// WriteObjectToFile writes the given object to a JSON file.
func WriteObjectToFile(obj interface{}, filename string) error {
	// Marshal the object to JSON
	jsonData, err := json.Marshal(obj)
	if err != nil {
		return err
	}

	// Write JSON data to the file
	err = os.WriteFile(filename, jsonData, 0644)
	if err != nil {
		return err
	}

	fmt.Println("JSON data saved to", filename)
	return nil
}

// ReadObjectFromFile reads the JSON data from the file and unmarshals it into the given object.
func ReadObjectFromFile(filename string, obj interface{}) error {
	// Read JSON data from the file
	jsonData, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	// Unmarshal JSON data into the object
	err = json.Unmarshal(jsonData, &obj)
	if err != nil {
		return err
	}

	return nil
}

func getCacheFileContents(filename string) string {
	filepath := getCacheFilepath(filename)
	content, err := os.ReadFile(filepath)
	if err != nil {
		fmt.Println("Error:", err)
		return ""
	}
	return string(content)
}

func getAndSaveHtml(url, filename string) {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()

	os.Mkdir("./"+cacheDir, os.ModePerm)
	dest := getCacheFilepath(filename)
	file, err := os.Create(dest)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Printf("HTML saved to %s\n", filename)
}

func urlToFilename(url string) string {
	return strings.ReplaceAll(url, "/", "_") + ".html"
}

func fileExists(filename string) bool {
	filePath := getCacheFilepath(filename)
	_, err := os.Stat(filePath)

	if err == nil {
		return true
	}

	if os.IsNotExist(err) {
		return false
	} else {
		fmt.Println("Unexpected Error occurred:", err)
	}
	return false
}

func getCacheFilepath(filename string) string {
	return fmt.Sprintf("./%s/%s", cacheDir, filename)
}
