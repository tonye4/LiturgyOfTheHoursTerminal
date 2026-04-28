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

// Need to init our Tabs or we get a runtime panic
func NewTabs() *Tabs {
	return &Tabs{active: [3]Tab{Yesterday, Today, Tomorrow}, currentIndex: 1}
}

func (t *Tabs) CurrentIndex() int {
	return t.currentIndex
}

func (t *Tabs) Next() {
	t.currentIndex++
	if t.currentIndex >= len(t.active) {
		t.currentIndex = 0
	}
}

func (t *Tabs) Prev() {
	t.currentIndex--
	if t.currentIndex < 0 {
		t.currentIndex = 2
	}
}

func (t *Tabs) CurrentTab() Tab {
	return t.active[t.currentIndex]
}
