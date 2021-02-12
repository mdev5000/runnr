package running

import (
	"github.com/blang/vfs"
)

func fileExists(fs vfs.Filesystem, path string) bool {
	f, err := fs.Lstat(path)
	if err != nil || f.IsDir() {
		return false
	}
	return true
}
