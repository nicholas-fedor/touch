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
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/nicholas-fedor/touch/internal/cli"
)

// ExitFunc is a variable for the exit function, allowing mocking in tests.
var ExitFunc = os.Exit

// rootCmd represents the base command when called without any subcommands.
// It configures the "touch" command to mimic the behavior of the GNU touch utility,
// allowing creation or timestamp updates for files with various options.
var rootCmd = &cobra.Command{
	Use:   "touch [flags] file...",
	Short: "Change file access and modification times",
	Long: `touch changes the access and/or modification times of the specified files.
If a file does not exist, it is created empty unless -c or --no-create is specified.
By default, the current time is used unless a specific time is provided via -d, -r, or -t.
Supported date formats for -d include RFC3339, YYYY-MM-DDTHH:MM:SS, YYYY-MM-DD HH:MM:SS, YYYY-MM-DDTHH:MM, YYYY-MM-DD, HH:MM:SS, HH:MM.

Examples:
  touch file.txt                  # Create or update file.txt with current time
  touch -a file.txt               # Change only access time
  touch -d "2025-07-13 14:30" file.txt  # Set specific date and time
  touch -r ref.txt file.txt       # Use times from ref.txt

For more details, see the GNU touch manual or use --help.`,
	RunE:          cli.RunTouch, // Delegate to cli.RunTouch for execution logic, allowing separation from Cobra setup.
	SilenceErrors: true,         // Prevent Cobra from printing errors automatically.
	SilenceUsage:  true,         // Prevent Cobra from printing usage on error automatically.
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
// It handles any errors by exiting with a non-zero status.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		if err.Error() == "missing operands" || err.Error() == "invalid time argument" {
			rootCmd.Usage()
		}
		ExitFunc(1)
	}
}

// SetVersionInfo sets the version information for the root command.
func SetVersionInfo(version, commit, date string) {
	rootCmd.Version = fmt.Sprintf("%s (Built on %s from Git SHA %s)", version, date, commit)
}
