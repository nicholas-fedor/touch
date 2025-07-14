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
	"bytes"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/nicholas-fedor/touch/internal/core"
	"github.com/nicholas-fedor/touch/internal/filesystem"
	"github.com/nicholas-fedor/touch/internal/filesystem/mocks"
	"github.com/nicholas-fedor/touch/internal/platform"
)

func Test_calculateTimestamps(t *testing.T) {
	fixedNow := time.Date(2025, 7, 13, 0, 0, 0, 0, time.Local)
	oldNow := core.Now

	defer func() { core.Now = oldNow }()

	core.Now = func() core.Time { return fixedNow }

	type args struct {
		noDeref     bool
		refFilePath string
		tStamp      string
		dateStr     string
		files       []string
	}

	tests := []struct {
		name        string
		args        args
		mockFSSetup func(*mocks.MockFS)
		setupEnv    func(*testing.T)
		wantAccess  core.Time
		wantMod     core.Time
		wantFiles   []string
		wantErr     bool
		wantStderr  string
	}{
		{
			name: "default current time",
			args: args{
				noDeref:     false,
				refFilePath: "",
				tStamp:      "",
				dateStr:     "",
				files:       []string{},
			},
			mockFSSetup: nil,
			setupEnv:    nil,
			wantAccess:  fixedNow,
			wantMod:     fixedNow,
			wantFiles:   []string{},
			wantErr:     false,
			wantStderr:  "",
		},
		{
			name: "from reference no deref false",
			args: args{
				noDeref:     false,
				refFilePath: "ref.txt",
				tStamp:      "",
				dateStr:     "",
				files:       []string{},
			},
			mockFSSetup: func(m *mocks.MockFS) {
				m.On("Stat", "ref.txt").
					Return(&mockFileInfo{access: time.Date(2025, 7, 13, 14, 0, 0, 0, time.Local), mod: time.Date(2025, 7, 13, 13, 0, 0, 0, time.Local)}, nil)
			},
			setupEnv: func(_ *testing.T) {
				platform.GetAtime = func(fi os.FileInfo) core.Time {
					return fi.(*mockFileInfo).access
				}
			},
			wantAccess: time.Date(2025, 7, 13, 14, 0, 0, 0, time.Local),
			wantMod:    time.Date(2025, 7, 13, 13, 0, 0, 0, time.Local),
			wantFiles:  []string{},
			wantErr:    false,
			wantStderr: "",
		},
		{
			name: "from reference no deref true",
			args: args{
				noDeref:     true,
				refFilePath: "ref.txt",
				tStamp:      "",
				dateStr:     "",
				files:       []string{},
			},
			mockFSSetup: func(m *mocks.MockFS) {
				m.On("Lstat", "ref.txt").
					Return(&mockFileInfo{access: time.Date(2025, 7, 13, 15, 0, 0, 0, time.Local), mod: time.Date(2025, 7, 13, 12, 0, 0, 0, time.Local)}, nil)
			},
			setupEnv: func(_ *testing.T) {
				platform.GetAtime = func(fi os.FileInfo) core.Time {
					return fi.(*mockFileInfo).access
				}
			},
			wantAccess: time.Date(2025, 7, 13, 15, 0, 0, 0, time.Local),
			wantMod:    time.Date(2025, 7, 13, 12, 0, 0, 0, time.Local),
			wantFiles:  []string{},
			wantErr:    false,
			wantStderr: "",
		},
		{
			name: "from stamp",
			args: args{
				noDeref:     false,
				refFilePath: "",
				tStamp:      "2507131430",
				dateStr:     "",
				files:       []string{},
			},
			mockFSSetup: nil,
			setupEnv:    nil,
			wantAccess:  time.Date(2025, 7, 13, 14, 30, 0, 0, time.Local),
			wantMod:     time.Date(2025, 7, 13, 14, 30, 0, 0, time.Local),
			wantFiles:   []string{},
			wantErr:     false,
			wantStderr:  "",
		},
		{
			name: "from date full",
			args: args{
				noDeref:     false,
				refFilePath: "",
				tStamp:      "",
				dateStr:     "2025-07-13 14:30:00",
				files:       []string{},
			},
			mockFSSetup: nil,
			setupEnv:    nil,
			wantAccess:  time.Date(2025, 7, 13, 14, 30, 0, 0, time.Local),
			wantMod:     time.Date(2025, 7, 13, 14, 30, 0, 0, time.Local),
			wantFiles:   []string{},
			wantErr:     false,
			wantStderr:  "",
		},
		{
			name: "from date time only",
			args: args{
				noDeref:     false,
				refFilePath: "",
				tStamp:      "",
				dateStr:     "14:30:00",
				files:       []string{},
			},
			mockFSSetup: nil,
			setupEnv:    nil,
			wantAccess:  time.Date(2025, 7, 14, 14, 30, 0, 0, time.Local),
			wantMod:     time.Date(2025, 7, 14, 14, 30, 0, 0, time.Local),
			wantFiles:   []string{},
			wantErr:     false,
			wantStderr:  "",
		},
		{
			name: "obsolete usage success",
			args: args{
				noDeref:     false,
				refFilePath: "",
				tStamp:      "",
				dateStr:     "",
				files:       []string{"2507131430", "file1.txt", "file2.txt"},
			},
			mockFSSetup: nil,
			setupEnv:    nil,
			wantAccess:  time.Date(2025, 7, 13, 14, 30, 0, 0, time.Local),
			wantMod:     time.Date(2025, 7, 13, 14, 30, 0, 0, time.Local),
			wantFiles:   []string{"file1.txt", "file2.txt"},
			wantErr:     false,
			wantStderr:  "warning: 'touch 2507131430' is obsolete; use 'touch -t'\n",
		},
		{
			name: "obsolete usage with POSIXLY_CORRECT no warn",
			args: args{
				noDeref:     false,
				refFilePath: "",
				tStamp:      "",
				dateStr:     "",
				files:       []string{"2507131430", "file1.txt"},
			},
			mockFSSetup: nil,
			setupEnv: func(t *testing.T) {
				t.Helper()
				t.Setenv("POSIXLY_CORRECT", "1")
			},
			wantAccess: time.Date(2025, 7, 13, 14, 30, 0, 0, time.Local),
			wantMod:    time.Date(2025, 7, 13, 14, 30, 0, 0, time.Local),
			wantFiles:  []string{"file1.txt"},
			wantErr:    false,
			wantStderr: "",
		},
		{
			name: "obsolete invalid fallback current",
			args: args{
				noDeref:     false,
				refFilePath: "",
				tStamp:      "",
				dateStr:     "",
				files:       []string{"invalid", "file1.txt"},
			},
			mockFSSetup: nil,
			setupEnv:    nil,
			wantAccess:  fixedNow,
			wantMod:     fixedNow,
			wantFiles:   []string{"invalid", "file1.txt"},
			wantErr:     false,
			wantStderr:  "",
		},
		{
			name: "error from ref",
			args: args{
				noDeref:     false,
				refFilePath: "invalid_ref.txt",
				tStamp:      "",
				dateStr:     "",
				files:       []string{},
			},
			mockFSSetup: func(m *mocks.MockFS) {
				m.On("Stat", "invalid_ref.txt").Return(nil, os.ErrNotExist)
			},
			setupEnv:   nil,
			wantAccess: core.Time{},
			wantMod:    core.Time{},
			wantFiles:  nil,
			wantErr:    true,
			wantStderr: "",
		},
		{
			name: "error from stamp",
			args: args{
				noDeref:     false,
				refFilePath: "",
				tStamp:      "invalid",
				dateStr:     "",
				files:       []string{},
			},
			mockFSSetup: nil,
			setupEnv:    nil,
			wantAccess:  core.Time{},
			wantMod:     core.Time{},
			wantFiles:   nil,
			wantErr:     true,
			wantStderr:  "",
		},
		{
			name: "error from date",
			args: args{
				noDeref:     false,
				refFilePath: "",
				tStamp:      "",
				dateStr:     "invalid",
				files:       []string{},
			},
			mockFSSetup: nil,
			setupEnv:    nil,
			wantAccess:  core.Time{},
			wantMod:     core.Time{},
			wantFiles:   nil,
			wantErr:     true,
			wantStderr:  "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockFS := mocks.NewMockFS(t)
			if tt.mockFSSetup != nil {
				tt.mockFSSetup(mockFS)
			}

			filesystem.Default = mockFS // Override default FS with mock.

			// Setup env if needed, defer unset not needed with t.Setenv.
			if tt.setupEnv != nil {
				tt.setupEnv(t)
			}

			// Capture stderr.
			oldStderr := os.Stderr
			r, w, _ := os.Pipe()
			os.Stderr = w

			got, got1, got2, err := calculateTimestamps(
				tt.args.noDeref,
				tt.args.refFilePath,
				tt.args.tStamp,
				tt.args.dateStr,
				tt.args.files,
			)

			w.Close()

			os.Stderr = oldStderr

			var buf bytes.Buffer
			buf.ReadFrom(r)
			stderrOutput := buf.String()

			if (err != nil) != tt.wantErr {
				t.Errorf("calculateTimestamps() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if !got.Equal(tt.wantAccess) {
				t.Errorf("calculateTimestamps() got = %v, want %v", got, tt.wantAccess)
			}

			if !got1.Equal(tt.wantMod) {
				t.Errorf("calculateTimestamps() got1 = %v, want %v", got1, tt.wantMod)
			}

			if !reflect.DeepEqual(got2, tt.wantFiles) {
				t.Errorf("calculateTimestamps() got2 = %v, want %v", got2, tt.wantFiles)
			}

			if stderrOutput != tt.wantStderr {
				t.Errorf("calculateTimestamps() stderr = %v, want %v", stderrOutput, tt.wantStderr)
			}
		})
	}
}
