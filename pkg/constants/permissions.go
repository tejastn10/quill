package constants

import "io/fs"

const (
	UnreadableFilePerms fs.FileMode = 0
	DirectoryPerms      fs.FileMode = fs.ModePerm & 0750
	ConfigFilePerms     fs.FileMode = fs.ModePerm & 0600
)
