//go:build darwin

package open_file

import "os/exec"

func reveal(fname string) {
	exec.Command("open", "--reveal", fname).Run()
}
