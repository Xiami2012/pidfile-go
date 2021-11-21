//go:build !windows
// +build !windows

// Since code of SendInterrupt on Windows is from Go standard library test and it is already
// tested by Golang, we do not make another test here.

package pidfile

import (
	"os"
	"os/signal"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInterruptSignal(t *testing.T) {
	chsig := make(chan os.Signal, 1)
	signal.Notify(chsig, os.Interrupt)
	defer signal.Reset(os.Interrupt)
	require.Nil(t, SendInterrupt(os.Getpid()))
	sig := <-chsig
	assert.Equal(t, os.Interrupt, sig)
}
