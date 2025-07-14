// Package cli handles CLI-specific logic for the touch command-line tool.
// It processes flags, calculates timestamps based on user input, and applies
// timestamp changes to files using concurrent goroutines for efficiency.
// The package separates CLI handling from core touch functionality to improve
// modularity and testability.
//
// Main Functions:
// - RunTouch: Orchestrates the entire touch operation, serving as the entry point for Cobra's RunE.
// - processFlags: Retrieves and validates command-line flags, computing the changeTimes mask.
// - calculateTimestamps: Determines access and modification times from flags or defaults to current time.
// - applyToFiles: Applies timestamp changes concurrently to the list of files.
//
// This package integrates with the core package for the actual timestamp application
// and uses the filesystem package for file operations. It also handles platform-specific
// behaviors, such as warnings for unsupported flags on Windows.
//
// For usage examples, see the RunTouch function and associated tests.
package cli
