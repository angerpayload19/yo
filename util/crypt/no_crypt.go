//go:build !crypt
// +build !crypt

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

// Package crypt is a builtin package that provides compile-time encoded string
// values to be decoded and used when first starting up.
//
// This package should only be used with the "crypt" tag, which is auto compiled
// during build.
package crypt

// Get returns the crypt value at the provided string index.
func Get(_ uint8) string {
	return ""
}
