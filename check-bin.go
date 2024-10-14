package gosxnotifier

import (
	"fmt"
	"os/exec"
	"runtime"
	"sync"
)

const (
	binName = "terminal-notifier"
)

var (
	binPath   string
	checkErr  error
	checkOnce sync.Once
)

func check() (err error) {
	checkOnce.Do(func() {
		if runtime.GOOS != "darwin" {
			checkErr = fmt.Errorf("%s does not support %s", runtime.GOOS, binName)
			return
		}

		binPath, err = exec.LookPath(binName)
		if err != nil {
			checkErr = fmt.Errorf("could not find %s executable: %w (try $ brew install terminal-notifier)", binName, err)
			return
		}
	})

	return checkErr
}
