//go:build windows && crypt

package task

import "github.com/iDigitalFlame/xmt/util/crypt"

var execD = crypt.Get(28) // *.jpg
