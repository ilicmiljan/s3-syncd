package local

import "os"

type TempFile struct {
	Path    string
	File    *os.File
	Cleanup func()
}
