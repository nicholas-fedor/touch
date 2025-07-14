package version

import (
	"runtime/debug"
	"strings"
	"testing"
	"time"
)

func TestGetVersionInfo(t *testing.T) {
	tests := []struct {
		name         string
		setVars      func()
		expect       Info
		partialMatch bool
	}{
		{
			name: "GoReleaser build",
			setVars: func() {
				Version = "0.0.1"
				Commit = "abc123"
				Date = "2025-05-07T00:00:00Z"
			},
			expect: Info{
				Version: "v0.0.1",
				Commit:  "abc123",
				Date:    "2025-05-07T00:00:00Z",
			},
		},
		{
			name: "Source build with default values",
			setVars: func() {
				Version = devVersion
				Commit = unknownValue
				Date = unknownValue
			},
			expect: Info{
				Version: unknownValue,
				Commit:  unknownValue,
				Date:    unknownValue,
			},
			partialMatch: true,
		},
		{
			name: "Source build with empty values",
			setVars: func() {
				Version = ""
				Commit = ""
				Date = ""
			},
			expect: Info{
				Version: unknownValue,
				Commit:  unknownValue,
				Date:    unknownValue,
			},
		},
		{
			name: "Invalid GoReleaser version",
			setVars: func() {
				Version = "v"
				Commit = ""
				Date = ""
			},
			expect: Info{
				Version: unknownValue,
				Commit:  unknownValue,
				Date:    unknownValue,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setVars()

			info := GetVersionInfo()

			if !tt.partialMatch {
				if info.Version != tt.expect.Version {
					t.Errorf("Version = %q, want %q", info.Version, tt.expect.Version)
				}

				if info.Commit != tt.expect.Commit {
					t.Errorf("Commit = %q, want %q", info.Commit, tt.expect.Commit)
				}

				if info.Date != tt.expect.Date {
					t.Errorf("Date = %q, want %q", info.Date, tt.expect.Date)
				}
			} else if info.Version != tt.expect.Version && !strings.Contains(info.Version, "+dirty") {
				t.Errorf("Version = %q, want %q or dirty variant", info.Version, tt.expect.Version)
			}
		})
	}
}

func TestGetVersionInfo_VCSData(t *testing.T) {
	Version = devVersion
	Commit = unknownValue
	Date = unknownValue

	info := GetVersionInfo()

	if buildInfo, ok := debug.ReadBuildInfo(); ok {
		var vcsRevision, vcsTime, vcsModified string

		for _, setting := range buildInfo.Settings {
			switch setting.Key {
			case "vcs.revision":
				vcsRevision = setting.Value
			case "vcs.time":
				vcsTime = setting.Value
			case "vcs.modified":
				vcsModified = setting.Value
			}
		}

		if vcsRevision != "" {
			if info.Commit == unknownValue {
				t.Errorf(
					"Expected commit %q, got %q; ensure repository has commit history",
					vcsRevision,
					info.Commit,
				)
			} else if info.Commit != vcsRevision {
				t.Errorf("Commit = %q, want %q", info.Commit, vcsRevision)
			}
		} else {
			t.Logf("No vcs.revision found; ensure repository has Git metadata to cover commit assignment")
		}

		if vcsTime != "" {
			if _, err := time.Parse(time.RFC3339, vcsTime); err == nil {
				if info.Date == unknownValue {
					t.Errorf(
						"Expected valid date, got %q; ensure vcs.time is a valid RFC3339 timestamp",
						info.Date,
					)
				}
			} else if info.Date != unknownValue {
				t.Logf("vcs.time %q is invalid; date should remain %q", vcsTime, unknownValue)
				t.Errorf("Expected date %q, got %q for invalid vcs.time", unknownValue, info.Date)
			}
		} else {
			t.Logf("No vcs.time found; ensure repository has commit timestamps to cover date assignment")
		}

		if vcsModified == trueValue && info.Version != unknownValue {
			if !strings.Contains(info.Version, "+dirty") {
				t.Errorf(
					"Expected version to contain '+dirty', got %q; ensure repository has uncommitted changes",
					info.Version,
				)
			}
		} else if vcsModified != trueValue {
			t.Logf("Repository is clean (vcs.modified=%q); make uncommitted changes to cover '+dirty' case", vcsModified)
		}
	} else {
		t.Logf("debug.ReadBuildInfo() failed; ensure tests run in a Git repository to cover VCS parsing")
	}
}

func TestGetVersionInfo_InvalidVCSTime(t *testing.T) {
	Version = devVersion
	Commit = unknownValue
	Date = unknownValue

	info := GetVersionInfo()

	if info.Date == "" || info.Date != unknownValue {
		t.Errorf("Expected date to be %q, got %q", unknownValue, info.Date)
	}
}

func TestContains(t *testing.T) {
	tests := []struct {
		name     string
		s        string
		substr   string
		expected bool
	}{
		{
			name:     "Substring found",
			s:        "v1.0.0+dirty",
			substr:   "+dirty",
			expected: true,
		},
		{
			name:     "Substring not found",
			s:        "v1.0.0",
			substr:   "+dirty",
			expected: false,
		},
		{
			name:     "Empty string",
			s:        "",
			substr:   "+dirty",
			expected: false,
		},
		{
			name:     "Empty substring",
			s:        "v1.0.0",
			substr:   "",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := contains(tt.s, tt.substr)
			if result != tt.expected {
				t.Errorf("contains(%q, %q) = %v, want %v", tt.s, tt.substr, result, tt.expected)
			}
		})
	}
}
