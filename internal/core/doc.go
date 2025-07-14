// Package core provides the main Touch function and utilities for changing file timestamps.
// It integrates with the filesystem package for file operations and the platform package
// for OS-specific behaviors, such as retrieving access times and setting timestamps without
// dereferencing symlinks.
//
// Main Functions:
// - Touch: Applies specified timestamps to a file, creating it if necessary (unless noCreate is true).
//   Supports partial updates by preserving existing times and handles no-dereference mode.
// - Now: A variable holding the function to get the current time, allowing mocking in tests.
// - BoolToInt: Converts a boolean to an integer (1 for true, 0 for false), used for flag counting.
// - Quote: Wraps a string in quotes for safe display in error messages.
//
// Constants:
// - ChAtime, ChMtime: Bit flags to determine which timestamps to update.
//
// This package is designed to be platform-agnostic, delegating OS-specific logic to the platform package.
// It is used by the cli package to perform the actual touch operations on files.
package core
