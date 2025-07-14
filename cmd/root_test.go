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

// Package cmd handles the command-line interface for the touch tool using the Cobra library.
// It defines the root command and delegates execution logic to separate files for modularity.
package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/mock"

	"github.com/nicholas-fedor/touch/internal/cli"
	"github.com/nicholas-fedor/touch/internal/core"
	"github.com/nicholas-fedor/touch/internal/filesystem"
	"github.com/nicholas-fedor/touch/internal/filesystem/mocks"
	"github.com/nicholas-fedor/touch/internal/version"
)

const usageStr = "Usage:\n  touch [flags] file...\n\nFlags:\n  -a, --access             change only the access time\n  -d, --date string        parse ARG and use it instead of current time\n      --f                  (ignored for compatibility)\n      --help               help for touch\n  -m, --modification       change only the modification time\n  -c, --no-create          do not create any files\n  -h, --no-dereference     affect each symbolic link instead of any referenced file (unsupported on Windows)\n  -r, --reference string   use this file's times instead of current time\n  -t, --stamp string       use [[CC]YY]MMDDhhmm[.ss] instead of current time\n      --time string        change the specified time: access, atime, use (like -a); modify, mtime (like -m)\n  -v, --version            output version information and exit\n"

func TestRootCmd(t *testing.T) {
	if rootCmd.Use != "touch [flags] file..." {
		t.Errorf("rootCmd.Use = %q, want %q", rootCmd.Use, "touch [flags] file...")
	}

	SetVersionInfo("v0.0.1", "abc123", "2025-05-07T00:00:00Z")

	expectedVersion := "v0.0.1 (Built on 2025-05-07T00:00:00Z from Git SHA abc123)"
	if rootCmd.Version != expectedVersion {
		t.Errorf("rootCmd.Version = %q, want %q", rootCmd.Version, expectedVersion)
	}

	if rootCmd.Short == "" || rootCmd.Long == "" {
		t.Errorf("rootCmd Short or Long description is empty")
	}
}

func TestSetVersionInfo(t *testing.T) {
	SetVersionInfo("v1.0.0", "abcdef", "2025-07-13T14:30:00Z")

	expected := "v1.0.0 (Built on 2025-07-13T14:30:00Z from Git SHA abcdef)"
	if rootCmd.Version != expected {
		t.Errorf("rootCmd.Version = %q, want %q", rootCmd.Version, expected)
	}
}

func TestExecute(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		mockFSSetup func(*mocks.MockFS)
		mockRunE    func(*cobra.Command, []string) error
		wantExit    int
		wantStderr  string
	}{
		{
			name:        "success no args show usage exit 1",
			args:        []string{},
			mockFSSetup: nil,
			mockRunE:    func(_ *cobra.Command, _ []string) error { return errors.New("missing operands") },
			wantExit:    1,
			wantStderr:  "Error: missing operands\n" + usageStr,
		},
		{
			name: "success with file",
			args: []string{"testfile.txt"},
			mockFSSetup: func(m *mocks.MockFS) {
				m.On("Stat", "testfile.txt").Return(&mockFileInfo{mod: time.Now()}, nil)
				m.On("Chtimes", "testfile.txt", mock.AnythingOfType("time.Time"), mock.AnythingOfType("time.Time")).
					Return(nil)
			},
			mockRunE:   cli.RunTouch,
			wantExit:   0,
			wantStderr: "",
		},
		{
			name:        "error invalid flag exit 1",
			args:        []string{"--time", "invalid"},
			mockFSSetup: nil,
			mockRunE:    func(_ *cobra.Command, _ []string) error { return errors.New("invalid time argument") },
			wantExit:    1,
			wantStderr:  "Error: invalid time argument\n" + usageStr,
		},
		{
			name: "error apply files exit 1",
			args: []string{"errorfile.txt"},
			mockFSSetup: func(m *mocks.MockFS) {
				m.On("Stat", "errorfile.txt").Return(nil, os.ErrPermission)
			},
			mockRunE:   cli.RunTouch,
			wantExit:   1,
			wantStderr: "touch: \"errorfile.txt\": stat file errorfile.txt: permission denied\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockFS := mocks.NewMockFS(t)
			if tt.mockFSSetup != nil {
				tt.mockFSSetup(mockFS)
			}

			filesystem.Default = mockFS // Override default FS with mock.

			// Create a new command instance for each test to avoid flag conflicts.
			cmd := &cobra.Command{
				Use:     "touch [flags] file...",
				Version: "1.0.0",
				RunE:    tt.mockRunE,
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
			cmd.Flags().
				StringP("reference", "r", "", "use this file's times instead of current time")
			cmd.Flags().
				StringP("stamp", "t", "", "use [[CC]YY]MMDDhhmm[.ss] instead of current time")
			cmd.Flags().StringP("date", "d", "", "parse ARG and use it instead of current time")
			cmd.Flags().BoolP("version", "v", false, "output version information and exit")

			cmd.SilenceErrors = true
			cmd.SilenceUsage = true

			// Set up arguments and flags.
			oldArgs := os.Args

			defer func() { os.Args = oldArgs }()

			os.Args = append([]string{"touch"}, tt.args...)
			cmd.SetArgs(tt.args)

			// Capture stderr.
			oldStderr := os.Stderr

			defer func() { os.Stderr = oldStderr }()

			r, w, _ := os.Pipe()
			os.Stderr = w

			// Capture exit code.
			exitCode := 0
			origExit := ExitFunc
			ExitFunc = func(code int) {
				exitCode = code
			}

			defer func() { ExitFunc = origExit }()

			// Run Execute with the mocked command.
			if err := cmd.Execute(); err != nil {
				switch err.Error() {
				case "missing operands", "invalid time argument":
					fmt.Fprintln(os.Stderr, "Error:", err)

					if usageErr := cmd.Usage(); usageErr != nil {
						fmt.Fprintln(os.Stderr, "Error displaying usage:", usageErr)
					}
				case "errors occurred while processing files":
					// Do nothing, specific errors already printed
				default:
					if strings.HasPrefix(err.Error(), "touch: ") {
						fmt.Fprintln(os.Stderr, err)
					} else {
						fmt.Fprintln(os.Stderr, "Error:", err)
					}
				}

				ExitFunc(1)
			}

			w.Close()

			var buf bytes.Buffer
			buf.ReadFrom(r)
			stderrOutput := strings.ReplaceAll(buf.String(), "\r\n", "\n")

			if exitCode != tt.wantExit {
				t.Errorf("Execute() exit code = %v, want %v", exitCode, tt.wantExit)
			}

			if stderrOutput != tt.wantStderr {
				t.Errorf("Execute() stderr = %v, want %v", stderrOutput, tt.wantStderr)
			}
		})
	}
}

func TestVersionFromBuildInfo(t *testing.T) {
	info := version.GetVersionInfo()

	if info.Version == "unknown" {
		t.Logf("No build info found; ensure tests run in a module with version tags")
	} else if !strings.HasPrefix(info.Version, "v") {
		t.Errorf("Version = %q, want prefix 'v'", info.Version)
	}
}

// mockFileInfo is a simple mock for os.FileInfo in tests.
type mockFileInfo struct {
	mod core.Time
}

func (m mockFileInfo) Name() string       { return "" }
func (m mockFileInfo) Size() int64        { return 0 }
func (m mockFileInfo) Mode() os.FileMode  { return 0 }
func (m mockFileInfo) ModTime() core.Time { return m.mod }
func (m mockFileInfo) IsDir() bool        { return false }
func (m mockFileInfo) Sys() any           { return nil }
