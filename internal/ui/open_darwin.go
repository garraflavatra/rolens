//go:build darwin

package ui

import "os/exec"

func reveal(fname string) {
	exec.Command("open", "--reveal", fname).Run()
}
