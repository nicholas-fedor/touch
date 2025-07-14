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

// Package platform provides platform-specific implementations for timestamp operations.
// It defines exported vars for GetAtime and SetTimesNoDeref, overridden by build tags.
package platform

import (
	"os"
	"time"

	"github.com/nicholas-fedor/touch/internal/errors"
)

// Time is an alias for time.Time, used for clarity in function signatures.
type Time = time.Time

// GetAtime retrieves the access time from file info, platform-specific.
var GetAtime func(os.FileInfo) Time

// SetTimesNoDeref sets times without dereferencing symlinks, platform-specific.
var SetTimesNoDeref func(string, Time, Time) error

// init sets fallback implementations.
func init() {
	GetAtime = func(fileInfo os.FileInfo) Time {
		return fileInfo.ModTime() // Default: use mod time if access unavailable.
	}
	SetTimesNoDeref = func(_ string, _ Time, _ Time) error {
		return errors.ErrNoDerefUnsupported // Default: unsupported.
	}
}
