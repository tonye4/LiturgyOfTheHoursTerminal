package main

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"

	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/tonye4/LiturgyOfTheHoursTerminal/prayers"
)

// TODO: add line folding. Entire paragraphs go off screen if there's no html element seperating each line.

// ─── View states ─────────────────────────────────────────────────────────────

type viewState int

const (
	menuState viewState = iota
	prayerState
)

// ─── Styles ──────────────────────────────────────────────────────────────────

var (
	appTitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#7D56F4")).
			Padding(0, 3)

	subtitleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#888888")).
			Italic(true)

	cursorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#7D56F4")).
			Bold(true)

	selectedItemStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FAFAFA")).
				Background(lipgloss.Color("#7D56F4"))

	normalItemStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#DDDDDD"))

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#555555"))

	vpTitleStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Right = "├"
		return lipgloss.NewStyle().BorderStyle(b).Padding(0, 1).Bold(true)
	}()

	vpFooterStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Left = "┤"
		return vpTitleStyle.BorderStyle(b)
	}()
)

// ─── Model ───────────────────────────────────────────────────────────────────

type model struct {
	prayerNames   []string
	prayerList    map[string]string
	prayerDate    string
	cursor        int
	state         viewState
	selectedTitle string
	viewport      viewport.Model
	ready         bool
	termWidth     int
	termHeight    int
}

// ─── Messages ────────────────────────────────────────────────────────────────

type prayerLoadedMsg struct {
	names      []string
	prayerList map[string]string
	date       string
}

type errMsg struct{ err error }

// ─── Commands ────────────────────────────────────────────────────────────────

func loadPrayersCmd() tea.Cmd {
	return func() tea.Msg {
		// Read in and unmarshall our cached_prayers.json into the prayerList map.
		// It should be populated but if it's not we need to fetch more prayers.
		fileBytes, err := os.ReadFile("cached_prayers.json")
		if err != nil {
			return errMsg{fmt.Errorf("could not read cached_prayers.json: %w", err)}
		}

		var prayersList prayers.ApiResponse
		if err := json.Unmarshal(fileBytes, &prayersList); err != nil {
			return errMsg{fmt.Errorf("could not parse prayers JSON: %w", err)}
		}

		// Check if our list of prayers contains our date.
		// If our today date doesn't exist, we scrape
		// for more prayers for today and future dates.
		// TODO: Create a test for this bit of logic.
		date := Today()
		day, ok := prayersList[date]
		if !ok {
			prayers.GetPrayers()
			day = prayersList[date]
		}

		var names []string
		list := make(map[string]string)
		for _, p := range day.Prayers {
			names = append(names, p.PostTitle)
			list[p.PostTitle] = prayers.FormatString(p.PostContent)
		}

		return prayerLoadedMsg{names, list, date}
	}
}

// ─── Init ────────────────────────────────────────────────────────────────────

func (m model) Init() tea.Cmd {
	// TODO: Get the program to have a loading page if the getting prayers takes a while.
	return loadPrayersCmd()
}

// ─── Update ──────────────────────────────────────────────────────────────────

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {

	// After the prayer loaded msg returns we set the returned values to our model struct fields for rendering.
	case prayerLoadedMsg:
		m.prayerNames = msg.names
		m.prayerList = msg.prayerList
		m.prayerDate = msg.date
		return m, nil

	case errMsg:
		fmt.Println("Error:", msg.err)
		return m, tea.Quit

	case tea.WindowSizeMsg:
		m.termWidth = msg.Width
		m.termHeight = msg.Height

		if m.state == prayerState && m.ready {
			m.viewport.SetWidth(msg.Width)
			m.viewport.SetHeight(msg.Height - lipgloss.Height(m.headerView()) - lipgloss.Height(m.footerView()))
		}

	case tea.KeyPressMsg:
		switch m.state {

		case menuState:
			switch msg.String() {
			case "ctrl+c", "q":
				return m, tea.Quit
			case "up", "k":
				if m.cursor > 0 {
					m.cursor--
				}
			case "down", "j":
				if m.cursor < len(m.prayerNames)-1 {
					m.cursor++
				}
			// When the user presses enter, we set the specific prayer
			// within the viewport.
			case "enter", " ":
				if len(m.prayerNames) == 0 {
					break
				}
				m.selectedTitle = m.prayerNames[m.cursor]
				m.state = prayerState

				headerH := lipgloss.Height(m.headerView())
				footerH := lipgloss.Height(m.footerView())
				m.viewport = viewport.New(
					viewport.WithWidth(m.termWidth),
					viewport.WithHeight(m.termHeight-headerH-footerH),
				)
				m.viewport.YPosition = headerH
				m.viewport.HighlightStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#ff0000"))
				m.viewport.SelectedHighlightStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#ff0000"))
				content := m.prayerList[m.selectedTitle]
				m.viewport.SetContent(content)
				m.viewport.SetHighlights(regexp.MustCompile(`\bChrist\b|\bJesus\b`).FindAllStringIndex(content, -1))
				m.ready = true
			}

		case prayerState:
			switch msg.String() {
			case "ctrl+c":
				return m, tea.Quit
			case "esc", "backspace", "q":
				m.state = menuState
				m.ready = false
				return m, nil
			}
		}
	}

	if m.state == prayerState && m.ready {
		m.viewport, cmd = m.viewport.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

// ─── View ────────────────────────────────────────────────────────────────────

func (m model) View() tea.View {
	// where all the actual rendering happens
	var v tea.View
	v.AltScreen = true
	v.MouseMode = tea.MouseModeCellMotion

	switch m.state {
	case menuState:
		v.SetContent(m.renderMenu())
	case prayerState:
		if !m.ready {
			v.SetContent("\n  Loading...")
		} else {
			v.SetContent(fmt.Sprintf("%s\n%s\n%s",
				m.headerView(), m.viewport.View(), m.footerView()))
		}
	}

	return v
}

// ─── Menu rendering ──────────────────────────────────────────────────────────

func (m model) renderMenu() string {
	if m.termWidth == 0 {
		return "\n  Loading..."
	}

	center := lipgloss.NewStyle().Width(m.termWidth).Align(lipgloss.Center)

	var lines []string
	lines = append(lines, appTitleStyle.Render("Divine Office"))
	lines = append(lines, subtitleStyle.Render(formatDate(m.prayerDate)))
	lines = append(lines, "")

	if len(m.prayerNames) == 0 {
		lines = append(lines, subtitleStyle.Render("Loading prayers..."))
	} else {
		// find the widest prayer name for consistent item width
		maxW := 0
		for _, name := range m.prayerNames {
			if w := lipgloss.Width(name); w > maxW {
				maxW = w
			}
		}

		// Rendering our cursor position.
		for i, name := range m.prayerNames {
			var item string
			if m.cursor == i {
				item = cursorStyle.Render("> ") + selectedItemStyle.Width(maxW+2).Render(name)
			} else {
				item = "  " + normalItemStyle.Width(maxW+2).Render(name)
			}
			lines = append(lines, item)
		}

		lines = append(lines, "")
		lines = append(lines, helpStyle.Render("↑/k up   ↓/j down   enter select   q quit"))
	}

	// vertical centering
	topPad := max(0, (m.termHeight-len(lines))/2)

	var b strings.Builder
	b.WriteString(strings.Repeat("\n", topPad))
	for _, l := range lines {
		b.WriteString(center.Render(l))
		b.WriteString("\n")
	}
	return b.String()
}

// ─── Viewport header / footer ────────────────────────────────────────────────

func (m model) headerView() string {
	title := vpTitleStyle.Render(m.selectedTitle)
	line := strings.Repeat("─", max(0, m.viewport.Width()-lipgloss.Width(title)))
	return lipgloss.JoinHorizontal(lipgloss.Center, title, line)
}

func (m model) footerView() string {
	info := vpFooterStyle.Render(
		fmt.Sprintf("%3.f%%  [esc] back  [q] quit", m.viewport.ScrollPercent()*100),
	)
	line := strings.Repeat("─", max(0, m.viewport.Width()-lipgloss.Width(info)))
	return lipgloss.JoinHorizontal(lipgloss.Center, line, info)
}

// ─── Helpers ─────────────────────────────────────────────────────────────────

// formatDate turns "20260412" into "April 12, 2026".
func formatDate(d string) string {
	if len(d) != 8 {
		return d
	}
	months := [...]string{
		"", "January", "February", "March", "April", "May", "June",
		"July", "August", "September", "October", "November", "December",
	}
	var month, day int
	fmt.Sscanf(d[4:6], "%d", &month)
	fmt.Sscanf(d[6:8], "%d", &day)
	if month < 1 || month > 12 {
		return d
	}
	return fmt.Sprintf("%s %d, %s", months[month], day, d[0:4])
}

// ─── Entry point ─────────────────────────────────────────────────────────────

func main() {
	p := tea.NewProgram(model{})
	if _, err := p.Run(); err != nil {
		fmt.Println("could not run program:", err)
		os.Exit(1)
	}
}
