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

	"github.com/spf13/cobra"
)

// RunTouch is the entry point for the root command's RunE.
// It orchestrates flag processing, timestamp calculation, and file application.
// Returns an error if any step fails.
func RunTouch(cmd *cobra.Command, args []string) error {
	changeTimes, noCreate, noDeref, refFilePath, tStamp, dateStr, err := processFlags(cmd)
	if err != nil {
		return err
	}

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

	if len(files) == 0 {
		return fmt.Errorf("missing operands")
	}

	return applyToFiles(changeTimes, noCreate, noDeref, accessTime, modTime, files)
}
