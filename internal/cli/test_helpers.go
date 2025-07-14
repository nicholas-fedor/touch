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
// This file contains shared test helpers for CLI tests.
package cli

import (
	"os"

	"github.com/nicholas-fedor/touch/internal/core"
)

// mockFileInfo is a simple mock for os.FileInfo in tests.
type mockFileInfo struct {
	access core.Time
	mod    core.Time
}

func (m mockFileInfo) Name() string       { return "" }
func (m mockFileInfo) Size() int64        { return 0 }
func (m mockFileInfo) Mode() os.FileMode  { return 0 }
func (m mockFileInfo) ModTime() core.Time { return m.mod }
func (m mockFileInfo) IsDir() bool        { return false }
func (m mockFileInfo) Sys() any           { return nil }
