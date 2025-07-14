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

// Package core provides the main Touch function and utilities, orchestrating file timestamp changes.
// It integrates with filesystem, timestamp, and platform subpackages for cross-platform support.
package core

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/nicholas-fedor/touch/internal/filesystem"
	"github.com/nicholas-fedor/touch/internal/platform"
)

// Constants for timestamp change masks.
// These bit flags determine which timestamps (access, modification) to update.
const (
	ChAtime = 1 << iota // Flag to change access time.
	ChMtime             // Flag to change modification time.
)

// Time is an alias for time.Time, used for clarity in function signatures.
type Time = time.Time

// Now is a variable holding the function to get current time, allowing mocking in tests.
var Now = time.Now

// BoolToInt converts a boolean to an integer (1 for true, 0 for false).
// Used for counting active flags in validation.
func BoolToInt(b bool) int {
	if b {
		return 1
	}

	return 0
}

// Quote wraps a string in quotes for safe error message display.
// Mimics shell quoting for filenames with special characters.
func Quote(s string) string {
	return fmt.Sprintf("%q", s)
}

// Touch updates the access and/or modification times of the file at path.
// If the file does not exist and noCreate is false, it creates an empty file.
// The change mask determines which times to update (ChAtime, ChMtime).
// If noDeref is true, it affects symlinks without following them (unsupported on Windows).
// Returns an error if the operation fails.
func Touch(
	file string,
	change int,
	noCreate, noDeref bool,
	accessTimeParam, modTimeParam Time,
) error {
	fileInfo, err := filesystem.Default.Stat(file)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			if noCreate {
				return nil // No creation requested; silently succeed.
			}

			newFile, err := filesystem.Default.Create(file)
			if err != nil {
				return fmt.Errorf("create file %s: %w", file, err)
			}
			defer newFile.Close()
			// Set times on the newly created file.
			if err := filesystem.Default.Chtimes(file, accessTimeParam, modTimeParam); err != nil {
				return fmt.Errorf("chtimes new file %s: %w", file, err)
			}

			return nil
		}

		return fmt.Errorf("stat file %s: %w", file, err)
	}

	// File exists; determine times to set, preserving unchanged ones.
	accessTime := accessTimeParam
	modTime := modTimeParam

	// If not changing access time, retrieve current access time using platform-specific function.
	if change&ChAtime == 0 {
		accessTime = platform.GetAtime(fileInfo)
	}

	// If not changing modification time, use existing ModTime.
	if change&ChMtime == 0 {
		modTime = fileInfo.ModTime()
	}

	// Apply the times.
	if noDeref {
		if err := platform.SetTimesNoDeref(file, accessTime, modTime); err != nil {
			return fmt.Errorf("set times no deref %s: %w", file, err)
		}

		return nil
	}

	if err := filesystem.Default.Chtimes(file, accessTime, modTime); err != nil {
		return fmt.Errorf("chtimes %s: %w", file, err)
	}

	return nil
}
