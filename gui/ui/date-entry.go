package ui

import (
	"fmt"
	"time"

	"fyne.io/fyne/v2/data/binding"
)

type DateEntry struct {
	BlurrableEntry

	time TimeBinding
}

func (de *DateEntry) getDateItems() (int, time.Month, int) {
	storedTime, _ := de.time.Get()
	year, month, day := storedTime.Date()
	return year, month, day
}

func (de *DateEntry) getTimeItems() (int, int) {
	storedTime, _ := de.time.Get()
	return storedTime.Hour(), storedTime.Minute()
}

func (de *DateEntry) refresh() {
	year, month, day := de.getDateItems()
	de.SetText(fmt.Sprintf("%d-%02d-%02d", year, month, day))
}

func (de *DateEntry) syncDateToBinding() {
	theDate, err := time.Parse("2006-01-02", de.Text)
	if err == nil {
		hours, minutes := de.getTimeItems()
		theDate = theDate.Add(time.Hour * time.Duration(hours))
		theDate = theDate.Add(time.Minute * time.Duration(minutes))
		de.time.Set(theDate)
	} else {
		de.refresh()
	}
}

func NewDateEntry(theTime TimeBinding) *DateEntry {
	entry := &DateEntry{
		time: theTime,
	}
	entry.ExtendBaseWidget(entry)

	entry.time.AddListener(binding.NewDataListener(entry.refresh))
	entry.OnSubmitted = func(s string) { entry.syncDateToBinding() }
	entry.OnFocusLost = func() { entry.syncDateToBinding() }

	return entry
}
