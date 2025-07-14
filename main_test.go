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

package main

import (
	"os"
	"testing"

	"github.com/nicholas-fedor/touch/cmd"
)

func TestMain(t *testing.T) {
	// Mock cmd.ExitFunc to capture the exit code.
	var exitCode int

	origExit := cmd.ExitFunc
	cmd.ExitFunc = func(code int) {
		exitCode = code
	}

	defer func() { cmd.ExitFunc = origExit }()

	// Set os.Args to ["touch", "--version"] to print version and exit 0 without calling RunE.
	oldArgs := os.Args
	os.Args = []string{"touch", "--version"}

	defer func() { os.Args = oldArgs }()

	// Run main.
	main()

	// Check exit code.
	if exitCode != 0 {
		t.Errorf("main() exited with code %d, want 0", exitCode)
	}
}
