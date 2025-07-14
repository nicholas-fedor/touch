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
	"os"
	"runtime"
	"strings"

	"github.com/spf13/cobra"

	"github.com/nicholas-fedor/touch/internal/core"
	"github.com/nicholas-fedor/touch/internal/errors"
)

// processFlags retrieves and validates flags from the Cobra command.
// It computes the changeTimes mask and handles flag conflicts.
// Returns the processed values or an error if validation fails.
func processFlags(cmd *cobra.Command) (int, bool, bool, string, string, string, error) {
	changeAccess, _ := cmd.Flags().GetBool("access")
	changeMod, _ := cmd.Flags().GetBool("modification")
	timeStr, _ := cmd.Flags().GetString("time")
	noCreate, _ := cmd.Flags().GetBool("no-create")
	noDeref, _ := cmd.Flags().GetBool("no-dereference")
	refFilePath, _ := cmd.Flags().GetString("reference")
	tStamp, _ := cmd.Flags().GetString("stamp")
	dateStr, _ := cmd.Flags().GetString("date")

	// Handle no-dereference on Windows: warn and disable if set, as it's unsupported.
	if noDeref && runtime.GOOS == "windows" {
		os.Stderr.WriteString(
			"Warning: -h/--no-dereference is not supported on Windows; symlinks will be followed\n",
		)

		noDeref = false
	}

	// Compute changeTimes mask based on flags.
	changeTimes := 0
	if changeAccess {
		changeTimes |= core.ChAtime
	}

	if changeMod {
		changeTimes |= core.ChMtime
	}

	if timeStr != "" {
		switch strings.ToLower(timeStr) {
		case "access", "atime", "use":
			changeTimes = core.ChAtime
		case "modify", "mtime":
			changeTimes = core.ChMtime
		default:
			return 0, false, false, "", "", "", errors.ErrInvalidTimeArg
		}
	}

	if changeTimes == 0 {
		changeTimes = core.ChAtime | core.ChMtime
	}

	// Validate exclusive use of time sources.
	useRef := refFilePath != ""
	flexDate := dateStr != ""

	hasTStamp := tStamp != ""
	if core.BoolToInt(useRef)+core.BoolToInt(flexDate)+core.BoolToInt(hasTStamp) > 1 {
		return 0, false, false, "", "", "", errors.ErrMultipleTimeSources
	}

	return changeTimes, noCreate, noDeref, refFilePath, tStamp, dateStr, nil
}
