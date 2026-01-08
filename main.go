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
	Body  []string // consider making it a slice of paragraphs.
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

		if n.Type == html.TextNode && noNl != "br" {
			builder.WriteString(n.Data + "\n")
		} else if n.Type == html.TextNode {
			builder.WriteString(n.Data)
		}
		if n.FirstChild != nil {
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				f(c)
			}
		}
	}
	for _, n := range s.Nodes {
		f(n)
	}

	return builder.String()
}

// c = new collector object. Manages network communicaiton and responsible for handling golang callbacks.
func main() {
	c := colly.NewCollector(
		colly.AllowedDomains("divineoffice.org"),
	)

	//readings := make([]Reading, 0, 10)

	cc := c.Clone()

	/* 	c.OnHTML("div[class=sidebar-intro]", func(e *colly.HTMLElement) {
		fmt.Println("In onHTML method.")

		element := e.ChildText("h3")

		if element == "" {
			fmt.Println("Nothing found sorry!")
		} else {
			fmt.Println(element)
		}
	}) */

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
		title := e.DOM.Find(".entry > h2").Text() // we can use goQuery to actually get stuff wow very nice.

		selectionObj := e.DOM.Find(".entry > p") // This isn't ideal bcos we now don't have a seperation of paragraphs for prettier formatting.
		body := TextL(selectionObj)

		// placeholder variable until I set up the actual object!

		/* 		pNodes := e.DOM.Find(".entry >p:nth-of-type(n+2)")
		   		for _, v := range pNodes.Nodes {
		   			if v.Type == html.TextNode {
		   				fmt.Println(v.Data)
		   			}
		   		} */

		// This is how .Text() is implemented make it append /n to EVERY element and use it instead of Text itself.

		// this is the genral idea, but it's prolly not how to go about making a slice of paragraphs.
		/* 		e.DOM.Find(".entry").Each(func(row int, s *goquery.Selection) {
			paragraphs = append(paragraphs, s.Find("p").Text()) // hopefully dat works lol
		}) */

		//fmt.Println(date)
		fmt.Println(title)

		fmt.Println(body)
		// print out child elements.
	})

	/* 	c.OnHTML("a[href*=ord-]", func(e *colly.HTMLElement) {
	   		link := e.Attr("href")
	   		fmt.Println("found link: ", link)
	   	})
	   	c.OnHTML("a[href*=-dp]", func(e *colly.HTMLElement) {
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
