package ui

import "github.com/ncruces/zenity"

func (u *UI) OpenDirectory(id, title string) string {
	if title == "" {
		title = "Choose a directory"
	}

	dir, err := zenity.SelectFile(zenity.Title(title), zenity.Directory(), zenity.Modal())

	if err != nil && err != zenity.ErrCanceled {
		zenity.Error("Error while opening directory", zenity.ErrorIcon)
	}

	return dir
}

func (u *UI) EnterText(title, info, defaultEntry string) string {
	input, err := zenity.Entry(info, zenity.Title(title), zenity.EntryText(defaultEntry), zenity.Modal())

	if err == zenity.ErrCanceled {
		return ""
	} else if err != nil {
		zenity.Error(err.Error(), zenity.Title("Encountered an error!"), zenity.ErrorIcon)
		return ""
	} else {
		return input
	}
}
