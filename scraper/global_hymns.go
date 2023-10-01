package main

import "fmt"

const rootUrl = "https://www.churchofjesuschrist.org"
const conferencePath = "/study/music/conference-music?lang=eng"

func GetConferenceHymns() {
	fullUrl := rootUrl + conferencePath
	soup := GetSoup(fullUrl)
	conferences := QueryAll(soup, "#main a")
	fmt.Println(len(conferences))
}
