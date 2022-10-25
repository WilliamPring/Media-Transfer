package filesystem

import (
	"os"
)

func checkFilePath() {
	if _, err := os.Stat("/path/to/whatever"); !os.IsNotExist(err) {
		// path/to/whatever exists
	}
}
