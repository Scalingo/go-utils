package tarball

import (
	"archive/tar"
	"errors"
	"fmt"
	"os"
)

const (
	c_ISUID  = 04000   // Set uid
	c_ISGID  = 02000   // Set gid
	c_ISVTX  = 01000   // Save text (sticky bit)
	c_ISDIR  = 040000  // Directory
	c_ISFIFO = 010000  // FIFO
	c_ISREG  = 0100000 // Regular file
	c_ISLNK  = 0120000 // Symbolic link
	c_ISBLK  = 060000  // Block special file
	c_ISCHR  = 020000  // Character special file
	c_ISSOCK = 0140000 // Socket
)

const (
	blockSize = 512

	// Types
	TypeReg           = '0'    // regular file
	TypeRegA          = '\x00' // regular file
	TypeLink          = '1'    // hard link
	TypeSymlink       = '2'    // symbolic link
	TypeChar          = '3'    // character device node
	TypeBlock         = '4'    // block device node
	TypeDir           = '5'    // directory
	TypeFifo          = '6'    // fifo node
	TypeCont          = '7'    // reserved
	TypeXHeader       = 'x'    // extended header
	TypeXGlobalHeader = 'g'    // global extended header
	TypeGNULongName   = 'L'    // Next file has a long name
	TypeGNULongLink   = 'K'    // Next file symlinks to a file w/ a long name
	TypeGNUSparse     = 'S'    // sparse file
)

// sysStat, if non-nil, populates h from system-dependent fields of fi.
var sysStat func(fi os.FileInfo, h *tar.Header) error

func FileInfoHeader(fi os.FileInfo, link string, fullpath string) (*tar.Header, error) {
	if fi == nil {
		return nil, errors.New("tar: FileInfo is nil")
	}
	fm := fi.Mode()
	h := &tar.Header{
		Name:    link,
		ModTime: fi.ModTime(),
		Mode:    int64(fm.Perm()), // or'd with c_IS* constants later
	}
	switch {
	case fm.IsRegular():
		h.Mode |= c_ISREG
		h.Typeflag = TypeReg
		h.Size = fi.Size()
	case fi.IsDir():
		h.Typeflag = TypeDir
		h.Mode |= c_ISDIR
		h.Name += "/"
	case fm&os.ModeSymlink != 0:
		h.Typeflag = TypeSymlink
		h.Mode |= c_ISLNK
		target, err := os.Readlink(fullpath)
		if err != nil {
			return nil, fmt.Errorf("archive/tar: fail to resolve symlink %v", fullpath)
		}
		h.Linkname = target
	case fm&os.ModeDevice != 0:
		if fm&os.ModeCharDevice != 0 {
			h.Mode |= c_ISCHR
			h.Typeflag = TypeChar
		} else {
			h.Mode |= c_ISBLK
			h.Typeflag = TypeBlock
		}
	case fm&os.ModeNamedPipe != 0:
		h.Typeflag = TypeFifo
		h.Mode |= c_ISFIFO
	case fm&os.ModeSocket != 0:
		h.Mode |= c_ISSOCK
	default:
		return nil, fmt.Errorf("archive/tar: unknown file mode %v", fm)
	}
	if fm&os.ModeSetuid != 0 {
		h.Mode |= c_ISUID
	}
	if fm&os.ModeSetgid != 0 {
		h.Mode |= c_ISGID
	}
	if fm&os.ModeSticky != 0 {
		h.Mode |= c_ISVTX
	}
	// If possible, populate additional fields from OS-specific
	// FileInfo fields.
	if sys, ok := fi.Sys().(*tar.Header); ok {
		// This FileInfo came from a Header (not the OS). Use the
		// original Header to populate all remaining fields.
		h.Uid = sys.Uid
		h.Gid = sys.Gid
		h.Uname = sys.Uname
		h.Gname = sys.Gname
		h.AccessTime = sys.AccessTime
		h.ChangeTime = sys.ChangeTime
		if sys.Xattrs != nil {
			h.Xattrs = make(map[string]string)
			for k, v := range sys.Xattrs {
				h.Xattrs[k] = v
			}
		}
		if sys.Typeflag == TypeLink {
			// hard link
			h.Typeflag = TypeLink
			h.Size = 0
			h.Linkname = sys.Linkname
		}
	}
	if sysStat != nil {
		return h, sysStat(fi, h)
	}
	return h, nil
}
