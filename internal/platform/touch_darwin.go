//go:build darwin

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

package platform

import (
	"fmt"
	"os"
	"syscall"
	"time"

	"golang.org/x/sys/unix"
)

// init assigns Darwin-specific implementations for GetAtime and SetTimesNoDeref.
func init() {
	GetAtime = func(fileInfo os.FileInfo) Time {
		if sysStat, ok := fileInfo.Sys().(*syscall.Stat_t); ok {
			return time.Unix(sysStat.Atimespec.Sec, sysStat.Atimespec.Nsec)
		}

		return fileInfo.ModTime() // Fallback if cast fails.
	}

	SetTimesNoDeref = func(file string, accessTime, modTime Time) error {
		timevals := []unix.Timeval{
			{Sec: accessTime.Unix(), Usec: int32(accessTime.UnixMicro() % 1000000)},
			{Sec: modTime.Unix(), Usec: int32(modTime.UnixMicro() % 1000000)},
		}
		if err := unix.Lutimes(file, timevals); err != nil {
			return fmt.Errorf("lutimes %s: %w", file, err)
		}

		return nil
	}
}
