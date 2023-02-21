//go:build windows

package ui

import "os/exec"

func reveal(fname string) {
	exec.Command("explorer", fname).Run()
}
