// Build constraints are from os/exec_unix.go

//go:build aix || darwin || dragonfly || freebsd || (js && wasm) || linux || netbsd || openbsd || solaris
// +build aix darwin dragonfly freebsd js,wasm linux netbsd openbsd solaris

package pidfile

import (
	"os"
)

// SendInterrupt sends a Ctrl-C event (SIGINT) to the specified process.
//
// It supports both unix and windows.
func SendInterrupt(pid int) error {
	p, _ := os.FindProcess(pid)
	return p.Signal(os.Interrupt)
}
