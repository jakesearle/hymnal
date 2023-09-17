package main

import (
	"fmt"
	"log"
	"strings"

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

func example() {
	var pricingHtml string = `
<div class="card mb-4 box-shadow">
	<div class="card-header">
		<h4 class="my-0 font-weight-normal">Free</h4>
	</div>
	<div class="card-body">
		<h1 class="card-title pricing-card-title">$0/mo</h1>
		<ul class="list-unstyled mt-3 mb-4">
			<li>10 users included</li>
			<li>2 GB of storage</li>
			<li><a href="https://example.com">See more</a></li>
		</ul>
	</div>
</div>

<div class="card mb-4 box-shadow">
	<div class="card-header">
		<h4 class="my-0 font-weight-normal">Pro</h4>
	</div>
	<div class="card-body">
		<h1 class="card-title pricing-card-title">$15/mo</h1>
		<ul class="list-unstyled mt-3 mb-4">
			<li>20 users included</li>
			<li>10 GB of storage</li>
			<li><a href="https://example.com">See more</a></li>
		</ul>
	</div>
</div>

<div class="card mb-4 box-shadow">
	<div class="card-header">
		<h4 class="my-0 font-weight-normal">Enterprise</h4>
	</div>
	<div class="card-body">
		<h1 class="card-title pricing-card-title">$29/mo</h1>
		<ul class="list-unstyled mt-3 mb-4">
			<li>30 users included</li>
			<li>15 GB of storage</li>
			<li><a>See more</a></li>
		</ul>
	</div>
</div>
`
	doc, err := html.Parse(strings.NewReader(pricingHtml))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("List of pricing plans:\n\n")
	for i, p := range QueryAll(doc, "div.card.mb-4.box-shadow") {
		planName := Query(p, "h4").FirstChild.Data
		price := Query(p, ".pricing-card-title").FirstChild.Data
		usersIncluded := Query(p, "li:first-child").FirstChild.Data
		storage := Query(p, "li:nth-child(2)").FirstChild.Data
		detailsUrl := AttrOr(Query(p, "li:last-child a"), "href", "(No link available)")
		fmt.Printf(
			"Plan #%d\nName: %s\nPrice: %s\nUsers: %s\nStorage: %s\nDetails: %s\n\n",
			i+1,
			planName,
			price,
			usersIncluded,
			storage,
			detailsUrl,
		)
	}
}
