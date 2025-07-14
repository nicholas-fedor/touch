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
// This file processes and validates command-line flags.
package cli

import (
	"bytes"
	"os"
	"runtime"
	"testing"

	"github.com/spf13/cobra"

	"github.com/nicholas-fedor/touch/internal/core"
	"github.com/nicholas-fedor/touch/internal/errors"
)

func Test_processFlags(t *testing.T) {
	tests := []struct {
		name         string
		flagSetup    func(*cobra.Command)
		wantChange   int
		wantNoCreate bool
		wantNoDeref  bool
		wantRef      string
		wantStamp    string
		wantDate     string
		wantErr      error
		wantStderr   string
	}{
		{
			name:         "default no flags",
			flagSetup:    func(cmd *cobra.Command) {},
			wantChange:   core.ChAtime | core.ChMtime,
			wantNoCreate: false,
			wantNoDeref:  false,
			wantRef:      "",
			wantStamp:    "",
			wantDate:     "",
			wantErr:      nil,
			wantStderr:   "",
		},
		{
			name: "access only",
			flagSetup: func(cmd *cobra.Command) {
				cmd.Flags().Set("access", "true")
			},
			wantChange:   core.ChAtime,
			wantNoCreate: false,
			wantNoDeref:  false,
			wantRef:      "",
			wantStamp:    "",
			wantDate:     "",
			wantErr:      nil,
			wantStderr:   "",
		},
		{
			name: "modification only",
			flagSetup: func(cmd *cobra.Command) {
				cmd.Flags().Set("modification", "true")
			},
			wantChange:   core.ChMtime,
			wantNoCreate: false,
			wantNoDeref:  false,
			wantRef:      "",
			wantStamp:    "",
			wantDate:     "",
			wantErr:      nil,
			wantStderr:   "",
		},
		{
			name: "time access",
			flagSetup: func(cmd *cobra.Command) {
				cmd.Flags().Set("time", "access")
			},
			wantChange:   core.ChAtime,
			wantNoCreate: false,
			wantNoDeref:  false,
			wantRef:      "",
			wantStamp:    "",
			wantDate:     "",
			wantErr:      nil,
			wantStderr:   "",
		},
		{
			name: "time modify",
			flagSetup: func(cmd *cobra.Command) {
				cmd.Flags().Set("time", "modify")
			},
			wantChange:   core.ChMtime,
			wantNoCreate: false,
			wantNoDeref:  false,
			wantRef:      "",
			wantStamp:    "",
			wantDate:     "",
			wantErr:      nil,
			wantStderr:   "",
		},
		{
			name: "invalid time",
			flagSetup: func(cmd *cobra.Command) {
				cmd.Flags().Set("time", "invalid")
			},
			wantChange:   0,
			wantNoCreate: false,
			wantNoDeref:  false,
			wantRef:      "",
			wantStamp:    "",
			wantDate:     "",
			wantErr:      errors.ErrInvalidTimeArg,
			wantStderr:   "",
		},
		{
			name: "no create",
			flagSetup: func(cmd *cobra.Command) {
				cmd.Flags().Set("no-create", "true")
			},
			wantChange:   core.ChAtime | core.ChMtime,
			wantNoCreate: true,
			wantNoDeref:  false,
			wantRef:      "",
			wantStamp:    "",
			wantDate:     "",
			wantErr:      nil,
			wantStderr:   "",
		},
		{
			name: "no deref non-windows",
			flagSetup: func(cmd *cobra.Command) {
				cmd.Flags().Set("no-dereference", "true")
			},
			wantChange:   core.ChAtime | core.ChMtime,
			wantNoCreate: false,
			wantNoDeref:  runtime.GOOS != "windows",
			wantRef:      "",
			wantStamp:    "",
			wantDate:     "",
			wantErr:      nil,
			wantStderr: func() string {
				if runtime.GOOS == "windows" {
					return "Warning: -h/--no-dereference is not supported on Windows; symlinks will be followed\n"
				}

				return ""
			}(),
		},
		{
			name: "reference file",
			flagSetup: func(cmd *cobra.Command) {
				cmd.Flags().Set("reference", "ref.txt")
			},
			wantChange:   core.ChAtime | core.ChMtime,
			wantNoCreate: false,
			wantNoDeref:  false,
			wantRef:      "ref.txt",
			wantStamp:    "",
			wantDate:     "",
			wantErr:      nil,
			wantStderr:   "",
		},
		{
			name: "stamp",
			flagSetup: func(cmd *cobra.Command) {
				cmd.Flags().Set("stamp", "2507131430")
			},
			wantChange:   core.ChAtime | core.ChMtime,
			wantNoCreate: false,
			wantNoDeref:  false,
			wantRef:      "",
			wantStamp:    "2507131430",
			wantDate:     "",
			wantErr:      nil,
			wantStderr:   "",
		},
		{
			name: "date",
			flagSetup: func(cmd *cobra.Command) {
				cmd.Flags().Set("date", "2025-07-13")
			},
			wantChange:   core.ChAtime | core.ChMtime,
			wantNoCreate: false,
			wantNoDeref:  false,
			wantRef:      "",
			wantStamp:    "",
			wantDate:     "2025-07-13",
			wantErr:      nil,
			wantStderr:   "",
		},
		{
			name: "multiple time sources ref and date",
			flagSetup: func(cmd *cobra.Command) {
				cmd.Flags().Set("reference", "ref.txt")
				cmd.Flags().Set("date", "2025-07-13")
			},
			wantChange:   0,
			wantNoCreate: false,
			wantNoDeref:  false,
			wantRef:      "",
			wantStamp:    "",
			wantDate:     "",
			wantErr:      errors.ErrMultipleTimeSources,
			wantStderr:   "",
		},
		{
			name: "multiple time sources stamp and date",
			flagSetup: func(cmd *cobra.Command) {
				cmd.Flags().Set("stamp", "2507131430")
				cmd.Flags().Set("date", "2025-07-13")
			},
			wantChange:   0,
			wantNoCreate: false,
			wantNoDeref:  false,
			wantRef:      "",
			wantStamp:    "",
			wantDate:     "",
			wantErr:      errors.ErrMultipleTimeSources,
			wantStderr:   "",
		},
		{
			name: "multiple time sources ref and stamp",
			flagSetup: func(cmd *cobra.Command) {
				cmd.Flags().Set("reference", "ref.txt")
				cmd.Flags().Set("stamp", "2507131430")
			},
			wantChange:   0,
			wantNoCreate: false,
			wantNoDeref:  false,
			wantRef:      "",
			wantStamp:    "",
			wantDate:     "",
			wantErr:      errors.ErrMultipleTimeSources,
			wantStderr:   "",
		},
		{
			name: "all time sources",
			flagSetup: func(cmd *cobra.Command) {
				cmd.Flags().Set("reference", "ref.txt")
				cmd.Flags().Set("stamp", "2507131430")
				cmd.Flags().Set("date", "2025-07-13")
			},
			wantChange:   0,
			wantNoCreate: false,
			wantNoDeref:  false,
			wantRef:      "",
			wantStamp:    "",
			wantDate:     "",
			wantErr:      errors.ErrMultipleTimeSources,
			wantStderr:   "",
		},
		{
			name: "ignored f flag",
			flagSetup: func(cmd *cobra.Command) {
				cmd.Flags().Set("f", "true")
			},
			wantChange:   core.ChAtime | core.ChMtime,
			wantNoCreate: false,
			wantNoDeref:  false,
			wantRef:      "",
			wantStamp:    "",
			wantDate:     "",
			wantErr:      nil,
			wantStderr:   "",
		},
		{
			name: "combined flags",
			flagSetup: func(cmd *cobra.Command) {
				cmd.Flags().Set("access", "true")
				cmd.Flags().Set("no-create", "true")
				cmd.Flags().Set("reference", "ref.txt")
			},
			wantChange:   core.ChAtime,
			wantNoCreate: true,
			wantNoDeref:  false,
			wantRef:      "ref.txt",
			wantStamp:    "",
			wantDate:     "",
			wantErr:      nil,
			wantStderr:   "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &cobra.Command{}
			// Define all flags as in init().
			cmd.Flags().Bool("help", false, "")
			cmd.Flags().BoolP("access", "a", false, "")
			cmd.Flags().BoolP("modification", "m", false, "")
			cmd.Flags().String("time", "", "")
			cmd.Flags().BoolP("no-create", "c", false, "")
			cmd.Flags().BoolP("no-dereference", "h", false, "")
			cmd.Flags().Bool("f", false, "")
			cmd.Flags().StringP("reference", "r", "", "")
			cmd.Flags().StringP("stamp", "t", "", "")
			cmd.Flags().StringP("date", "d", "", "")
			cmd.Flags().BoolP("version", "v", false, "")

			if tt.flagSetup != nil {
				tt.flagSetup(cmd)
			}

			// Capture stderr for warnings.
			oldStderr := os.Stderr
			r, w, _ := os.Pipe()
			os.Stderr = w

			got, got1, got2, got3, got4, got5, err := processFlags(cmd)

			w.Close()
			os.Stderr = oldStderr

			var buf bytes.Buffer
			buf.ReadFrom(r)
			stderrOutput := buf.String()

			if (err != nil) != (tt.wantErr != nil) {
				t.Errorf("processFlags() error = %v, wantErr %v", err, tt.wantErr)
			} else if tt.wantErr != nil && err.Error() != tt.wantErr.Error() {
				t.Errorf("processFlags() error = %v, want %v", err, tt.wantErr)
			}
			if got != tt.wantChange {
				t.Errorf("processFlags() got = %v, want %v", got, tt.wantChange)
			}
			if got1 != tt.wantNoCreate {
				t.Errorf("processFlags() got1 = %v, want %v", got1, tt.wantNoCreate)
			}
			if got2 != tt.wantNoDeref {
				t.Errorf("processFlags() got2 = %v, want %v", got2, tt.wantNoDeref)
			}
			if got3 != tt.wantRef {
				t.Errorf("processFlags() got3 = %v, want %v", got3, tt.wantRef)
			}
			if got4 != tt.wantStamp {
				t.Errorf("processFlags() got4 = %v, want %v", got4, tt.wantStamp)
			}
			if got5 != tt.wantDate {
				t.Errorf("processFlags() got5 = %v, want %v", got5, tt.wantDate)
			}
			if stderrOutput != tt.wantStderr {
				t.Errorf("processFlags() stderr = %v, want %v", stderrOutput, tt.wantStderr)
			}
		})
	}
}
