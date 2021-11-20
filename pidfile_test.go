package pidfile

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func getTempFilePath() string {
	f, err := os.CreateTemp("", "pidfile-")
	if err != nil {
		panic(err)
	}
	f.Close()
	res := f.Name()
	//runtime.SetFinalizer(&res, func(s *string) { os.Remove(*s) })
	return res
}

func TestPIDFile(t *testing.T) {
	fp := getTempFilePath()
	defer os.Remove(fp)
	assert.NotNil(t, Remove(fp))
	assert.Nil(t, Write(fp))
	pid, err := GetRunningPIDValid(fp)
	assert.Nil(t, err)
	assert.EqualValues(t, os.Getpid(), pid)
	assert.Nil(t, Write(fp))
	assert.Nil(t, Remove(fp))
}

func TestOtherProcess(t *testing.T) {
	fp := getTempFilePath()
	defer os.Remove(fp)
	// Empty file is invalid
	assert.NotNil(t, Remove(fp))
	assert.Nil(t, os.Remove(fp))
	// Remove a non-existing file is ok
	assert.Nil(t, Remove(fp))
	assert.Nil(t, os.WriteFile(fp, []byte("1"), 0644))
	pid, err := GetRunningPID(fp)
	assert.Nil(t, err)
	assert.EqualValues(t, 1, pid)
	_, err = GetRunningPIDValid(fp)
	assert.NotNil(t, err)
}

func TestInvalidFileContent(t *testing.T) {
	fp := getTempFilePath()
	defer os.Remove(fp)
	assert.Nil(t, os.WriteFile(fp, []byte("2147483647"), 0644))
	_, err := GetRunningPID(fp)
	assert.NotNil(t, err)
	assert.NotNil(t, Remove(fp))
	assert.Nil(t, os.WriteFile(fp, []byte("0xABCDEFGH"), 0644))
	_, err = GetRunningPID(fp)
	assert.NotNil(t, err)
	assert.NotNil(t, Remove(fp))
}

func TestInvalidPIDFile(t *testing.T) {
	assert.NotNil(t, PIDFile{}.Write())
	assert.NotNil(t, WriteForce(""))
	assert.NotNil(t, RemoveForce(""))
	assert.NotNil(t, WriteForce("/notexists"))
	assert.Nil(t, RemoveForce("/notexists"))
	wd, err := os.Getwd()
	assert.Nil(t, err)
	_, err = PIDFile{wd}.GetRunningPIDValid()
	assert.NotNil(t, err)
}
