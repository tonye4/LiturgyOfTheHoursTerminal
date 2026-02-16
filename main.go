package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/gocolly/colly"
	"golang.org/x/net/html"
)

// TODO:
type ApiResponse map[string]struct {
	Prayers []struct {
		PostTitle   string `json:"post_title"`
		PostContent string `json:"post_content"`
	} `json:"prayers"`
}

// Function that uses /net/html package to recurse through the html string and build to a new string
// only the values that are text contents and not of type element, giving us pure pre-formatted text.

// date format -> [yearmonthday]
// TODO: Use string builder to spit out URL.
// TODO: Get current day and day in 7 days.
// API Query -> "https://divineoffice.org/wp-json/do/v1/prayers/?date_start=[yearmonthday]&date_end=[yearmonthday]"
func getURL() string {
	var b strings.Builder

	// Current day and day seven days from now.
	currentDate := time.Now().Format(time.DateOnly)
	futureDate := time.Now().AddDate(0, 0, 7).Format(time.DateOnly)

	currentDateFormatted := strings.ReplaceAll(currentDate, "-", "")
	futureDateFormatted := strings.ReplaceAll(futureDate, "-", "")

	b.WriteString("https://divineoffice.org/wp-json/do/v1/prayers/?date_start=")
	b.WriteString(currentDateFormatted)
	b.WriteString("&date_end=")
	b.WriteString(futureDateFormatted)

	url := b.String()

	return url
}

func formatString(str string) string {
	// Can just use net/html
	doc, err := html.Parse(strings.NewReader(str))
	if err != nil {
		return ""
	}

	var b strings.Builder
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.TextNode {
			b.WriteString(n.Data)
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}

	f(doc)
	return b.String()
}

func main() {

	url := getURL()

	c := colly.NewCollector(
		colly.AllowedDomains("divineoffice.org"),
		//colly.Async(true),
	)

	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Accept", "application/json")
		r.Headers.Set("X-Requested-With", "XMLHttpRequest")
		r.Headers.Set("Referer", "https://divineoffice.org/")
	})

	c.OnResponse(func(r *colly.Response) {
		var apiResp ApiResponse

		fmt.Println("RAW RESPONSE:")
		fmt.Println(string(r.Body))

		err := json.Unmarshal(r.Body, &apiResp)
		if err != nil {
			log.Printf("JSON parsing error: %v", err)
			return
		}

		for date, day := range apiResp {
			fmt.Println(date)
			fmt.Println(day)
			for _, p := range day.Prayers {
				p.PostContent = formatString(p.PostContent)
				fmt.Printf("title: %s, content: %s\n", p.PostTitle, p.PostContent)
			}
		}
	})

	// Visit the API endpoint directly
	c.Visit(url)
	fmt.Println("this is the url: ", url)
	// TODO: Get current date each run of the script and format it in a similar way.
	// the api query is dynamic in that it changes each day, so need it to be updated.

	// Can implement caching, each week sunday, data is pulled into a .json file and then is just used accordingly to each day.
	// Will reduce the load on the server via batching.
	//https://divineoffice.org/wp-json/do/v1/prayers/?date_start={{20260129}}
}
