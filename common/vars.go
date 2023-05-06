package common

import "os"

//nolint:gochecknoglobals,gomnd
var (
	DirPerm  = os.FileMode(0o755)
	FilePerm = os.FileMode(0o644)
)
