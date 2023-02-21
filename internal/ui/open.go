package ui

// Reveal reveals the specified file in the Finder.
func (u *UI) Reveal(fname string) {
	reveal(fname)
}
