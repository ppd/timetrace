package state

import (
	"time"

	"fyne.io/fyne/v2/data/binding"
)

type boundTime struct {
	binding.Untyped
}

type TimeBinding interface {
	binding.DataItem
	Get() (time.Time, error)
	Set(time.Time) error
}

func (b *boundTime) Get() (time.Time, error) {
	theTime, error := b.Untyped.Get()
	return theTime.(time.Time), error
}

func (b *boundTime) Set(theTime time.Time) error {
	return b.Untyped.Set(theTime)
}

func NewBoundTime() TimeBinding {
	return &boundTime{
		Untyped: binding.NewUntyped(),
	}
}

func NewBoundTimeWithData(initial time.Time) TimeBinding {
	bt := NewBoundTime()
	bt.Set(initial)
	return bt
}
