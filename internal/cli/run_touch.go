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

// Package cli handles CLI-specific logic, separated from core touch functionality for modularity.
// This file orchestrates the runTouch logic by calling refactored components.
package cli

import (
	"fmt"
	"os"
	"runtime"

	"github.com/spf13/cobra"

	"github.com/nicholas-fedor/touch/internal/errors"
)

// RunTouch is the entry point for the root command's RunE function.
// It processes flags, calculates timestamps, and applies changes to files.
// It handles warnings for obsolete usage or platform-specific limitations.
func RunTouch(cmd *cobra.Command, args []string) error {
	// Process and validate command-line flags.
	changeTimes, noCreate, noDeref, refFilePath, tStamp, dateStr, err := processFlags(cmd)
	if err != nil {
		return err
	}

	// Warn if -h/--no-dereference is used on Windows, where it's unsupported.
	if noDeref && runtime.GOOS == "windows" {
		fmt.Fprintln(
			os.Stderr,
			"Warning: -h/--no-dereference is not supported on Windows; symlinks will be followed",
		)
	}

	// Calculate timestamps and update args if using obsolete format (e.g., `touch 202507131430 file.txt`).
	accessTime, modTime, files, err := calculateTimestamps(
		noDeref,
		refFilePath,
		tStamp,
		dateStr,
		args,
	)
	if err != nil {
		return err
	}

	// If no files are provided, return an error (will trigger usage display).
	if len(files) == 0 {
		return errors.ErrMissingOperands
	}

	// Apply the touch operation to the list of files concurrently.
	return applyToFiles(changeTimes, noCreate, noDeref, accessTime, modTime, files)
}
