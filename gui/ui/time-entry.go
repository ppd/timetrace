package ui

import (
	"fmt"
	"time"

	"fyne.io/fyne/v2/data/binding"
	"github.com/dominikbraun/timetrace/gui/state"
)

type TimeEntry struct {
	BlurrableEntry
	time state.TimeBinding
}

func (te *TimeEntry) getDateItems() (int, time.Month, int) {
	storedTime, _ := te.time.Get()
	year, month, day := storedTime.Date()
	return year, month, day
}

func (te *TimeEntry) getTimeItems() (int, int) {
	storedTime, _ := te.time.Get()
	return storedTime.Hour(), storedTime.Minute()
}

func (te *TimeEntry) syncTimeToBinding() {
	theTime, err := time.Parse("15:04", te.Text)
	if err == nil {
		year, month, day := te.getDateItems()
		theDate := time.Date(year, month, day, theTime.Hour(), theTime.Minute(), 0, 0, time.Local)
		te.time.Set(theDate)
	} else {
		te.refresh()
	}
}

func (te *TimeEntry) refresh() {
	hour, minute := te.getTimeItems()
	te.SetText(fmt.Sprintf("%02d:%02d", hour, minute))
}

func NewTimeEntry(theTime state.TimeBinding) *TimeEntry {
	entry := &TimeEntry{
		time: theTime,
	}
	entry.ExtendBaseWidget(entry)

	entry.time.AddListener(binding.NewDataListener(entry.refresh))
	entry.OnSubmitted = func(s string) { entry.syncTimeToBinding() }
	entry.OnFocusLost = func() { entry.syncTimeToBinding() }

	return entry
}
