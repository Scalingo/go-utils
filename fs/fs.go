package fs

import (
	"os"

	"github.com/spf13/afero"
)

// Fs interface contains all methods implemented in MemFs and OsFs
type Fs interface {
	afero.Fs
	Chown(string, int, int) error
	Link(string, string) error
	Symlink(string, string) error
}

type MemFs struct {
	afero.Fs
}

func NewMemFs() MemFs {
	return MemFs{Fs: afero.NewMemMapFs()}
}

func (fs MemFs) Chown(string, int, int) error {
	return nil
}

func (fs MemFs) Link(string, string) error {
	return nil
}

func (fs MemFs) Symlink(string, string) error {
	return nil
}

type OsFs struct {
	afero.Fs
}

func NewOsFs() OsFs {
	return OsFs{Fs: afero.NewOsFs()}
}

func (fs OsFs) Chown(path string, uid int, gid int) error {
	return os.Chown(path, uid, gid)
}

func (fs OsFs) Link(old, new string) error {
	return os.Link(old, new)
}

func (fs OsFs) Symlink(old, new string) error {
	return os.Symlink(old, new)
}
