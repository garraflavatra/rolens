package ui

import "github.com/ncruces/zenity"

func (u *UI) OpenDirectory(title string) string {
	if title == "" {
		title = "Choose a directory"
	}

	dir, err := zenity.SelectFile(zenity.Title(title), zenity.Directory(), zenity.Modal())

	if err != nil && err != zenity.ErrCanceled {
		zenity.Error("Error while opening directory", zenity.ErrorIcon)
	}

	return dir
}
