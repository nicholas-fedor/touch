// Package timestamp handles timestamp parsing and retrieval for the touch tool.
// It provides functions to parse POSIX timestamps, flexible date strings, and
// retrieve access/modification times from reference files, supporting dereference control.
//
// Main Functions:
// - ParsePosixTime: Parses POSIX timestamp format [[CC]YY]MMDDhhmm[.ss], handling century/year variations.
// - ParseDate: Parses date strings in formats like RFC3339, YYYY-MM-DDTHH:MM:SS, and time-only variants.
// - GetTimesFromRef: Retrieves access and modification times from a reference file, using Stat or Lstat based on noDeref.
//
// This package is used by the cli package to compute timestamps from user input or reference files.
// It assumes local timezone for all parsing and integrates with the filesystem and platform packages
// for file info retrieval.
package timestamp
