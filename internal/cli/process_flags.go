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
// This file processes and validates command-line flags.
package cli

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/spf13/cobra"

	"github.com/nicholas-fedor/touch/internal/core"
	"github.com/nicholas-fedor/touch/internal/errors"
)

// Constants for repeated string values.
const (
	timeAccess = "access"
	timeAtime  = "atime"
	timeUse    = "use"
	timeModify = "modify"
	timeMtime  = "mtime"
	osWindows  = "windows"
)

// processFlags processes and validates command-line flags from the Cobra command.
// It returns the flags as parameters for the touch operation and checks for invalid combinations.
// It also emits warnings for platform-specific limitations (e.g., no-dereference on Windows).
func processFlags(cmd *cobra.Command) (int, bool, bool, string, string, string, error) {
	// Initialize defaults: change both access and modification times.
	changeTimes := core.ChAtime | core.ChMtime

	// Handle -a and -m flags for changing specific timestamps.
	access, _ := cmd.Flags().GetBool("access")
	modification, _ := cmd.Flags().GetBool("modification")
	timeFlag, _ := cmd.Flags().GetString("time")

	// Validate and set changeTimes based on -a, -m, or --time.
	switch {
	case timeFlag != "":
		switch strings.ToLower(timeFlag) {
		case timeAccess, timeAtime, timeUse:
			changeTimes = core.ChAtime
		case timeModify, timeMtime:
			changeTimes = core.ChMtime
		default:
			return 0, false, false, "", "", "", errors.ErrInvalidTimeArg
		}
	case access && !modification:
		changeTimes = core.ChAtime
	case modification && !access:
		changeTimes = core.ChMtime
	}

	// Handle -c/--no-create flag.
	noCreate, _ := cmd.Flags().GetBool("no-create")

	// Handle -h/--no-dereference flag, warn if used on Windows.
	noDeref, _ := cmd.Flags().GetBool("no-dereference")
	if noDeref && runtime.GOOS == osWindows {
		fmt.Fprintln(
			os.Stderr,
			"Warning: -h/--no-dereference is not supported on Windows; symlinks will be followed",
		)

		noDeref = false
	}

	// Handle time source flags: -r, -t, -d.
	refFilePath, _ := cmd.Flags().GetString("reference")
	tStamp, _ := cmd.Flags().GetString("stamp")
	dateStr, _ := cmd.Flags().GetString("date")

	// Check for multiple time sources, which is invalid.
	timeSources := core.BoolToInt(
		refFilePath != "",
	) + core.BoolToInt(
		tStamp != "",
	) + core.BoolToInt(
		dateStr != "",
	)
	if timeSources > 1 {
		return 0, false, false, "", "", "", errors.ErrMultipleTimeSources
	}

	return changeTimes, noCreate, noDeref, refFilePath, tStamp, dateStr, nil
}
