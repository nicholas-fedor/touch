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

// Package timestamp handles reference file timestamp retrieval.
package timestamp

import (
	"os"
	"testing"
	"time"

	"github.com/nicholas-fedor/touch/internal/filesystem"
	"github.com/nicholas-fedor/touch/internal/filesystem/mocks"
	"github.com/nicholas-fedor/touch/internal/platform"
)

// mockFileInfo is a simple mock for os.FileInfo in tests.
type mockFileInfo struct {
	mod Time
	sys any
}

func (m mockFileInfo) Name() string      { return "" }
func (m mockFileInfo) Size() int64       { return 0 }
func (m mockFileInfo) Mode() os.FileMode { return 0 }
func (m mockFileInfo) ModTime() Time     { return m.mod }
func (m mockFileInfo) IsDir() bool       { return false }
func (m mockFileInfo) Sys() any          { return m.sys }

func TestGetTimesFromRef(t *testing.T) {
	type args struct {
		refFilePath string
		noDeref     bool
	}

	tests := []struct {
		name         string
		args         args
		mockSetup    func(*mocks.MockFS)
		mockGetAtime func(os.FileInfo) Time
		wantAccess   Time
		wantMod      Time
		wantErr      bool
	}{
		{
			name: "stat existing no deref false",
			args: args{refFilePath: "testref.txt", noDeref: false},
			mockSetup: func(m *mocks.MockFS) {
				m.On("Stat", "testref.txt").
					Return(&mockFileInfo{mod: time.Date(2025, 7, 13, 13, 0, 0, 0, time.Local)}, nil)
			},
			mockGetAtime: func(_ os.FileInfo) Time {
				return time.Date(2025, 7, 13, 14, 0, 0, 0, time.Local) // Custom access time.
			},
			wantAccess: time.Date(2025, 7, 13, 14, 0, 0, 0, time.Local),
			wantMod:    time.Date(2025, 7, 13, 13, 0, 0, 0, time.Local),
			wantErr:    false,
		},
		{
			name: "lstat existing no deref true",
			args: args{refFilePath: "testref.txt", noDeref: true},
			mockSetup: func(m *mocks.MockFS) {
				m.On("Lstat", "testref.txt").
					Return(&mockFileInfo{mod: time.Date(2025, 7, 13, 12, 0, 0, 0, time.Local)}, nil)
			},
			mockGetAtime: func(_ os.FileInfo) Time {
				return time.Date(2025, 7, 13, 15, 0, 0, 0, time.Local) // Custom access time.
			},
			wantAccess: time.Date(2025, 7, 13, 15, 0, 0, 0, time.Local),
			wantMod:    time.Date(2025, 7, 13, 12, 0, 0, 0, time.Local),
			wantErr:    false,
		},
		{
			name: "stat error",
			args: args{refFilePath: "invalid.txt", noDeref: false},
			mockSetup: func(m *mocks.MockFS) {
				m.On("Stat", "invalid.txt").Return(nil, os.ErrNotExist)
			},
			wantAccess: Time{},
			wantMod:    Time{},
			wantErr:    true,
		},
		{
			name: "lstat error",
			args: args{refFilePath: "invalid.txt", noDeref: true},
			mockSetup: func(m *mocks.MockFS) {
				m.On("Lstat", "invalid.txt").Return(nil, os.ErrPermission)
			},
			wantAccess: Time{},
			wantMod:    Time{},
			wantErr:    true,
		},
		{
			name: "fallback getAtime if platform fails",
			args: args{refFilePath: "testref.txt", noDeref: false},
			mockSetup: func(m *mocks.MockFS) {
				m.On("Stat", "testref.txt").
					Return(&mockFileInfo{mod: time.Date(2025, 7, 13, 13, 0, 0, 0, time.Local)}, nil)
			},
			wantAccess: time.Date(2025, 7, 13, 13, 0, 0, 0, time.Local),
			wantMod:    time.Date(2025, 7, 13, 13, 0, 0, 0, time.Local),
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockFS := mocks.NewMockFS(t)
			if tt.mockSetup != nil {
				tt.mockSetup(mockFS)
			}

			filesystem.Default = mockFS // Override default FS with mock.
			oldGetAtime := platform.GetAtime

			defer func() { platform.GetAtime = oldGetAtime }()

			if tt.mockGetAtime != nil {
				platform.GetAtime = tt.mockGetAtime
			}

			got, got1, err := GetTimesFromRef(tt.args.refFilePath, tt.args.noDeref)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetTimesFromRef() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if !got.Equal(tt.wantAccess) {
				t.Errorf("GetTimesFromRef() got = %v, want %v", got, tt.wantAccess)
			}

			if !got1.Equal(tt.wantMod) {
				t.Errorf("GetTimesFromRef() got1 = %v, want %v", got1, tt.wantMod)
			}
		})
	}
}
