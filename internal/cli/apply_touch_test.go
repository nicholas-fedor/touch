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
	"bytes"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"

	"github.com/nicholas-fedor/touch/internal/core"
	"github.com/nicholas-fedor/touch/internal/filesystem"
	"github.com/nicholas-fedor/touch/internal/filesystem/mocks"
)

func Test_applyToFiles(t *testing.T) {
	type args struct {
		changeTimes int
		noCreate    bool
		noDeref     bool
		accessTime  core.Time
		modTime     core.Time
		files       []string
	}

	tests := []struct {
		name        string
		args        args
		mockFSSetup func(*mocks.MockFS)
		wantErr     bool
		wantStderr  string
	}{
		{
			name: "no files",
			args: args{
				changeTimes: core.ChAtime | core.ChMtime,
				noCreate:    false,
				noDeref:     false,
				accessTime:  time.Now(),
				modTime:     time.Now(),
				files:       []string{},
			},
			mockFSSetup: nil,
			wantErr:     false,
			wantStderr:  "",
		},
		{
			name: "single file success",
			args: args{
				changeTimes: core.ChAtime | core.ChMtime,
				noCreate:    false,
				noDeref:     false,
				accessTime:  time.Date(2025, 7, 13, 14, 0, 0, 0, time.Local),
				modTime:     time.Date(2025, 7, 13, 13, 0, 0, 0, time.Local),
				files:       []string{"testfile.txt"},
			},
			mockFSSetup: func(m *mocks.MockFS) {
				m.On("Stat", "testfile.txt").
					Return(&mockFileInfo{mod: time.Date(2025, 7, 13, 12, 0, 0, 0, time.Local)}, nil)
				m.On("Chtimes", "testfile.txt", mock.AnythingOfType("time.Time"), mock.AnythingOfType("time.Time")).
					Return(nil)
			},
			wantErr:    false,
			wantStderr: "",
		},
		{
			name: "multiple files success",
			args: args{
				changeTimes: core.ChAtime | core.ChMtime,
				noCreate:    false,
				noDeref:     false,
				accessTime:  time.Date(2025, 7, 13, 14, 0, 0, 0, time.Local),
				modTime:     time.Date(2025, 7, 13, 13, 0, 0, 0, time.Local),
				files:       []string{"file1.txt", "file2.txt"},
			},
			mockFSSetup: func(m *mocks.MockFS) {
				m.On("Stat", "file1.txt").
					Return(&mockFileInfo{mod: time.Date(2025, 7, 13, 12, 0, 0, 0, time.Local)}, nil)
				m.On("Chtimes", "file1.txt", mock.AnythingOfType("time.Time"), mock.AnythingOfType("time.Time")).
					Return(nil)
				m.On("Stat", "file2.txt").
					Return(&mockFileInfo{mod: time.Date(2025, 7, 13, 11, 0, 0, 0, time.Local)}, nil)
				m.On("Chtimes", "file2.txt", mock.AnythingOfType("time.Time"), mock.AnythingOfType("time.Time")).
					Return(nil)
			},
			wantErr:    false,
			wantStderr: "",
		},
		{
			name: "single file error",
			args: args{
				changeTimes: core.ChAtime | core.ChMtime,
				noCreate:    false,
				noDeref:     false,
				accessTime:  time.Date(2025, 7, 13, 14, 0, 0, 0, time.Local),
				modTime:     time.Date(2025, 7, 13, 13, 0, 0, 0, time.Local),
				files:       []string{"errorfile.txt"},
			},
			mockFSSetup: func(m *mocks.MockFS) {
				m.On("Stat", "errorfile.txt").Return(nil, os.ErrPermission)
			},
			wantErr:    true,
			wantStderr: "touch: \"errorfile.txt\": stat file errorfile.txt: permission denied\n",
		},
		{
			name: "multiple files one error",
			args: args{
				changeTimes: core.ChAtime | core.ChMtime,
				noCreate:    false,
				noDeref:     false,
				accessTime:  time.Date(2025, 7, 13, 14, 0, 0, 0, time.Local),
				modTime:     time.Date(2025, 7, 13, 13, 0, 0, 0, time.Local),
				files:       []string{"file1.txt", "errorfile.txt"},
			},
			mockFSSetup: func(m *mocks.MockFS) {
				m.On("Stat", "file1.txt").
					Return(&mockFileInfo{mod: time.Date(2025, 7, 13, 12, 0, 0, 0, time.Local)}, nil)
				m.On("Chtimes", "file1.txt", mock.AnythingOfType("time.Time"), mock.AnythingOfType("time.Time")).
					Return(nil)
				m.On("Stat", "errorfile.txt").Return(nil, os.ErrPermission)
			},
			wantErr:    true,
			wantStderr: "touch: \"errorfile.txt\": stat file errorfile.txt: permission denied\n",
		},
		{
			name: "create new file",
			args: args{
				changeTimes: core.ChAtime | core.ChMtime,
				noCreate:    false,
				noDeref:     false,
				accessTime:  time.Date(2025, 7, 13, 14, 0, 0, 0, time.Local),
				modTime:     time.Date(2025, 7, 13, 13, 0, 0, 0, time.Local),
				files:       []string{"newfile.txt"},
			},
			mockFSSetup: func(m *mocks.MockFS) {
				m.On("Stat", "newfile.txt").Return(nil, os.ErrNotExist)
				m.On("Create", "newfile.txt").Return(&os.File{}, nil)
				m.On("Chtimes", "newfile.txt", mock.AnythingOfType("time.Time"), mock.AnythingOfType("time.Time")).
					Return(nil)
			},
			wantErr:    false,
			wantStderr: "",
		},
		{
			name: "no create missing file",
			args: args{
				changeTimes: core.ChAtime | core.ChMtime,
				noCreate:    true,
				noDeref:     false,
				accessTime:  time.Date(2025, 7, 13, 14, 0, 0, 0, time.Local),
				modTime:     time.Date(2025, 7, 13, 13, 0, 0, 0, time.Local),
				files:       []string{"missing.txt"},
			},
			mockFSSetup: func(m *mocks.MockFS) {
				m.On("Stat", "missing.txt").Return(nil, os.ErrNotExist)
			},
			wantErr:    false,
			wantStderr: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockFS := mocks.NewMockFS(t)
			if tt.mockFSSetup != nil {
				tt.mockFSSetup(mockFS)
			}

			filesystem.Default = mockFS // Override default FS with mock.

			// Capture stderr.
			oldStderr := os.Stderr
			r, w, _ := os.Pipe()
			os.Stderr = w

			err := applyToFiles(
				tt.args.changeTimes,
				tt.args.noCreate,
				tt.args.noDeref,
				tt.args.accessTime,
				tt.args.modTime,
				tt.args.files,
			)

			w.Close()

			os.Stderr = oldStderr

			var buf bytes.Buffer
			buf.ReadFrom(r)
			stderrOutput := buf.String()

			if (err != nil) != tt.wantErr {
				t.Errorf("applyToFiles() error = %v, wantErr %v", err, tt.wantErr)
			}

			if stderrOutput != tt.wantStderr {
				t.Errorf("applyToFiles() stderr = %v, want %v", stderrOutput, tt.wantStderr)
			}
		})
	}
}
