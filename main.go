package main

import (
	"fmt"
	//"strings"
	"github.com/gocolly/colly"
)

type Reading struct {
	Date  string
	Title string
	Body  string
	URL   string
}

// c = new collector object. Manages network communicaiton and responsible for handling golang callbacks.
func main() {
	c := colly.NewCollector(
		colly.AllowedDomains("divineoffice.org"),
	)

	//readings := make([]Reading, 0, 10)

	cc := c.Clone()
	//var prayerItems []string

	// Get to print out links of all .prayers-grid-items
	// Current issue: OnHTML not even being entered!
	//.page-content-wrapper > #frontpage_content > .wrapper > .primary > .community > .sidebar-intro > .intro-title-lite"
	// div[class=sidebar-intro]
	c.OnHTML("div[class=sidebar-intro]", func(e *colly.HTMLElement) {
		fmt.Println("In onHTML method.")
		// beautiful now we know how to print out text!
		// All we need to know now is how to go through specific links!
		element := e.ChildText("h3") // e -> pointer html element node -> ChildText would be a subElement

		if element == "" {
			fmt.Println("Nothing found sorry!")
		} else {
			fmt.Println(element)
		}
	})

	// common pattern -> have a default colly object go through websites and then have a clone of it scrape specific details.
	// About today prayer and invatory and stuff has all the same class and -> has href which can be visited by c. detail scraper
	// can actually visit the href links.

	/* 	c.OnHTML(".home > .page-content-wrapper > #frontpage_content > .wrapper > .primary > .today-prayer", func(e *colly.HTMLElement) {
		fmt.Println("HELLURR", e.Name, e.Index)
		// theory: the html package which goQuery is built upon, has a fn called parse rejects elements more than 512 jjj
	}) */

	// usefult fact: the dates on the today-prayer hrefs
	c.OnHTML("a[href*=about-]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		fmt.Println("found link: ", link)

		cc.Visit(link)
		// find pattern for visiting and further scraping links.
	})

	cc.OnHTML(".entry", func(e *colly.HTMLElement) {
		// make look nicer after a lickle bit!
		date := e.DOM.Find(".entry > p:nth-of-type(1)").Text()
		title := e.DOM.Find(".entry > h2").Text() // we can use goQuery to actually get stuff wow very nice.
		// the body should be one big appended string -> need to account for the fact that the # of p's per page can vary!
		body := e.DOM.Find(".entry > p:nth-child(n+2)").Text()

		fmt.Println(date)
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
