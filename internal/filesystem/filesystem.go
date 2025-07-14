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

// Package filesystem defines the FS interface and its default implementation for file operations.
package filesystem

import (
	"fmt"
	"os"
	"time"
)

// Time is an alias for time.Time, used for clarity in function signatures.
type Time = time.Time

// FS abstracts file system operations for testability and modularity.
type FS interface {
	Stat(
		path string,
	) (info os.FileInfo, err error) // Retrieves file info, following path symlinks.
	Lstat(
		path string,
	) (info os.FileInfo, err error) // Retrieves file info without following path symlinks.
	Create(path string) (file *os.File, err error) // Creates a new file at path.
	Chtimes(
		path string,
		atime Time,
		mtime Time,
	) error // Changes path's access and mod times, following symlinks.
}

// defaultFS is the default implementation using os package functions.
type defaultFS struct{}

// Default is the default file system implementation, using standard os functions.
var Default FS = defaultFS{}

// Stat implements FS.Stat using os.Stat.
func (defaultFS) Stat(path string) (os.FileInfo, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("stat %s: %w", path, err)
	}

	return info, nil
}

// Lstat implements FS.Lstat using os.Lstat.
func (defaultFS) Lstat(path string) (os.FileInfo, error) {
	info, err := os.Lstat(path)
	if err != nil {
		return nil, fmt.Errorf("lstat %s: %w", path, err)
	}

	return info, nil
}

// Create implements FS.Create using os.Create.
func (defaultFS) Create(path string) (*os.File, error) {
	file, err := os.Create(path)
	if err != nil {
		return nil, fmt.Errorf("create %s: %w", path, err)
	}

	return file, nil
}

// Chtimes implements FS.Chtimes using os.Chtimes.
func (defaultFS) Chtimes(path string, atime Time, mtime Time) error {
	if err := os.Chtimes(path, atime, mtime); err != nil {
		return fmt.Errorf("chtimes %s: %w", path, err)
	}

	return nil
}
