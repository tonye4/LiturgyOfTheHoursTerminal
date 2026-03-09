package main

/*
MVP:
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
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"
)

type modelz struct {
	choices []string // items on the to-do list
	cursor  int      // which to-do list item our cursor is pointing at selected map[int]struct{} // which to-do items are selected }
}

type prayerJsonz map[string]struct {
	Prayers []struct { // slice of structs because any date could have x amt of prayers with the same construct: PostTitle and PostContent.
		PostTitle   string `json:"post_title"`
		PostContent string `json:"post_content"`
	} `json:"prayers"`
}

func getDay() string {
	currentDate := time.Now().Format(time.DateOnly)
	currentDateFormatted := strings.ReplaceAll(currentDate, "-", "")

	return currentDateFormatted
}

var prayerss prayerJson // should we init this in the model?
var todayDate = getDay()

func initialModel() model {
	//todayDate := getDay()

	fileBytes, err1 := os.ReadFile("cached_prayers.json")
	if err1 != nil {
		log.Fatal(err1)
	}

	//var prayers prayerJson // make it a global var.

	err2 := json.Unmarshal(fileBytes, &prayers)
	if err2 != nil {
		log.Fatal(err2)
	}

	//fmt.Println(prayers[todayDate].Prayers[0].PostTitle)

	//
	// Clean up the code after you get it working, don't want pretty code right
	// now just want something that works.
	// Note: The slice needs to be dynamic.
	//

	var prayerTitles []string
	prayerTitles = make([]string, 10) // make dynamic.
	number := 0

	for number < len(prayers[todayDate].Prayers) {
		//fmt.Println(prayers[todayDate].Prayers[number].PostTitle)
		prayerTitles[number] = prayers[todayDate].Prayers[number].PostTitle
		//fmt.Println(prayerTitles[number])

		number++
	}

	return model{
		// This is a slice, try making a map and unpack the title of the prayers into there.
		choices: prayerTitles[:8],
		//choices: prayers[todayDate].Prayers,

		// A map which indicates which choices are selected. We're using
		// the  map like a mathematical set. The keys refer to the indexes
		// of the `choices` slice, above.
		selected: make(map[int]struct{}),
	}
}

// Init I/O will be reading our .JSON file.
func (m model) Initz() tea.Cmd {
	// Just return `nil`, which means "no I/O Sright now, please."
	return nil
}

func (m model) Updatez(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	// Is it a key press?
	case tea.KeyPressMsg:

		// Cool, what was the actual key pressed?
		switch msg.String() {

		// These keys should exit the program.
		case "ctrl+c", "q":
			return m, tea.Quit

		// The "up" and "k" keys move the cursor up
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		// The "down" and "j" keys move the cursor down
		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}

			// The "enter" key and the space bar toggle the selected state
			// for the item that the cursor is pointing at.
		case "enter", "space":
			_, ok := m.selected[m.cursor]
			if ok {
				delete(m.selected, m.cursor) // we'll keep this for now
			} else {
				m.selected[m.cursor] = struct{}{} // pretty much denotes that the selected choice has been selected.
			}
		}
	}

	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return m, nil
}

func (m model) Viewz() tea.View {
	// The header
	// This can be our variable we change in order to print out the actual prayers.

	//
	// 's' should be updated to reflect the prayer that was selected.
	//

	s := "Select Prayer\n\n"

	// Iterate over our choices
	for i, choice := range m.choices {

		// Is the cursor pointing at this choice?
		cursor := " " // no cursor
		if m.cursor == i {
			cursor = ">" // cursor!
		}

		// Is this choice selected?
		//checked := " " // not selected
		// does the m.selected index EXIST! (refer to update ie was it initialized with an empty struct or not nil sturct{}{})
		if _, ok := m.selected[i]; ok {
			s = prayers[todayDate].Prayers[0].PostContent
		}

		// Render the row
		s += fmt.Sprintf("%s %s\n", cursor, choice)
	}

	// The footer
	s += "\nPress q to quit.\n"

	// Send the UI for rendering
	return tea.NewView(s)
}

func main2() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
