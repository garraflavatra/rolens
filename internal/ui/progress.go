package ui

import "github.com/ncruces/zenity"

func (u *UI) StartProgressBar(title string) {
	if u.progress != nil {
		// already loading
		return
	}
	if title == "" {
		// default title
		title = "Loading"
	}
	p, err := zenity.Progress(zenity.Title(title), zenity.Pulsate(), zenity.NoCancel(), zenity.Modal())
	if err != nil {
		return
	}
	u.progress = p
}

func (u *UI) StopProgressBar() {
	if u.progress == nil {
		return
	}
	u.progress.Complete()
	u.progress.Close()
	u.progress = nil
}
