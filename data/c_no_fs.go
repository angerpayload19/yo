//go:build !go1.16
// +build !go1.16

// Copyright (C) 2020 - 2023 iDigitalFlame
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.
//

package data

import (
	"os"
	"sort"
	"strings"
	"syscall"

	"github.com/iDigitalFlame/xmt/util"
)

type dirEntry struct {
	os.FileInfo
}

// A DirEntry is an entry read from a directory
// (using the ReadDir function or a ReadDirFile's ReadDir method).
//
// Copy of "fs.DirEntry" for go1.15 and below.
//
// This is a pre go1.16 compatibility helper.
type DirEntry interface {
	Name() string
	IsDir() bool
	Type() os.FileMode
	Info() (os.FileInfo, error)
}

func (d dirEntry) Type() os.FileMode {
	return d.Mode()
}

// ReadDir reads the named directory, returning all its directory entries sorted
// by filename. If an error occurs reading the directory, ReadDir returns the
// entries it was able to read before the error, along with the error.
//
// Alias of "os.ReadDir".
//
// This is a pre go1.16 compatibility helper.
func ReadDir(n string) ([]DirEntry, error) {
	f, err := os.OpenFile(n, 0, 0)
	if err != nil {
		return nil, err
	}
	d, err := f.Readdir(0)
	f.Close()
	sort.Slice(d, func(i, j int) bool {
		return d[i].Name() < d[j].Name()
	})
	r := make([]DirEntry, len(d))
	for i := range d {
		r[i] = dirEntry{d[i]}
	}
	return r, nil
}
func (d dirEntry) Info() (os.FileInfo, error) {
	return d.FileInfo, nil
}

// CreateTemp creates a new temporary file in the directory dir, opens the file
// for reading and writing, and returns the resulting file.
//
// The filename is generated by taking pattern and adding a random string to the
// end.
//
// If pattern includes a "*", the random string replaces the last "*".
//
// If dir is the empty string, CreateTemp uses the default directory for temporary
// files, as returned by TempDir.
//
// Multiple programs or goroutines calling CreateTemp simultaneously will not choose
// the same file.
//
// The caller can use the file's Name method to find the pathname of the file.
// It is the caller's responsibility to remove the file when it is no longer needed.
//
// Alias of "os.CreateTemp".
//
// This is a pre go1.16 compatibility helper.
func CreateTemp(d, p string) (*os.File, error) {
	if len(d) == 0 {
		d = os.TempDir()
	}
	s, e, err := prefixAndSuffix(p)
	if err != nil {
		return nil, err
	}
	if os.IsPathSeparator(d[len(d)-1]) {
		s = d + s
	} else {
		s = d + string(os.PathSeparator) + s
	}
	for i := 0; ; {
		// 0xC2 - os.O_RDWR | os.O_CREATE | os.O_EXCL
		f, err := os.OpenFile(s+util.Uitoa(uint64(util.FastRand()))+e, 0xC2, 0600)
		if err != nil && os.IsExist(err) {
			if i++; i < 0x2710 {
				continue
			}
			return nil, syscall.EAGAIN
		}
		return f, err
	}
}
func prefixAndSuffix(v string) (string, string, error) {
	for i := 0; i < len(v); i++ {
		if os.IsPathSeparator(v[i]) {
			return "", "", syscall.EBADF
		}
	}
	if i := strings.LastIndexByte(v, '*'); i != -1 {
		return v[:i], v[i+1:], nil
	}
	return v, "", nil
}