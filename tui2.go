package main

/* MVP:
List of selectable prayers
Should look like doom emacs start menu
Prayers are displayed without anything fancy just scrolling
if required.
*/

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	tea "charm.land/bubbletea/v2"
)

type model struct {
	loadedPrayer string            // The user selected prayer.
	prayerList   map[string]string // Contains the prayers for a single day.
	prayerNames  []string          // Slice containing the names of the current day's prayers. Acts as a key to the prayerList in relation to the cursor position.
	selectedDay  string            // The day of the selected prayers.
	cursor       int               // Position of the cursor.
}

// NOTE: This should in theory be the same deal as the original tui file's prayerJson type but just more extensible.
/* type prayerJson map[string]struct {
	SpecPrayerList []specPrayer `json:"prayers"`
}

type specPrayer struct { // slice of structs because any date could have x amt of prayers with the same construct: PostTitle and PostContent.
	PostTitle   string `json:"post_title"`
	PostContent string `json:"post_content"`
} */

type prayerJson map[string]struct {
	Prayers []struct { // slice of structs because any date could have x amt of prayers with the same construct: PostTitle and PostContent.
		PostTitle   string `json:"post_title"`
		PostContent string `json:"post_content"`
	} `json:"prayers"`
}

// some globals to hold cached prayers.
var isCached bool = false
var prayers prayerJson

//
// Defining some Msg types to trigger the update and therefore the UI
//

type prayerMsg struct {
	day         string
	prayerNames []string
	prayerList  map[string]string
}

type daychangeMsg string

//
// main is where all the magic happens!
//

func main() {
	p := tea.NewProgram(model{})
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}

func (m model) Init() tea.Cmd {
	return loadJson()
}

//
// Cmd that grabs prayers for a specific day from the JSON prayer cache
// On initial launch, the current date is the default prayer set.
// Later on, checking for cached dates should be implemented.
//

// consider making date a type.
func loadJson(date string) tea.Cmd {
	return func() tea.Msg {
		if isCached == false {
			fileBytes, err1 := os.ReadFile("cached_prayers.json")
			if err1 != nil {
				log.Fatal(err1)
			}

			err2 := json.Unmarshal(fileBytes, &prayers)
			if err2 != nil {
				log.Fatal(err2)
			}

			isCached = true
		}

		var prayerTitles []string
		prayerTitles = make([]string, 6, 12)
		prayerList := make(map[string]string)

		// The loop populates... (finish documenting later)
		for i := 0; i < len(prayers[date].Prayers); i++ {
			prayerTitles[i] = prayers[date].Prayers[i].PostTitle
			prayerList[prayerTitles[i]] = prayers[date].Prayers[i].PostContent
		}

		return prayerMsg{
			date,
			prayerTitles,
			prayerList,
		}
	}
}
