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

// Package timestamp handles reference file timestamp retrieval.
package timestamp

import (
	"fmt"
	"os"

	"github.com/nicholas-fedor/touch/internal/filesystem"
	"github.com/nicholas-fedor/touch/internal/platform"
)

// GetTimesFromRef retrieves the access and modification times from a reference file.
// If noDeref is true, it uses Lstat to avoid following symlinks.
// Uses platform-specific GetAtime for access time; returns times or an error.
func GetTimesFromRef(refFilePath string, noDeref bool) (Time, Time, error) {
	var (
		fileInfo os.FileInfo
		err      error
	)

	if noDeref {
		fileInfo, err = filesystem.Default.Lstat(refFilePath)
	} else {
		fileInfo, err = filesystem.Default.Stat(refFilePath)
	}

	if err != nil {
		return Time{}, Time{}, fmt.Errorf("get file info for %s: %w", refFilePath, err)
	}

	modTime := fileInfo.ModTime()
	accessTime := platform.GetAtime(fileInfo)

	return accessTime, modTime, nil
}
