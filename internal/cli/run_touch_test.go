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
// This file orchestrates the runTouch logic by calling refactored components.
package cli

import (
	"bytes"
	"os"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/mock"

	"github.com/nicholas-fedor/touch/internal/filesystem"
	"github.com/nicholas-fedor/touch/internal/filesystem/mocks"
	"github.com/nicholas-fedor/touch/internal/platform"
)

func TestRunTouch(t *testing.T) {
	type args struct {
		cmd  *cobra.Command
		args []string
	}

	tests := []struct {
		name        string
		args        args
		mockFSSetup func(*mocks.MockFS)
		setupEnv    func()
		wantErr     bool
		wantStdout  string
		wantStderr  string
	}{
		{
			name: "no files show usage",
			args: args{
				cmd:  createTestCmd(),
				args: []string{},
			},
			mockFSSetup: nil,
			setupEnv:    nil,
			wantErr:     true,
			wantStdout:  "",
			wantStderr:  "",
		},
		{
			name: "single file default current",
			args: args{
				cmd:  createTestCmd(),
				args: []string{"file.txt"},
			},
			mockFSSetup: func(m *mocks.MockFS) {
				m.On("Stat", "file.txt").Return(&mockFileInfo{mod: time.Now()}, nil)
				m.On("Chtimes", "file.txt", mock.AnythingOfType("time.Time"), mock.AnythingOfType("time.Time")).
					Return(nil)
			},
			setupEnv:   nil,
			wantErr:    false,
			wantStdout: "",
			wantStderr: "",
		},
		{
			name: "error from processFlags invalid time",
			args: args{
				cmd: createTestCmd(
					func(cmd *cobra.Command) { cmd.Flags().Set("time", "invalid") },
				),
				args: []string{},
			},
			mockFSSetup: nil,
			setupEnv:    nil,
			wantErr:     true,
			wantStdout:  "",
			wantStderr:  "",
		},
		{
			name: "error from calculateTimestamps invalid date",
			args: args{
				cmd: createTestCmd(
					func(cmd *cobra.Command) { cmd.Flags().Set("date", "invalid") },
				),
				args: []string{"file.txt"},
			},
			mockFSSetup: nil,
			setupEnv:    nil,
			wantErr:     true,
			wantStdout:  "",
			wantStderr:  "",
		},
		{
			name: "error from applyToFiles",
			args: args{
				cmd:  createTestCmd(),
				args: []string{"errorfile.txt"},
			},
			mockFSSetup: func(m *mocks.MockFS) {
				m.On("Stat", "errorfile.txt").Return(nil, os.ErrPermission)
			},
			setupEnv:   nil,
			wantErr:    true,
			wantStdout: "",
			wantStderr: "touch: \"errorfile.txt\": stat file errorfile.txt: permission denied\n",
		},
		{
			name: "obsolete usage with warning",
			args: args{
				cmd:  createTestCmd(),
				args: []string{"2507131430", "file.txt"},
			},
			mockFSSetup: func(m *mocks.MockFS) {
				m.On("Stat", "file.txt").Return(&mockFileInfo{mod: time.Now()}, nil)
				m.On("Chtimes", "file.txt", mock.AnythingOfType("time.Time"), mock.AnythingOfType("time.Time")).
					Return(nil)
			},
			setupEnv:   nil,
			wantErr:    false,
			wantStdout: "",
			wantStderr: "warning: 'touch 2507131430' is obsolete; use 'touch -t'\n",
		},
		{
			name: "no deref on windows warning",
			args: args{
				cmd: createTestCmd(
					func(cmd *cobra.Command) { cmd.Flags().Set("no-dereference", "true") },
				),
				args: []string{"file.txt"},
			},
			mockFSSetup: func(m *mocks.MockFS) {
				m.On("Stat", "file.txt").Return(nil, os.ErrNotExist)
				m.On("Create", "file.txt").Return(&os.File{}, nil)
				m.On("Chtimes", "file.txt", mock.AnythingOfType("time.Time"), mock.AnythingOfType("time.Time")).
					Return(nil)
			},
			setupEnv: func() {
				if runtime.GOOS != osWindows {
					// On non-Windows, mock SetTimesNoDeref for no-dereference behavior.
					platform.SetTimesNoDeref = func(_ string, _, _ time.Time) error {
						return nil
					}
				}
			},
			wantErr:    false,
			wantStdout: "",
			wantStderr: func() string {
				if runtime.GOOS == osWindows {
					return "Warning: -h/--no-dereference is not supported on Windows; symlinks will be followed\n"
				}

				return ""
			}(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockFS := mocks.NewMockFS(t)
			if tt.mockFSSetup != nil {
				tt.mockFSSetup(mockFS)
			}

			filesystem.Default = mockFS // Override default FS with mock.

			// Setup env if needed.
			if tt.setupEnv != nil {
				tt.setupEnv()
			}

			// Capture stdout and stderr.
			oldStdout := os.Stdout
			oldStderr := os.Stderr
			rOut, wOut, _ := os.Pipe()
			rErr, wErr, _ := os.Pipe()
			os.Stdout = wOut
			os.Stderr = wErr

			err := RunTouch(tt.args.cmd, tt.args.args)

			wOut.Close()
			wErr.Close()

			os.Stdout = oldStdout
			os.Stderr = oldStderr

			var bufOut, bufErr bytes.Buffer
			bufOut.ReadFrom(rOut)
			bufErr.ReadFrom(rErr)

			stdoutOutput := bufOut.String()
			stderrOutput := strings.ReplaceAll(bufErr.String(), "\r\n", "\n")

			if (err != nil) != tt.wantErr {
				t.Errorf("RunTouch() error = %v, wantErr %v", err, tt.wantErr)
			}

			if stdoutOutput != tt.wantStdout {
				t.Errorf("RunTouch() stdout = %v, want %v", stdoutOutput, tt.wantStdout)
			}

			if stderrOutput != tt.wantStderr {
				t.Errorf("RunTouch() stderr = %v, want %v", stderrOutput, tt.wantStderr)
			}
		})
	}
}

// createTestCmd creates a test Cobra command with flags defined.
func createTestCmd(flagSetup ...func(*cobra.Command)) *cobra.Command {
	cmd := &cobra.Command{
		Use: "touch [flags] file...",
	}
	cmd.Flags().Bool("help", false, "help for touch")
	cmd.Flags().BoolP("access", "a", false, "change only the access time")
	cmd.Flags().BoolP("modification", "m", false, "change only the modification time")
	cmd.Flags().
		String("time", "", "change the specified time: access, atime, use (like -a); modify, mtime (like -m)")
	cmd.Flags().BoolP("no-create", "c", false, "do not create any files")
	cmd.Flags().
		BoolP("no-dereference", "h", false, "affect each symbolic link instead of any referenced file (unsupported on Windows)")
	cmd.Flags().Bool("f", false, "(ignored for compatibility)")
	cmd.Flags().StringP("reference", "r", "", "use this file's times instead of current time")
	cmd.Flags().StringP("stamp", "t", "", "use [[CC]YY]MMDDhhmm[.ss] instead of current time")
	cmd.Flags().StringP("date", "d", "", "parse ARG and use it instead of current time")
	cmd.Flags().BoolP("version", "v", false, "output version information and exit")

	for _, setup := range flagSetup {
		setup(cmd)
	}

	return cmd
}
