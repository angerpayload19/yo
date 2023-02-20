//go:build implant && !loader
// +build implant,!loader

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

package cmd

// LoaderEnabled is a compile-time constant that is true if the "DLLToASM" function
// will modify the provided bytes slice to ASM, otherwise this will be false.
const LoaderEnabled = false

// DLLToASM will patch the DLL raw bytes and convert it into shellcode
// using the SRDi launcher.
//
//	SRDi GitHub: https://github.com/monoxgas/sRDI
//
// The first string param is the function name which can be empty if not
// needed.
//
// The resulting byte slice can be used in an 'Asm' struct to directly load and
// run the DLL.
func DLLToASM(_ string, b []byte) []byte {
	return b
}
