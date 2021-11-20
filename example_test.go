package pidfile_test

import (
	"errors"
	"fmt"
	"os"

	"github.com/Xiami2012/pidfile-go"
)

// The most simple usage is to call Write on startup and Remove on shutdown.
func Example() {
	// In func main
	if err := pidfile.Write("/run/your_instance_name.pid"); err != nil {
		if errors.Is(err, pidfile.ErrPIDFileInUse) {
			fmt.Printf("Instance/Service is already running. Exiting.")
			os.Exit(1)
		}
		panic(fmt.Sprintf("Failed to create PIDFile: %s\n", err.Error()))
	}
	defer pidfile.Remove("/run/your_instance_name.pid")
}

// To avoid writing pidfile's path every time, use PIDFile struct.
func Example_useStruct() {
	// In func main
	mypidfile := pidfile.PIDFile{Filepath: "/run/your_instance_name.pid"}
	if err := mypidfile.Write(); err != nil {
		if errors.Is(err, pidfile.ErrPIDFileInUse) {
			fmt.Printf("Instance/Service is already running. Exiting.")
			os.Exit(1)
		}
		panic(fmt.Sprintf("Failed to create PIDFile: %s\n", err.Error()))
	}
	defer mypidfile.Remove()
}

// Use GetRunningPIDValid to read a PID file and verifies it is valid (process is running and has
// the same executable name with current process)
func Example_getRunningPIDValid() {
	pid, err := pidfile.GetRunningPIDValid("/run/your_instance_name.pid")
	if err == nil {
		fmt.Printf("Found an existing instance/service with pid: %d\n", pid)
	}
}
