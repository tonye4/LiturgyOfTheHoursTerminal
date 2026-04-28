package tabs

type Tab int

const (
	Today Tab = iota
	Tomorrow
	Yesterday
)

type Tabs struct {
	base         []Tab
	active       []Tab
	currentIndex int
}

func (t *Tabs) Next() {
	t.currentIndex++
	if t.currentIndex >= len(t.active) {
		t.currentIndex = 0
	}
}

func (t *Tabs) Prev() {
	t.currentIndex--
	if t.currentIndex < len(t.active) {
		t.currentIndex = 0
	}
}

func (t *Tabs) CurrentTab() Tab {
	return t.active[t.currentIndex]
}
