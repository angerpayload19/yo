//go:build !go1.14
// +build !go1.14

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

// Package limits contains many options for setting Global limits on how the
// overall application behaves. Many of these options are configured by build tags.
//
// Other functions include re-implemented standard library functions.
package limits

import (
	"os"
	"os/signal"
)

// StopNotify will stop the signal handling loop from running and will cause
// all signal handling to stop.
//
// This function will block until the Goroutine closes.
//
// This function has no effect if the loop is not started or stopped.
//
// The supplied chan can be nil but if non-nil will be passed to 'signal.Stop'
// for convince.
//
// If the Go version is 1.13 or less this function is just a wrapper for
// 'signal.Stop'.
func StopNotify(c chan<- os.Signal) {
	if c == nil {
		return
	}
	signal.Stop(c)
}

// Notify causes package signal to relay incoming signals to c.
// If no signals are provided, all incoming signals will be relayed to c.
// Otherwise, just the provided signals will.
//
// Package signal will not block sending to c: the caller must ensure
// that c has sufficient buffer space to keep up with the expected
// signal rate. For a channel used for notification of just one signal value,
// a buffer of size 1 is sufficient.
//
// It is allowed to call Notify multiple times with the same channel:
// each call expands the set of signals sent to that channel.
// The only way to remove signals from the set is to call Stop.
//
// It is allowed to call Notify multiple times with different channels
// and the same signals: each channel receives copies of incoming
// signals independently.
//
// This version will stop the signal handling loop once the 'StopNotify'
// function has been called.
//
// If the Go version is 1.13 or less this function is just a wrapper for
// 'signal.Notify'.
func Notify(c chan<- os.Signal, s ...os.Signal) {
	signal.Notify(c, s...)
}
