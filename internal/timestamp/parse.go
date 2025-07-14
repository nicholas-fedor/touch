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

// Package timestamp handles timestamp parsing for POSIX and flexible date formats.
package timestamp

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/nicholas-fedor/touch/internal/errors"
)

// Constants for POSIX timestamp parsing.
const (
	posixSecondsLength = 2  // Length of [.ss] part.
	posixFullLength    = 12 // Length for CCYYMMDDhhmm.
	posixYearLength    = 10 // Length for YYMMDDhhmm.
	posixMonthLength   = 8  // Length for MMDDhhmm.
	centuryMultiplier  = 100
	y2kBase            = 1900
	y2kPivot           = 69
	y2kShift           = 100
	minMonth           = 1
	maxMonth           = 12
	minDay             = 1
	maxDay             = 31
	minHour            = 0
	maxHour            = 23
	minuteMin          = 0
	minuteMax          = 59
	minSecond          = 0
	maxSecond          = 61 // Allow for leap seconds.
)

// Time is an alias for time.Time, used for clarity in function signatures.
type Time = time.Time

// ParsePosixTime parses the POSIX timestamp format [[CC]YY]MMDDhhmm[.ss].
// Handles century/year variations and validates component ranges.
// Returns a time.Time in the local timezone or an error if invalid.
func ParsePosixTime(timestampStr string) (Time, error) {
	dotIndex := strings.Index(timestampStr, ".")
	second := 0
	if dotIndex != -1 {
		secondsStr := timestampStr[dotIndex+1:]
		if len(secondsStr) != posixSecondsLength {
			return Time{}, fmt.Errorf("%w: %s", errors.ErrInvalidSeconds, secondsStr)
		}
		var err error
		second, err = strconv.Atoi(secondsStr)
		if err != nil {
			return Time{}, fmt.Errorf("atoi seconds: %w", err)
		}
		if second < minSecond || second > maxSecond {
			return Time{}, fmt.Errorf("%w: %d", errors.ErrInvalidSeconds, second)
		}
		timestampStr = timestampStr[:dotIndex]
	}

	length := len(timestampStr)
	var year, month, day, hour, minuteValue int
	var err error
	switch length {
	case posixFullLength: // CCYYMMDDhhmm
		century, err := strconv.Atoi(timestampStr[0:2])
		if err != nil {
			return Time{}, fmt.Errorf("atoi century: %w", err)
		}
		year2, err := strconv.Atoi(timestampStr[2:4])
		if err != nil {
			return Time{}, fmt.Errorf("atoi year2: %w", err)
		}
		year = century*centuryMultiplier + year2
		timestampStr = timestampStr[4:]
	case posixYearLength: // YYMMDDhhmm
		year2, err := strconv.Atoi(timestampStr[0:2])
		if err != nil {
			return Time{}, fmt.Errorf("atoi year2: %w", err)
		}
		year = y2kBase + year2
		if year2 < y2kPivot {
			year += y2kShift
		}
		timestampStr = timestampStr[2:]
	case posixMonthLength: // MMDDhhmm
		year = time.Now().Year()
	default:
		return Time{}, fmt.Errorf("%w: %s", errors.ErrInvalidPosixLength, timestampStr)
	}

	month, err = strconv.Atoi(timestampStr[0:2])
	if err != nil {
		return Time{}, fmt.Errorf("atoi month: %w", err)
	}
	day, err = strconv.Atoi(timestampStr[2:4])
	if err != nil {
		return Time{}, fmt.Errorf("atoi day: %w", err)
	}
	hour, err = strconv.Atoi(timestampStr[4:6])
	if err != nil {
		return Time{}, fmt.Errorf("atoi hour: %w", err)
	}
	minuteValue, err = strconv.Atoi(timestampStr[6:8])
	if err != nil {
		return Time{}, fmt.Errorf("atoi minute: %w", err)
	}

	if month < minMonth || month > maxMonth || day < minDay || day > maxDay || hour < minHour ||
		hour > maxHour ||
		minuteValue < minuteMin ||
		minuteValue > minuteMax {
		return Time{}, errors.ErrInvalidDateTimeValues
	}

	return time.Date(year, time.Month(month), day, hour, minuteValue, second, 0, time.Local), nil
}

// ParseDate parses a date string using predefined formats.
// Supports RFC3339, YYYY-MM-DDTHH:MM:SS, YYYY-MM-DD HH:MM:SS, YYYY-MM-DDTHH:MM, YYYY-MM-DD, HH:MM:SS, HH:MM.
// Assumes local timezone; returns a time.Time or an error if the format is unsupported.
func ParseDate(dateStr string) (Time, error) {
	formats := []string{
		time.RFC3339,
		"2006-01-02T15:04:05",
		"2006-01-02 15:04:05",
		"2006-01-02T15:04",
		"2006-01-02",
		"15:04:05",
		"15:04",
	}
	var parsedTime time.Time
	var parseErr error
	now := time.Now()
	isTimeOnly := false
	for _, format := range formats {
		parsedTime, parseErr = time.ParseInLocation(format, dateStr, time.Local)
		if parseErr == nil {
			if format == "15:04:05" || format == "15:04" {
				isTimeOnly = true
			}

			break
		}
	}
	if parseErr != nil {
		return Time{}, errors.ErrUnsupportedDateFormat
	}
	if isTimeOnly {
		parsedTime = time.Date(
			now.Year(),
			now.Month(),
			now.Day(),
			parsedTime.Hour(),
			parsedTime.Minute(),
			parsedTime.Second(),
			0,
			time.Local,
		)
	}

	return parsedTime, nil
}
