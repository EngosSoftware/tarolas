package server

import (
	"path/filepath"
)

var (
	FlagTrue  = true
	FlagFalse = false
)

// prepareAbsolutePath creates full, absolute path to file or directory
// that is prepended with root directory path. Aditionally the given
// file or directory name is cleaned from all characters like '..'
// or doubled '//'. See filepath.Clean method description for more details.
func prepareAbsolutePath(cfg *Configuration, name string) string {
	return filepath.Join(cfg.RootDirectory, filepath.Clean(name))
}
