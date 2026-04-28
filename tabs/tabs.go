package tabs

type Tab int

const (
	Today Tab = iota
	Yesterday
	Tomorrow
)

type Tabs struct {
	active       [3]Tab
	currentIndex int
}

func NewTabs() *Tabs {
	return &Tabs{active: [3]Tab{}}
}

func (t *Tabs) Next() {
	t.currentIndex++
	if t.currentIndex >= len(t.active) {
		t.currentIndex = 1
	}
}

func (t *Tabs) CurrentTab() Tab {
	return t.active[t.currentIndex]
}
