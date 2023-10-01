package main

import (
	"fmt"
	"strconv"

	"github.com/andybalholm/cascadia"
	"golang.org/x/net/html"
)

func Query(n *html.Node, query string) *html.Node {
	sel, err := cascadia.Parse(query)
	if err != nil {
		return &html.Node{}
	}
	return cascadia.Query(n, sel)
}

func QueryAll(n *html.Node, query string) []*html.Node {
	sel, err := cascadia.Parse(query)
	if err != nil {
		return []*html.Node{}
	}
	return cascadia.QueryAll(n, sel)
}

func AttrOr(n *html.Node, attrName, or string) string {
	for _, a := range n.Attr {
		if a.Key == attrName {
			return a.Val
		}
	}
	return or
}

func GetText(n *html.Node) string {
	// Drill down to the firstmost descendant and get the data
	if n.FirstChild == nil {
		return n.Data
	}
	return GetText(n.FirstChild)
}

func GetInt(n *html.Node) int {
	ret, err := strconv.Atoi(GetText(n))
	if err != nil {
		fmt.Println("Error parsing string to int:", err)
		panic(err)
	}
	return ret
}
