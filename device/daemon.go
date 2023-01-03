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

package device

import (
	"context"

	"github.com/iDigitalFlame/xmt/util/xerr"
)

// ErrQuit is an error that can be returned from the DaemonFunction that
// will indicate a clean (non-error) break of the Daemon loop.
var ErrQuit = xerr.Sub("quit", 0x1F)

// DaemonFunc is a function type that can be used as a Daemon. This function
// should return nil to indicate a successful run or ErrQuit to break out of
// a 'DaemonTicker' loop.
//
// Any non-nil errors will be interpreted as exit code '1'.
type DaemonFunc func(context.Context) error
