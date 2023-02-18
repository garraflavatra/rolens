//go:build windows

package open_file

import "os/exec"

func reveal(fname string) {
	exec.Command("explorer", fname).Run()
}
