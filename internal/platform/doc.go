// Package platform provides platform-specific implementations for timestamp operations
// in the touch tool. It defines exported variables GetAtime and SetTimesNoDeref, which
// are overridden by build tags for different operating systems (Unix, Darwin, Windows).
//
// Main Components:
// - GetAtime: Function to retrieve the access time from file info, using OS-specific structures.
// - SetTimesNoDeref: Function to set timestamps without dereferencing symlinks, using OS-specific calls.
// - init: Sets fallback implementations for unsupported platforms or default behaviors.
//
// Build Tags:
// - touch_unix.go: For Unix-like systems (non-Windows, non-Darwin), uses syscall.Stat_t and unix.UtimesNanoAt.
// - touch_darwin.go: For Darwin (macOS), uses syscall.Stat_t and unix.Lutimes.
// - touch_windows.go: For Windows, uses windows.Win32FileAttributeData and a custom filetimeToTime conversion.
//
// This package is used by the core package to handle OS-specific logic in a modular way,
// allowing the core Touch function to remain platform-agnostic.
package platform
