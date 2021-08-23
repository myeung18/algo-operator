package util

import (
	"io/ioutil"
	"os"
)

// CreateTempDir creates a temporary directory and returns
// the directory name and also a function for removing the directory.
// The function is often deferred for directory removal.
func CreateTempDir(baseDirectory string) (string, func()) {
	tmpDir, err := ioutil.TempDir("", baseDirectory)
	if err != nil {
		panic(err)
	}
	return tmpDir, func() {
		if err := os.RemoveAll(tmpDir); err != nil {
			panic(err)
		}
	}
}
