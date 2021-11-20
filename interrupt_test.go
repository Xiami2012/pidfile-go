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
	require.Nil(t, SendInterrupt(os.Getpid()))
	sig := <-chsig
	assert.Equal(t, os.Interrupt, sig)
}
