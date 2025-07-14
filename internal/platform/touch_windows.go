//go:build windows

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
	"os"
	"time"

	"golang.org/x/sys/windows"
)

// Constants for filetimeToTime calculations.
const (
	HighDateTimeShift       = 32
	FiletimeToUnixDivisor   = 10000000
	FiletimeToUnixRemainder = 100
	EpochOffset100ns        = 116444736000000000 // 100ns intervals from 1601 to 1970.
)

// init assigns Windows-specific implementations for GetAtime and SetTimesNoDeref.
func init() {
	GetAtime = func(fileInfo os.FileInfo) Time {
		if winStat, ok := fileInfo.Sys().(*windows.Win32FileAttributeData); ok {
			return filetimeToTime(winStat.LastAccessTime)
		}

		return fileInfo.ModTime() // Fallback if cast fails.
	}
}

// filetimeToTime converts a Windows Filetime to time.Time.
// Filetime represents the number of 100-nanosecond intervals since January 1, 1601 UTC.
func filetimeToTime(ft windows.Filetime) Time {
	nsec := int64(ft.HighDateTime)<<HighDateTimeShift + int64(ft.LowDateTime)
	nsec -= EpochOffset100ns
	sec := nsec / FiletimeToUnixDivisor
	nsec = (nsec % FiletimeToUnixDivisor) * FiletimeToUnixRemainder

	return time.Unix(sec, nsec)
}
