//go:build linux

package ui

import "os/exec"

func reveal(fname string) {
	exec.Command("xdg-open", fname).Run()
}
