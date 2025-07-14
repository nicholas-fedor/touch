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

// Package filesystem defines the FS interface and its default implementation for file operations.
package filesystem

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func Test_defaultFS_Stat(t *testing.T) {
	type args struct {
		path string
	}

	tests := []struct {
		name    string
		d       defaultFS
		args    args
		setup   func() string // Setup temp file if needed.
		wantErr bool
	}{
		{
			name: "existing file",
			d:    defaultFS{},
			setup: func() string {
				tmpFile, _ := os.CreateTemp(t.TempDir(), "test_stat_*")
				defer tmpFile.Close()

				return tmpFile.Name()
			},
			wantErr: false,
		},
		{
			name: "non-existing file",
			d:    defaultFS{},
			args: args{path: "non_existing.txt"},
			setup: func() string {
				return ""
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := tt.args.path
			if tt.setup != nil {
				path = tt.setup()
				if !tt.wantErr {
					defer os.Remove(path) // Cleanup if created.
				}
			}

			tt.args.path = path

			_, err := tt.d.Stat(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("defaultFS.Stat() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_defaultFS_Lstat(t *testing.T) {
	type args struct {
		path string
	}

	tests := []struct {
		name    string
		d       defaultFS
		args    args
		setup   func() string // Setup temp file if needed.
		wantErr bool
	}{
		{
			name: "existing file",
			d:    defaultFS{},
			setup: func() string {
				tmpFile, _ := os.CreateTemp(t.TempDir(), "test_lstat_*")
				defer tmpFile.Close()

				return tmpFile.Name()
			},
			wantErr: false,
		},
		{
			name: "non-existing file",
			d:    defaultFS{},
			args: args{path: "non_existing.txt"},
			setup: func() string {
				return ""
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := tt.args.path
			if tt.setup != nil {
				path = tt.setup()
				if !tt.wantErr {
					defer os.Remove(path) // Cleanup if created.
				}
			}

			tt.args.path = path

			_, err := tt.d.Lstat(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("defaultFS.Lstat() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_defaultFS_Create(t *testing.T) {
	type args struct {
		path string
	}

	tests := []struct {
		name    string
		d       defaultFS
		args    args
		wantErr bool
	}{
		{
			name:    "create new file",
			d:       defaultFS{},
			args:    args{path: filepath.Join(os.TempDir(), "test_create.txt")},
			wantErr: false,
		},
		{
			name:    "create in invalid dir",
			d:       defaultFS{},
			args:    args{path: "/invalid_dir/test_create.txt"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.d.Create(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("defaultFS.Create() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if got != nil {
				got.Close() // Close if created.
			}

			if !tt.wantErr {
				if _, statErr := os.Stat(tt.args.path); statErr != nil {
					t.Errorf("defaultFS.Create() file should exist but stat failed: %v", statErr)
				}

				os.Remove(tt.args.path) // Cleanup after stat.
			}
		})
	}
}

func Test_defaultFS_Chtimes(t *testing.T) {
	type args struct {
		path  string
		atime Time
		mtime Time
	}

	tests := []struct {
		name    string
		d       defaultFS
		args    args
		setup   func() string // Setup temp file.
		wantErr bool
	}{
		{
			name: "change times on existing file",
			d:    defaultFS{},
			setup: func() string {
				tmpFile, _ := os.CreateTemp(t.TempDir(), "test_chtimes_*")
				defer tmpFile.Close()

				return tmpFile.Name()
			},
			args: args{
				atime: time.Date(2025, 7, 13, 14, 0, 0, 0, time.Local),
				mtime: time.Date(2025, 7, 13, 13, 0, 0, 0, time.Local),
			},
			wantErr: false,
		},
		{
			name: "non-existing file",
			d:    defaultFS{},
			args: args{
				path:  "non_existing.txt",
				atime: time.Now(),
				mtime: time.Now(),
			},
			setup: func() string {
				return ""
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := tt.args.path
			if tt.setup != nil {
				path = tt.setup()
				if !tt.wantErr {
					defer os.Remove(path) // Cleanup.
				}
			}

			tt.args.path = path

			err := tt.d.Chtimes(tt.args.path, tt.args.atime, tt.args.mtime)
			if (err != nil) != tt.wantErr {
				t.Errorf("defaultFS.Chtimes() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if tt.wantErr {
				return
			}

			info, statErr := os.Stat(tt.args.path)
			if statErr != nil {
				t.Errorf("defaultFS.Chtimes() stat after = %v", statErr)

				return
			}

			if !info.ModTime().Equal(tt.args.mtime) {
				t.Errorf(
					"defaultFS.Chtimes() mod time = %v, want %v",
					info.ModTime(),
					tt.args.mtime,
				)
			}
		})
	}
}
