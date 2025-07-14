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

package cmd

// init initializes the root command by defining all supported flags.
// Flags are bound using Cobra's flag definitions, mirroring GNU touch options.
func init() {
	// Add help flag without shorthand to avoid conflict with custom -h.
	rootCmd.Flags().Bool("help", false, "help for touch")

	// Define flags for changing specific timestamps.
	rootCmd.Flags().BoolP("access", "a", false, "change only the access time")
	rootCmd.Flags().BoolP("modification", "m", false, "change only the modification time")
	rootCmd.Flags().
		String("time", "", "change the specified time: access, atime, use (like -a); modify, mtime (like -m)")

	// Flags for controlling file creation.
	rootCmd.Flags().BoolP("no-create", "c", false, "do not create any files")

	// Flags for symlink handling.
	rootCmd.Flags().
		BoolP("no-dereference", "h", false, "affect each symbolic link instead of any referenced file (unsupported on Windows)")

	// Ignored flag for compatibility.
	rootCmd.Flags().Bool("f", false, "(ignored for compatibility)")

	// Flags for specifying reference file or timestamps.
	rootCmd.Flags().StringP("reference", "r", "", "use this file's times instead of current time")
	rootCmd.Flags().StringP("stamp", "t", "", "use [[CC]YY]MMDDhhmm[.ss] instead of current time")
	rootCmd.Flags().StringP("date", "d", "", "parse ARG and use it instead of current time")

	// Enable version flag with shorthand.
	rootCmd.Flags().BoolP("version", "v", false, "output version information and exit")
}
