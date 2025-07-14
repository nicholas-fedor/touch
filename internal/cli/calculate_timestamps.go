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
// This file computes timestamps based on provided sources.
package cli

import (
	"fmt"
	"os"

	"github.com/nicholas-fedor/touch/internal/core"
	"github.com/nicholas-fedor/touch/internal/timestamp"
)

// calculateTimestamps computes the access and modification times based on flags and args.
// Handles reference, stamp, date, obsolete usage, or defaults to current time.
// Returns the computed times and updated files list or an error.
func calculateTimestamps(
	noDeref bool,
	refFilePath, tStamp, dateStr string,
	files []string,
) (core.Time, core.Time, []string, error) {
	var accessTime, modTime core.Time

	dateSet := false

	var err error

	// Use switch to determine timestamp source, addressing ifElseChain lint rule.
	switch {
	case refFilePath != "":
		accessTime, modTime, err = timestamp.GetTimesFromRef(refFilePath, noDeref)
		if err != nil {
			return core.Time{}, core.Time{}, nil, fmt.Errorf("get reference times: %w", err)
		}

		dateSet = true
	case tStamp != "":
		accessTime, err = timestamp.ParsePosixTime(tStamp)
		if err != nil {
			return core.Time{}, core.Time{}, nil, fmt.Errorf("parse POSIX stamp: %w", err)
		}

		modTime = accessTime
		dateSet = true
	case dateStr != "":
		newTime, err := timestamp.ParseDate(dateStr)
		if err != nil {
			return core.Time{}, core.Time{}, nil, fmt.Errorf("parse date: %w", err)
		}

		accessTime = newTime
		modTime = newTime
		dateSet = true
	}

	// Handle obsolete usage if no source set: treat first arg as POSIX timestamp.
	if !dateSet && len(files) >= 1 {
		t, err := timestamp.ParsePosixTime(files[0])
		if err == nil {
			accessTime = t
			modTime = t
			dateSet = true

			if os.Getenv("POSIXLY_CORRECT") == "" {
				fmt.Fprintf(
					os.Stderr,
					"warning: 'touch %s' is obsolete; use 'touch -t'\n",
					files[0],
				)
			}

			files = files[1:]
		}
	}

	// Default to current time if still not set.
	if !dateSet {
		now := core.Now()
		accessTime = now
		modTime = now
	}

	return accessTime, modTime, files, nil
}
