//go:build windows

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

package platform

import (
	"testing"
	"time"

	"golang.org/x/sys/windows"
)

func Test_filetimeToTime(t *testing.T) {
	type args struct {
		ft windows.Filetime
	}
	tests := []struct {
		name string
		args args
		want Time
	}{
		{
			name: "unix epoch",
			args: args{ft: windows.Filetime{HighDateTime: 27111902, LowDateTime: 3577643008}},
			want: time.Unix(0, 0).UTC(),
		},
		{
			name: "one second after unix epoch",
			args: args{ft: windows.Filetime{HighDateTime: 27111902, LowDateTime: 3587643008}},
			want: time.Unix(1, 0).UTC(),
		},
		{
			name: "half second after unix epoch",
			args: args{ft: windows.Filetime{HighDateTime: 27111902, LowDateTime: 3582643008}},
			want: time.Unix(0, 500000000).UTC(),
		},
		{
			name: "future date 2025-07-13 14:30:00 UTC",
			args: args{ft: windows.Filetime{HighDateTime: 31192066, LowDateTime: 2635326464}},
			want: time.Date(2025, 7, 13, 14, 30, 0, 0, time.UTC),
		},
		{
			name: "pre-unix epoch 1969-12-31 23:59:59 UTC",
			args: args{ft: windows.Filetime{HighDateTime: 27111902, LowDateTime: 3567643008}},
			want: time.Unix(-1, 0).UTC(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := filetimeToTime(tt.args.ft); !got.Equal(tt.want) {
				t.Errorf("filetimeToTime() = %v, want %v", got, tt.want)
			}
		})
	}
}
