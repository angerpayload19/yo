//go:build windows && crypt

package evade

import (
	"github.com/iDigitalFlame/xmt/device/winapi"
	"github.com/iDigitalFlame/xmt/util/crypt"
)

var sect = crypt.Get(104) // .text

func fullPath(n string) string {
	if !isBaseName(n) {
		return n
	}
	d, err := winapi.GetSystemDirectory()
	if err != nil {
		d = crypt.Get(105) // C:\Windows\System32
	}
	return d + "\\" + n
}