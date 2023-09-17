package main

import (
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
