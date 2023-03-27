package ui

import "fyne.io/fyne/v2/widget"

type BlurrableEntry struct {
	widget.Entry
	OnFocusLost func()
}

func (e *BlurrableEntry) FocusLost() {
	e.Entry.FocusLost()
	if e.OnFocusLost != nil {
		e.OnFocusLost()
	}
}

func NewBlurrableEntry() *BlurrableEntry {
	entry := &BlurrableEntry{}
	entry.ExtendBaseWidget(entry)
	return entry
}
