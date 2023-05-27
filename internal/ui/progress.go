package ui

import (
	"time"

	"github.com/ncruces/zenity"
)

// @todo: this takes ~0.5 seconds. Improve?
func (u *UI) StartProgressBar(id uint, title string) {
	if title == "" {
		// default title
		title = "Loading…"
	}
	p, err := zenity.Progress(zenity.Title(title), zenity.Pulsate(), zenity.Modal())
	if err != nil {
		return
	}
	u.progressBars[id] = p
}

func (u *UI) StopProgressBar(id uint) {
	for try := 0; try < 10; try++ {
		if p := u.progressBars[id]; p != nil {
			p.Complete()
			p.Close()
			p = nil
			return
		}

		println("Progress dialog not found:", id, try)
		time.Sleep(100 * time.Millisecond)
	}
}
