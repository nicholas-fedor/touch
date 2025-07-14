// Package cmd handles the command-line interface for the touch tool using the Cobra library.
// It defines the root command, sets up flags, and delegates execution to the cli package for
// processing. The package also manages version information and error handling for the CLI.
//
// Main Functions:
// - Execute: Runs the root command, handling errors by printing to stderr, displaying usage if appropriate, and exiting with a non-zero status via ExitFunc.
// - SetVersionInfo: Sets the version string for the root command, incorporating build details like commit and date.
//
// Exported Variables:
// - ExitFunc: A variable for the exit function (defaults to os.Exit), allowing mocking in tests.
//
// This package is called from main.go, where version information is set before executing the command.
// Flags are defined in flags.go and initialized in init().
package cmd
