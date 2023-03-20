//go:build windows && !go1.13
// +build windows,!go1.13

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

package winapi

import (
	// Importing unsafe to use "linkname"
	_ "unsafe"
)

func removeCtrlHandler() {
	syscallN(setConsoleCtrlHandler, funcPC(ctrlhandler), 0)
}

//go:linkname ctrlhandler runtime.ctrlhandler
func ctrlhandler(uint32) uint32

//go:linkname funcPC runtime.funcPC
func funcPC(interface{}) uintptr
