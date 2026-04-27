package prayers

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/gocolly/colly"
	"golang.org/x/net/html"
)

// TODO: Change up the naming convention of main.go and the packages.
/*
TODO: Change the name of this file to getPrayers.
The actual main file is going to be the TUI frontend and it's going to be
using the backend ie prompting caching through a function call.
We can get our frontend to check during initialization if more prayers are
required ie, the last cached day was passed or there are no cached prayers
in general, which would call these various functions.

Before we link these together ie make this file a library of prayer functions
I want to use the .json file seperately and just run the TUI to
see how it works and how it's gonna look and then we can connect them together.

TUI is going to be Charm Bracelet and Lip Gloss (for styles).
*/
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
	futureDate := time.Now().AddDate(0, 0, 5).Format(time.DateOnly)

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

func cachePrayers(prayers ApiResponse) {
	// Marshal map into JSON object
	jsonedPrayers, err1 := json.MarshalIndent(prayers, "", "\t")

	if err1 != nil {
		panic(err1)
	}

	//fmt.Println(string(jsonedPrayers))
	// Write object into file
	// Create second file that holds data
	file, err2 := os.Create("cached_prayers.json")

	// logic to create a cache_history file.
	// check if file exists.

	if err2 != nil {
		panic(err2)
	}

	_, err3 := file.WriteAt(jsonedPrayers, 0)

	if err3 != nil {
		panic(err3)
	}

	fmt.Println("Succesfuly wrote to file: cached_prayers.json")
	// return int representing success or fail.
}

func GetPrayers() {
	// Before caling endpoint, check if cached_prayers.json exists
	// and also check if it's last day (last 2 digits of the string)

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
		// do I nedd to run make to give the map some memory?
		var apiResp ApiResponse

		err := json.Unmarshal(r.Body, &apiResp)
		if err != nil {
			log.Printf("JSON parsing error: %v", err)
			return
		}

		cachePrayers(apiResp)
	})

	// Visit the API endpoint directly
	c.Visit(url)
	// the api query is dynamic in that it changes each day, so need it to be updated.

	// Can implement caching, each week sunday, data is pulled into a .json file and then is just used accordingly to each day.
	// Will reduce the load on the server via batching.
	//https://divineoffice.org/wp-json/do/v1/prayers/?date_start={{20260129}}
}
