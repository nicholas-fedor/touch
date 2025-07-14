//go:build !windows && !darwin

/*
Copyright © 2025 Nicholas Fedor <nick@nickfedor.com>

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/

// Package platform provides platform-specific implementations for timestamp operations.
// It defines exported vars for GetAtime and SetTimesNoDeref, overridden by build tags.
package platform

import (
	"fmt"
	"os"
	"syscall"
	"time"

	"golang.org/x/sys/unix"
)

// init assigns Unix-specific (non-Darwin) implementations for GetAtime and SetTimesNoDeref.
func init() {
	GetAtime = func(fileInfo os.FileInfo) Time {
		if sysStat, ok := fileInfo.Sys().(*syscall.Stat_t); ok {
			// Cast to int64 to support 32-bit architectures (386, arm) where Sec and Nsec are int32.
			// On 64-bit systems, these are already int64, but the cast is safe and avoids type errors.
			//nolint:unconvert // Necessary for 32-bit compatibility.
			return time.Unix(int64(sysStat.Atim.Sec), int64(sysStat.Atim.Nsec))
		}

		return fileInfo.ModTime() // Fallback if cast fails.
	}

	SetTimesNoDeref = func(file string, accessTime, modTime Time) error {
		ts := []unix.Timespec{
			unix.NsecToTimespec(accessTime.UnixNano()),
			unix.NsecToTimespec(modTime.UnixNano()),
		}
		if err := unix.UtimesNanoAt(unix.AT_FDCWD, file, ts, unix.AT_SYMLINK_NOFOLLOW); err != nil {
			return fmt.Errorf("utimesnanoat %s: %w", file, err)
		}

		return nil
	}
}
