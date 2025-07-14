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
// This file applies the touch operation concurrently to files.
package cli

import (
	"fmt"
	"os"
	"sync"
	"sync/atomic"

	"github.com/nicholas-fedor/touch/internal/core"
	"github.com/nicholas-fedor/touch/internal/errors"
)

// applyToFiles applies the touch operation concurrently to the list of files.
// Uses goroutines for parallel processing; prints errors to stderr and returns an error if any fail.
func applyToFiles(
	changeTimes int,
	noCreate, noDeref bool,
	accessTime, modTime core.Time,
	files []string,
) error {
	var (
		wg       sync.WaitGroup
		hadError atomic.Bool
	)

	for _, file := range files {
		wg.Add(1)

		go func(currentFile string) {
			defer wg.Done()

			if err := core.Touch(currentFile, changeTimes, noCreate, noDeref, accessTime, modTime); err != nil {
				fmt.Fprintf(os.Stderr, "touch: %s: %v\n", core.Quote(currentFile), err)
				hadError.Store(true)
			}
		}(file)
	}

	wg.Wait()

	if hadError.Load() {
		return errors.ErrProcessingFiles
	}

	return nil
}
