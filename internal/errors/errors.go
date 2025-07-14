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

// Package errors defines static error variables for common failure cases.
// These are used to wrap specific errors, enabling error checking with errors.Is.
package errors

import "errors"

// ErrMultipleTimeSources is returned when times are specified from more than one source (-d, -r, -t).
var ErrMultipleTimeSources = errors.New("cannot specify times from more than one source")

// ErrInvalidTimeArg is returned for invalid --time arguments.
var ErrInvalidTimeArg = errors.New("invalid time argument")

// ErrProcessingFiles is returned when errors occur during file processing.
var ErrProcessingFiles = errors.New("errors occurred while processing files")

// ErrInvalidSeconds is returned for invalid seconds in POSIX timestamps.
var ErrInvalidSeconds = errors.New("invalid seconds")

// ErrInvalidPosixLength is returned for invalid length in POSIX timestamps.
var ErrInvalidPosixLength = errors.New("invalid length for POSIX time")

// ErrInvalidDateTimeValues is returned for out-of-range date/time components.
var ErrInvalidDateTimeValues = errors.New("invalid date/time values")

// ErrUnsupportedDateFormat is returned for unsupported date formats in ParseDate.
var ErrUnsupportedDateFormat = errors.New(
	"unsupported date format; try RFC3339, YYYY-MM-DDTHH:MM:SS, YYYY-MM-DD HH:MM:SS, YYYY-MM-DDTHH:MM, YYYY-MM-DD, HH:MM:SS, HH:MM",
)

// ErrNoDerefUnsupported is returned when no-dereference is requested on unsupported platforms.
var ErrNoDerefUnsupported = errors.New("no-dereference not supported on this platform")
