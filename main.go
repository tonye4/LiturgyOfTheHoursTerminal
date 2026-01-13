package main

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	"golang.org/x/net/html"
)

type Reading struct {
	Date  string
	Title string
	Body  string // consider making it a slice of paragraphs.
	URL   string
}

// This is the Text() method from goquery package, I just want to tweak it so that
// after each node it'll append a newline.
func TextL(s *goquery.Selection) string {
	var builder strings.Builder

	var f func(*html.Node)
	f = func(n *html.Node) {
		// n.Data prints out what the actual element is depending on what kind of Node it is. For element nodes -> possibilities are: p, span and br.
		// if next sibling of a text block is a <br> , "-" or "ant" don't append a new line.
		var noNl string

		if n.NextSibling != nil {
			noNl = n.NextSibling.Data
		}

		// Return bool if n.Data == "br" and then if it

		// Because span will typically always have a lone child text node, we gotta explicitly handle it.
		/* 		if n.Parent.Data == "span" {
			builder.WriteString(n.Data)
		} else */
		if n.Type == html.TextNode && noNl != "br" {
			builder.WriteString(n.Data + "\n")
		} else if n.Type == html.TextNode {
			builder.WriteString(n.Data)

		}

		if n.FirstChild != nil { // go through all of n's children recursively.
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				f(c) // append string here based off bool value.
			}
		}
	}
	for _, n := range s.Nodes {
		f(n)
	}
	return builder.String()
}

/* func TextLL(s *goquery.Selection) string {
	var builder strings.Builder

	// trim the front and end of the string.
	var f func(*html.Node) bool
	f = func(n *html.Node) bool {

		// Return bool if n.Data == "br" and then if it
		if n.Data == "br" {
			return true
		}
		// Reason about this later (only an hour...)
		// Is the issue w/ printint to do with using n within that for loop below?
		if n.FirstChild != nil { // go through all of n's children recursively.
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				brEncountered := f(c) // append string here based off bool value.

				if n.Type == html.TextNode {
					fmt.Println("I'm literally a textnode!")
				}

				if n.Type == html.TextNode && brEncountered == true {
					builder.WriteString(n.Data)
				} else if n.Type == html.TextNode && brEncountered == false {
					builder.WriteString(n.Data + "\n")
				}
			}
		}
		return false
	}
	for _, n := range s.Nodes {
		f(n)
	}

	return builder.String()
} */

// TODO: Make functions for each type of prayer -> Invitatory, office of readings etc...
// TODO: Slap text etc.. into json then have a file w/ bubbletea.

/* func invitatory(collect *)  */

// c = new collector object. Manages network communicaiton and responsible for handling golang callbacks.

func main() {
	c := colly.NewCollector(
		colly.AllowedDomains("divineoffice.org"),
	)

	cc := c.Clone()

	var body string

	// common pattern -> have a default colly object go through websites and then have a clone of it scrape specific details.
	// About today prayer and invatory and stuff has all the same class and -> has href which can be visited by c. detail scraper
	// can actually visit the href links.

	// useful fact: the dates on the today-prayer hrefs
	c.OnHTML("a[href*=ip-]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		fmt.Println("found link: ", link)

		cc.Visit(link)

	})

	// On each new <p> and <span> append a \n so that when it's all appended by .Text() it prints out in it's intended formatting.
	cc.OnHTML(".entry", func(e *colly.HTMLElement) {

		// selection object contains a slice of Node objects.
		/* 		date := e.DOM.Find(".entry > p:nth-of-type(1)").Text() */

		/*		This thang not working blud...
				title := e.DOM.Find("h1.intro-title")

				if title.Length() == 0 {
					fmt.Println("No elements found blud")
					return
				} */

		selectionObj := e.DOM.Find(".entry > p") // This isn't ideal bcos we now don't have a seperation of paragraphs for prettier formatting.
		body = TextL(selectionObj)

		//fmt.Println(date)

		fmt.Println(body)
		// print out child elements.
	})

	/* 	c.OnHTML("a[href*=-or]", func(e *colly.HTMLElement) {
	   		link := e.Attr("href")
	   		fmt.Println("found link: ", link)

	   		cc.Visit(link)
	   	})

	   	cc.OnHTML(".entry", func(e *colly.HTMLElement) {
	   		selectionObj := e.DOM.Find(".entry > p") // This isn't ideal bcos we now don't have a seperation of paragraphs for prettier formatting.
	   		body = TextL(selectionObj)

	   		//fmt.Println(date)

	   		fmt.Println(body)
	   		// print out child elements.
	   	}) */
	/* 	   	c.OnHTML("a[href*=-dp]", func(e *colly.HTMLElement) {
	   	   		link := e.Attr("href")
	   	   		fmt.Println("found link: ", link)
	   	   	})
	   	   	c.OnHTML("a[href*=-ep]", func(e *colly.HTMLElement) {
	   	   		link := e.Attr("href")
	   	   		fmt.Println("found link: ", link)
	   	   	})
	   	   	c.OnHTML("a[href*=-np-]", func(e *colly.HTMLElement) {
	   	   		link := e.Attr("href")
	   	   		fmt.Println("found link: ", link)
	   	   	}) */
	/* 	c.OnHTML(".prayers-grid", func(e *colly.HTMLElement) {
		fmt.Println("inside prayers-grid w")
	}) */

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	cc.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.Visit("https://divineoffice.org/")
}
