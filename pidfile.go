// Package pidfile provides a way to manage your PID file.
//
// PID file is a file (usually with .pid file extension) with your running program's PID as its
// content. It is used for both avoiding running multiple same instances and managing the current
// running instance/service (sending signals).
package pidfile

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/mitchellh/go-ps"
	atomicfile "github.com/natefinch/atomic"
)

var (
	ErrEmpty                 = errors.New("empty PIDFile struct (no Filepath given)")
	ErrProcNotExists         = errors.New("process not exists")
	ErrProcNotSameExecutable = errors.New("process is likely not running the same executable")
	ErrPIDFileInUse          = errors.New("pidfile is in use by another process")
)

// PIDFile holds a file path for further related operations.
type PIDFile struct {
	Filepath string
}

func (f PIDFile) checkEmpty() error {
	if f.Filepath == "" {
		return ErrEmpty
	}
	return nil
}

// GetRunningPID reads PID from file and verifies it indicates a running process.
//
// On success, it returns the pid of running process. Otherwise, PID 0 and error is returned.
func (f PIDFile) GetRunningPID() (int, error) {
	if err := f.checkEmpty(); err != nil {
		return 0, err
	}

	fb, err := os.ReadFile(f.Filepath)
	if err != nil {
		return 0, err
	}
	pid, err := strconv.Atoi(string(fb))
	if err != nil {
		return 0, err
	}

	if p, _ := ps.FindProcess(pid); p == nil {
		return 0, ErrProcNotExists
	}

	return pid, nil
}

// GetRunningPIDValid acts like GetRunningPID and additionally verifies it has the same executable
// name with the current process.
func (f PIDFile) GetRunningPIDValid() (int, error) {
	pid, err := f.GetRunningPID()
	if err != nil {
		return 0, err
	}

	p, _ := ps.FindProcess(pid)
	myp, _ := ps.FindProcess(os.Getpid())
	if p.Executable() == myp.Executable() {
		return pid, nil
	}
	return 0, ErrProcNotSameExecutable
}

func (f PIDFile) doWrite(pid int) error {
	return atomicfile.WriteFile(f.Filepath, strings.NewReader(strconv.Itoa(pid)))
}

// Write writes PID of current process to the file if GetRunningPIDValid returns PID 0.
func (f PIDFile) Write() error {
	savedpid, err := f.GetRunningPIDValid()
	if err == nil && savedpid != os.Getpid() {
		return ErrPIDFileInUse
	}
	if savedpid == os.Getpid() {
		return nil
	}

	return f.doWrite(os.Getpid())
}

// WriteForce writes the PID file with PID of current process.
func (f PIDFile) WriteForce() error {
	if err := f.checkEmpty(); err != nil {
		return err
	}

	return f.doWrite(os.Getpid())
}

func (f PIDFile) doRemove() error {
	if err := os.Remove(f.Filepath); err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}
	return nil
}

// Remove removes the PID file if GetRunningPIDValid returns PID of current process.
func (f PIDFile) Remove() error {
	savedpid, err := f.GetRunningPIDValid()
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("pidfile content invalid: %w", err)
	}
	if err == nil && savedpid != os.Getpid() {
		return fmt.Errorf("removing pidfile written by another running processes: %d", savedpid)
	}

	return f.doRemove()
}

// RemoveForce removes the PID file. This is almost the same with os.Remove except it returns nil
// if the PID file does not exist.
func (f PIDFile) RemoveForce() error {
	if err := f.checkEmpty(); err != nil {
		return err
	}

	return f.doRemove()
}

// GetRunningPID constructs a PIDFile and calls its GetRunningPID.
// See PIDFile.GetRunningPID for details.
func GetRunningPID(filepath string) (int, error) {
	return PIDFile{filepath}.GetRunningPID()
}

// GetRunningPIDValid constructs a PIDFile and calls its GetRunningPIDValid.
// See PIDFile.GetRunningPIDValid for details.
func GetRunningPIDValid(filepath string) (int, error) {
	return PIDFile{filepath}.GetRunningPIDValid()
}

// Write constructs a PIDFile and calls its Write.
// See PIDFile.Write for details.
func Write(filepath string) error {
	return PIDFile{filepath}.Write()
}

// WriteForce constructs a PIDFile and calls its WriteForce.
// See PIDFile.WriteForce for details.
func WriteForce(filepath string) error {
	return PIDFile{filepath}.WriteForce()
}

// Remove constructs a PIDFile and calls its Remove.
// See PIDFile.Remove for details.
func Remove(filepath string) error {
	return PIDFile{filepath}.Remove()
}

// RemoveForce constructs a PIDFile and calls its RemoveForce.
// See PIDFile.RemoveForce for details.
func RemoveForce(filepath string) error {
	return PIDFile{filepath}.RemoveForce()
}
