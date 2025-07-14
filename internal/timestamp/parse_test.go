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
	"testing"
	"time"
)

func TestParsePosixTime(t *testing.T) {
	type args struct {
		timestampStr string
	}
	tests := []struct {
		name    string
		args    args
		want    Time
		wantErr bool
	}{
		{
			name:    "full format with seconds",
			args:    args{timestampStr: "202507131430.30"},
			want:    time.Date(2025, 7, 13, 14, 30, 30, 0, time.Local),
			wantErr: false,
		},
		{
			name:    "full format without seconds",
			args:    args{timestampStr: "202507131430"},
			want:    time.Date(2025, 7, 13, 14, 30, 0, 0, time.Local),
			wantErr: false,
		},
		{
			name:    "year format with century pivot",
			args:    args{timestampStr: "2507131430"},
			want:    time.Date(2025, 7, 13, 14, 30, 0, 0, time.Local),
			wantErr: false,
		},
		{
			name:    "year format below pivot",
			args:    args{timestampStr: "6807131430"},
			want:    time.Date(2068, 7, 13, 14, 30, 0, 0, time.Local),
			wantErr: false,
		},
		{
			name:    "month format",
			args:    args{timestampStr: "07131430"},
			want:    time.Date(2025, 7, 13, 14, 30, 0, 0, time.Local), // Using current year 2025.
			wantErr: false,
		},
		{
			name:    "invalid length short",
			args:    args{timestampStr: "0713143"},
			want:    Time{},
			wantErr: true,
		},
		{
			name:    "invalid length long",
			args:    args{timestampStr: "2025071314300"},
			want:    Time{},
			wantErr: true,
		},
		{
			name:    "invalid seconds value high",
			args:    args{timestampStr: "202507131430.99"},
			want:    Time{},
			wantErr: true,
		},
		{
			name:    "invalid seconds value low",
			args:    args{timestampStr: "202507131430.-1"},
			want:    Time{},
			wantErr: true,
		},
		{
			name:    "invalid seconds format",
			args:    args{timestampStr: "202507131430.abc"},
			want:    Time{},
			wantErr: true,
		},
		{
			name:    "invalid month high",
			args:    args{timestampStr: "202513131430"},
			want:    Time{},
			wantErr: true,
		},
		{
			name:    "invalid day low",
			args:    args{timestampStr: "202507001430"},
			want:    Time{},
			wantErr: true,
		},
		{
			name:    "invalid hour high",
			args:    args{timestampStr: "202507132430"},
			want:    Time{},
			wantErr: true,
		},
		{
			name:    "invalid minute low",
			args:    args{timestampStr: "2025071314-1"},
			want:    Time{},
			wantErr: true,
		},
		{
			name:    "non-numeric input",
			args:    args{timestampStr: "abcd"},
			want:    Time{},
			wantErr: true,
		},
		{
			name:    "leap second allowed",
			args:    args{timestampStr: "202507131430.60"},
			want:    time.Date(2025, 7, 13, 14, 30, 60, 0, time.Local),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParsePosixTime(tt.args.timestampStr)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParsePosixTime() error = %v, wantErr %v", err, tt.wantErr)

				return
			}
			if !got.Equal(tt.want) {
				t.Errorf("ParsePosixTime() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseDate(t *testing.T) {
	type args struct {
		dateStr string
	}
	tests := []struct {
		name    string
		args    args
		want    Time
		wantErr bool
	}{
		{
			name:    "RFC3339 full",
			args:    args{dateStr: "2025-07-13T14:30:00Z"},
			want:    time.Date(2025, 7, 13, 14, 30, 0, 0, time.UTC),
			wantErr: false,
		},
		{
			name:    "YYYY-MM-DDTHH:MM:SS",
			args:    args{dateStr: "2025-07-13T14:30:00"},
			want:    time.Date(2025, 7, 13, 14, 30, 0, 0, time.Local),
			wantErr: false,
		},
		{
			name:    "YYYY-MM-DD HH:MM:SS",
			args:    args{dateStr: "2025-07-13 14:30:00"},
			want:    time.Date(2025, 7, 13, 14, 30, 0, 0, time.Local),
			wantErr: false,
		},
		{
			name:    "YYYY-MM-DDTHH:MM",
			args:    args{dateStr: "2025-07-13T14:30"},
			want:    time.Date(2025, 7, 13, 14, 30, 0, 0, time.Local),
			wantErr: false,
		},
		{
			name:    "YYYY-MM-DD",
			args:    args{dateStr: "2025-07-13"},
			want:    time.Date(2025, 7, 13, 0, 0, 0, 0, time.Local),
			wantErr: false,
		},
		{
			name:    "HH:MM:SS",
			args:    args{dateStr: "14:30:00"},
			want:    time.Date(2025, 7, 13, 14, 30, 0, 0, time.Local), // Using current date.
			wantErr: false,
		},
		{
			name:    "HH:MM",
			args:    args{dateStr: "14:30"},
			want:    time.Date(2025, 7, 13, 14, 30, 0, 0, time.Local), // Using current date.
			wantErr: false,
		},
		{
			name:    "invalid format",
			args:    args{dateStr: "2025/07/13"},
			want:    Time{},
			wantErr: true,
		},
		{
			name:    "empty string",
			args:    args{dateStr: ""},
			want:    Time{},
			wantErr: true,
		},
		{
			name:    "invalid date values",
			args:    args{dateStr: "2025-13-13T14:30:00"},
			want:    Time{},
			wantErr: true,
		},
		{
			name:    "leap year date",
			args:    args{dateStr: "2024-02-29"},
			want:    time.Date(2024, 2, 29, 0, 0, 0, 0, time.Local),
			wantErr: false,
		},
		{
			name:    "non-leap year invalid",
			args:    args{dateStr: "2025-02-29"},
			want:    Time{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseDate(tt.args.dateStr)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseDate() error = %v, wantErr %v", err, tt.wantErr)

				return
			}
			if !got.Equal(tt.want) {
				t.Errorf("ParseDate() got = %v, want %v", got, tt.want)
			}
		})
	}
}
