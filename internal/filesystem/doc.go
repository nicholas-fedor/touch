// Package filesystem defines the FS interface for abstracting file system operations,
// allowing for testability and modularity in file interactions. It provides a default
// implementation using the os package and supports operations like retrieving file info
// (Stat/Lstat), creating files, and changing timestamps (Chtimes).
//
// Main Components:
// - FS: Interface for file system operations, including Stat, Lstat, Create, and Chtimes.
// - Default: The default FS implementation using standard os functions.
//
// This package is used by the core package to perform file operations in a way that
// can be mocked during testing. It wraps os functions with error formatting for consistency.
//
// For usage examples, see the core.Touch function and associated tests.
package filesystem
