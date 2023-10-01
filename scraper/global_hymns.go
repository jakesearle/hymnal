package main

import (
	"golang.org/x/net/html"
)

const rootUrl = "https://www.churchofjesuschrist.org"
const conferencePath = "/study/music/conference-music?lang=eng"

func GetConferenceHymns() (hymnNames []string) {
	fullUrl := rootUrl + conferencePath
	soup := GetSoup(fullUrl)
	conferences := QueryAll(soup, "#main a")
	conferenceUrlPaths := Map2(conferences, getUrl)
	for _, path := range conferenceUrlPaths {
		hymnNames = append(hymnNames, getHymnNames(rootUrl+path)...)
	}
	return
}

func getHymnNames(url string) []string {
	soup := GetSoup(url)
	hymnNameNodes := QueryAll(soup, "h4 p.title")
	return Map2(hymnNameNodes, GetText)
}

func getUrl(node *html.Node) string {
	return AttrOr(node, "href", "")
}
