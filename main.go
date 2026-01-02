package main

import (
	"fmt"
	//"strings"
	"github.com/gocolly/colly"
)

// c = new collector object. Manages network communicaiton and responsible for handling golang callbacks.
func main() {
	c := colly.NewCollector(
		colly.AllowedDomains("divineoffice.org"),
	)

	//var prayerItems []string

	// Get to print out links of all .prayers-grid-items
	// Current issue: OnHTML not even being entered!
	//.page-content-wrapper > #frontpage_content > .wrapper > .primary > .community > .sidebar-intro > .intro-title-lite"
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
		//prayerItems = append(prayerItems, link)
		// If link start with browse or includes either signup or login return from callback
		/*
			if !strings.HasPrefix(link, "/browse") || strings.Contains(link, "=signup") > -1 || strings.Contains(link, "=login") > -1 {
				return
			}*/

		//e.Request.Visit(link)
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.Visit("https://divineoffice.org/")
}
