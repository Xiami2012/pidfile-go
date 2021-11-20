package pidfile

import (
	"fmt"
	"syscall"
)

// SendInterrupt sends a Ctrl-C event (SIGINT) to the specified process.
//
// It supports both unix and windows.
func SendInterrupt(pid int) error {
	// Code from https://stackoverflow.com/questions/40498371/how-to-send-an-interrupt-signal
	dll, err := syscall.LoadDLL("kernel32.dll")
	if err != nil {
		return fmt.Errorf("LoadDLL: %v", err)
	}
	p, err := dll.FindProc("GenerateConsoleCtrlEvent")
	if err != nil {
		return fmt.Errorf("FindProc: %v", err)
	}
	r, _, err := p.Call(syscall.CTRL_BREAK_EVENT, uintptr(pid))
	if r == 0 {
		return fmt.Errorf("GenerateConsoleCtrlEvent: %v", err)
	}
	return nil
}
