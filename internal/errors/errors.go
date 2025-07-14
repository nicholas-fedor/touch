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

// Package errors defines custom error types used across the touch CLI tool.
// It provides specific error values for consistent error handling in various packages.
package errors

import "errors"

// ErrInvalidDateTimeValues indicates that the provided date or time components are out of valid ranges.
var ErrInvalidDateTimeValues = errors.New("invalid date or time values")

// ErrInvalidPosixLength indicates that the POSIX timestamp string has an invalid length.
var ErrInvalidPosixLength = errors.New("invalid POSIX timestamp length")

// ErrInvalidSeconds indicates that the seconds component in a POSIX timestamp is invalid.
var ErrInvalidSeconds = errors.New("invalid seconds value")

// ErrInvalidTimeArg indicates that the --time flag received an invalid argument.
var ErrInvalidTimeArg = errors.New("invalid time argument")

// ErrMissingOperands indicates that no files were provided as arguments when required.
var ErrMissingOperands = errors.New("missing operands")

// ErrMultipleTimeSources indicates that multiple time source flags (-r, -t, -d) were specified simultaneously.
var ErrMultipleTimeSources = errors.New("multiple time sources specified")

// ErrNoDerefUnsupported indicates that the --no-dereference option is not supported on the current platform.
var ErrNoDerefUnsupported = errors.New("no-dereference is not supported on this platform")

// ErrProcessingFiles indicates that errors occurred while processing one or more files.
var ErrProcessingFiles = errors.New("errors occurred while processing files")

// ErrUnsupportedDateFormat indicates that the provided date string does not match any supported format.
var ErrUnsupportedDateFormat = errors.New("unsupported date format")
