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

// Package core provides the main Touch function and utilities, orchestrating file timestamp changes.
// It integrates with filesystem, timestamp, and platform subpackages for cross-platform support.
package core

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"

	"github.com/nicholas-fedor/touch/internal/errors"
	"github.com/nicholas-fedor/touch/internal/filesystem"
	"github.com/nicholas-fedor/touch/internal/filesystem/mocks"
	"github.com/nicholas-fedor/touch/internal/platform"
)

func TestBoolToInt(t *testing.T) {
	type args struct {
		b bool
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "true to 1",
			args: args{b: true},
			want: 1,
		},
		{
			name: "false to 0",
			args: args{b: false},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := BoolToInt(tt.args.b); got != tt.want {
				t.Errorf("BoolToInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestQuote(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "simple string no special chars",
			args: args{s: "testfile.txt"},
			want: "\"testfile.txt\"",
		},
		{
			name: "string with space",
			args: args{s: "test file.txt"},
			want: "\"test file.txt\"",
		},
		{
			name: "string with quote",
			args: args{s: "test\"file.txt"},
			want: "\"test\\\"file.txt\"",
		},
		{
			name: "string with backslash",
			args: args{s: "test\\file.txt"},
			want: "\"test\\\\file.txt\"",
		},
		{
			name: "empty string",
			args: args{s: ""},
			want: "\"\"",
		},
		{
			name: "string with newline",
			args: args{s: "test\nfile.txt"},
			want: "\"test\\nfile.txt\"",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Quote(tt.args.s); got != tt.want {
				t.Errorf("Quote() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTouch(t *testing.T) {
	type args struct {
		file            string
		change          int
		noCreate        bool
		noDeref         bool
		accessTimeParam Time
		modTimeParam    Time
	}
	tests := []struct {
		name           string
		args           args
		mockFSSetup    func(*mocks.MockFS)
		mockGetAtime   func(os.FileInfo) Time
		mockSetNoDeref func(string, Time, Time) error
		wantErr        bool
	}{
		{
			name: "touch existing change both",
			args: args{
				file:            "existing.txt",
				change:          ChAtime | ChMtime,
				noCreate:        false,
				noDeref:         false,
				accessTimeParam: time.Date(2025, 7, 13, 14, 0, 0, 0, time.Local),
				modTimeParam:    time.Date(2025, 7, 13, 13, 0, 0, 0, time.Local),
			},
			mockFSSetup: func(m *mocks.MockFS) {
				m.On("Stat", "existing.txt").
					Return(&mockFileInfo{mod: time.Date(2025, 7, 13, 12, 0, 0, 0, time.Local)}, nil)
				m.On("Chtimes", "existing.txt", mock.AnythingOfType("time.Time"), mock.AnythingOfType("time.Time")).
					Return(nil)
			},
			mockGetAtime:   nil, // Use default.
			mockSetNoDeref: nil,
			wantErr:        false,
		},
		{
			name: "touch existing change only atime",
			args: args{
				file:            "existing.txt",
				change:          ChAtime,
				noCreate:        false,
				noDeref:         false,
				accessTimeParam: time.Date(2025, 7, 13, 14, 0, 0, 0, time.Local),
				modTimeParam:    time.Date(2025, 7, 13, 13, 0, 0, 0, time.Local), // Ignored.
			},
			mockFSSetup: func(m *mocks.MockFS) {
				m.On("Stat", "existing.txt").
					Return(&mockFileInfo{mod: time.Date(2025, 7, 13, 12, 0, 0, 0, time.Local)}, nil)
				m.On("Chtimes", "existing.txt", mock.AnythingOfType("time.Time"), time.Date(2025, 7, 13, 12, 0, 0, 0, time.Local)).
					Return(nil)
			},
			mockGetAtime:   nil,
			mockSetNoDeref: nil,
			wantErr:        false,
		},
		{
			name: "touch existing change only mtime",
			args: args{
				file:            "existing.txt",
				change:          ChMtime,
				noCreate:        false,
				noDeref:         false,
				accessTimeParam: time.Date(2025, 7, 13, 14, 0, 0, 0, time.Local), // Ignored.
				modTimeParam:    time.Date(2025, 7, 13, 13, 0, 0, 0, time.Local),
			},
			mockFSSetup: func(m *mocks.MockFS) {
				m.On("Stat", "existing.txt").
					Return(&mockFileInfo{mod: time.Date(2025, 7, 13, 12, 0, 0, 0, time.Local)}, nil)
				m.On("Chtimes", "existing.txt", time.Date(2025, 7, 13, 11, 0, 0, 0, time.Local), mock.AnythingOfType("time.Time")).
					Return(nil)
			},
			mockGetAtime: func(fi os.FileInfo) Time {
				return time.Date(2025, 7, 13, 11, 0, 0, 0, time.Local)
			},
			mockSetNoDeref: nil,
			wantErr:        false,
		},
		{
			name: "create new file",
			args: args{
				file:            "new.txt",
				change:          ChAtime | ChMtime,
				noCreate:        false,
				noDeref:         false,
				accessTimeParam: time.Date(2025, 7, 13, 14, 0, 0, 0, time.Local),
				modTimeParam:    time.Date(2025, 7, 13, 13, 0, 0, 0, time.Local),
			},
			mockFSSetup: func(m *mocks.MockFS) {
				m.On("Stat", "new.txt").Return(nil, os.ErrNotExist)
				m.On("Create", "new.txt").Return(&os.File{}, nil)
				m.On("Chtimes", "new.txt", mock.AnythingOfType("time.Time"), mock.AnythingOfType("time.Time")).
					Return(nil)
			},
			mockGetAtime:   nil,
			mockSetNoDeref: nil,
			wantErr:        false,
		},
		{
			name: "no create on missing",
			args: args{
				file:            "missing.txt",
				change:          ChAtime | ChMtime,
				noCreate:        true,
				noDeref:         false,
				accessTimeParam: time.Date(2025, 7, 13, 14, 0, 0, 0, time.Local),
				modTimeParam:    time.Date(2025, 7, 13, 13, 0, 0, 0, time.Local),
			},
			mockFSSetup: func(m *mocks.MockFS) {
				m.On("Stat", "missing.txt").Return(nil, os.ErrNotExist)
			},
			mockGetAtime:   nil,
			mockSetNoDeref: nil,
			wantErr:        false,
		},
		{
			name: "no deref unsupported",
			args: args{
				file:            "symlink.txt",
				change:          ChAtime | ChMtime,
				noCreate:        false,
				noDeref:         true,
				accessTimeParam: time.Date(2025, 7, 13, 14, 0, 0, 0, time.Local),
				modTimeParam:    time.Date(2025, 7, 13, 13, 0, 0, 0, time.Local),
			},
			mockFSSetup: func(m *mocks.MockFS) {
				m.On("Stat", "symlink.txt").
					Return(&mockFileInfo{mod: time.Date(2025, 7, 13, 12, 0, 0, 0, time.Local)}, nil)
			},
			mockGetAtime: nil,
			mockSetNoDeref: func(string, Time, Time) error {
				return errors.ErrNoDerefUnsupported
			},
			wantErr: true,
		},
		{
			name: "error on stat",
			args: args{
				file:            "error.txt",
				change:          ChAtime | ChMtime,
				noCreate:        false,
				noDeref:         false,
				accessTimeParam: time.Date(2025, 7, 13, 14, 0, 0, 0, time.Local),
				modTimeParam:    time.Date(2025, 7, 13, 13, 0, 0, 0, time.Local),
			},
			mockFSSetup: func(m *mocks.MockFS) {
				m.On("Stat", "error.txt").Return(nil, os.ErrPermission)
			},
			mockGetAtime:   nil,
			mockSetNoDeref: nil,
			wantErr:        true,
		},
		{
			name: "error on create",
			args: args{
				file:            "new_error.txt",
				change:          ChAtime | ChMtime,
				noCreate:        false,
				noDeref:         false,
				accessTimeParam: time.Date(2025, 7, 13, 14, 0, 0, 0, time.Local),
				modTimeParam:    time.Date(2025, 7, 13, 13, 0, 0, 0, time.Local),
			},
			mockFSSetup: func(m *mocks.MockFS) {
				m.On("Stat", "new_error.txt").Return(nil, os.ErrNotExist)
				m.On("Create", "new_error.txt").Return(nil, os.ErrPermission)
			},
			mockGetAtime:   nil,
			mockSetNoDeref: nil,
			wantErr:        true,
		},
		{
			name: "error on chtimes existing",
			args: args{
				file:            "existing_error.txt",
				change:          ChAtime | ChMtime,
				noCreate:        false,
				noDeref:         false,
				accessTimeParam: time.Date(2025, 7, 13, 14, 0, 0, 0, time.Local),
				modTimeParam:    time.Date(2025, 7, 13, 13, 0, 0, 0, time.Local),
			},
			mockFSSetup: func(m *mocks.MockFS) {
				m.On("Stat", "existing_error.txt").
					Return(&mockFileInfo{mod: time.Date(2025, 7, 13, 12, 0, 0, 0, time.Local)}, nil)
				m.On("Chtimes", "existing_error.txt", mock.AnythingOfType("time.Time"), mock.AnythingOfType("time.Time")).
					Return(os.ErrPermission)
			},
			mockGetAtime:   nil,
			mockSetNoDeref: nil,
			wantErr:        true,
		},
		{
			name: "error on chtimes new",
			args: args{
				file:            "new_chtimes_error.txt",
				change:          ChAtime | ChMtime,
				noCreate:        false,
				noDeref:         false,
				accessTimeParam: time.Date(2025, 7, 13, 14, 0, 0, 0, time.Local),
				modTimeParam:    time.Date(2025, 7, 13, 13, 0, 0, 0, time.Local),
			},
			mockFSSetup: func(m *mocks.MockFS) {
				m.On("Stat", "new_chtimes_error.txt").Return(nil, os.ErrNotExist)
				m.On("Create", "new_chtimes_error.txt").Return(&os.File{}, nil)
				m.On("Chtimes", "new_chtimes_error.txt", mock.AnythingOfType("time.Time"), mock.AnythingOfType("time.Time")).
					Return(os.ErrPermission)
			},
			mockGetAtime:   nil,
			mockSetNoDeref: nil,
			wantErr:        true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockFS := mocks.NewMockFS(t)
			if tt.mockFSSetup != nil {
				tt.mockFSSetup(mockFS)
			}
			filesystem.Default = mockFS // Override default FS with mock.
			if tt.mockGetAtime != nil {
				oldGetAtime := platform.GetAtime
				platform.GetAtime = tt.mockGetAtime
				defer func() { platform.GetAtime = oldGetAtime }()
			}
			if tt.mockSetNoDeref != nil {
				oldSetNoDeref := platform.SetTimesNoDeref
				platform.SetTimesNoDeref = tt.mockSetNoDeref
				defer func() { platform.SetTimesNoDeref = oldSetNoDeref }()
			}

			err := Touch(
				tt.args.file,
				tt.args.change,
				tt.args.noCreate,
				tt.args.noDeref,
				tt.args.accessTimeParam,
				tt.args.modTimeParam,
			)
			if (err != nil) != tt.wantErr {
				t.Errorf("Touch() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// mockFileInfo is a simple mock for os.FileInfo in tests.
type mockFileInfo struct {
	mod Time
}

func (m mockFileInfo) Name() string      { return "" }
func (m mockFileInfo) Size() int64       { return 0 }
func (m mockFileInfo) Mode() os.FileMode { return 0 }
func (m mockFileInfo) ModTime() Time     { return m.mod }
func (m mockFileInfo) IsDir() bool       { return false }
func (m mockFileInfo) Sys() any          { return nil }
