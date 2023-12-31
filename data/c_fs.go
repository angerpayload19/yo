//go:build go1.16
// +build go1.16

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

package data

import (
	"io/fs"
	"os"
)

// A DirEntry is an entry read from a directory
// (using the ReadDir function or a ReadDirFile's ReadDir method).
//
// Alias of "fs.DirEntry".
//
// This is a pre go1.16 compatibility helper.
type DirEntry = fs.DirEntry

// ReadDir reads the named directory, returning all its directory entries sorted
// by filename. If an error occurs reading the directory, ReadDir returns the
// entries it was able to read before the error, along with the error.
//
// Alias of "os.ReadDir".
//
// This is a pre go1.16 compatibility helper.
func ReadDir(n string) ([]DirEntry, error) {
	return os.ReadDir(n)
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
	return os.CreateTemp(d, p)
}
