<!-- markdownlint-disable -->
<div align="center">

# touch

<img src="/.github/assets/logo.svg" alt="touch Logo" width="150">
<!-- markdownlint-restore -->

A Go implementation of the GNU touch utility

[![Latest Version](https://img.shields.io/github/tag/nicholas-fedor/touch.svg)](https://github.com/nicholas-fedor/touch/releases)
[![CircleCI](https://dl.circleci.com/status-badge/img/gh/nicholas-fedor/touch/tree/main.svg?style=shield)](https://dl.circleci.com/status-badge/redirect/gh/nicholas-fedor/touch/tree/main)
[![Codecov](https://codecov.io/gh/nicholas-fedor/touch/branch/main/graph/badge.svg)](https://codecov.io/gh/nicholas-fedor/touch)
[![Codacy Badge](https://app.codacy.com/project/badge/Grade/ffbca83bd14d48669260bb9bb38668a8)](https://www.codacy.com/gh/nicholas-fedor/touch/dashboard?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=nicholas-fedor/touch&amp;utm_campaign=Badge_Grade)
[![GoDoc](https://godoc.org/github.com/nicholas-fedor/touch?status.svg)](https://godoc.org/github.com/nicholas-fedor/touch)
[![Go Report Card](https://goreportcard.com/badge/github.com/nicholas-fedor/touch)](https://goreportcard.com/report/github.com/nicholas-fedor/touch)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/nicholas-fedor/touch)
[![License](https://img.shields.io/badge/License-AGPL%20v3-blue.svg)](https://www.gnu.org/licenses/agpl-3.0)

</div>

`touch` is a command-line tool that changes the access and/or modification times of the specified files. If a file does not exist, it is created empty unless the `--no-create` flag is specified. It mimics the behavior of the GNU touch utility, with support for various timestamp sources like reference files, POSIX stamps, or date strings.

## Features

- Create new files or update timestamps of existing files.
- Support for GNU-compatible flags and options, including `--access`, `--modification`, `--date`, `--reference`, `--stamp`, and more.
- Modular design with separate packages for CLI handling (Cobra), core logic, filesystem interactions, timestamp parsing, and platform-specific functionality.
- Comprehensive unit tests for CLI logic, timestamp calculation, and filesystem operations.
- Cross-platform compatibility, with notes for Windows-specific limitations (e.g., `--no-dereference` is unsupported on Windows).

## Installation

Install `touch` using Go:

```bash
go install github.com/nicholas-fedor/touch@latest
```

This places the touch binary in your `$GOPATH/bin` (e.g., `~/go/bin/`).

## Usage

### Basic Usage

Create or update a file with current time:

```bash
touch file.txt
```

### Flags

| Flag                   | Description                                                                        |
|------------------------|------------------------------------------------------------------------------------|
| -a, --access           | Change only the access time.                                                       |
| -m, --modification     | Change only the modification time.                                                 |
| --time string          | Change the specified time: access, atime, use (like -a); modify, mtime (like -m).  |
| -c, --no-create        | Do not create any files.                                                           |
| -h, --no-dereference   | Affect each symbolic link instead of any referenced file (unsupported on Windows). |
| --f                    | (Ignored for compatibility with GNU touch).                                        |
| -r, --reference string | Use this file's times instead of current time.                                     |
| -t, --stamp string     | Use [[CC]YY]MMDDhhmm[.ss] instead of current time.                                 |
| -d, --date string      | Parse ARG and use it instead of current time.                                      |
| -v, --version          | Output version information and exit.                                               |
| --help                 | Show help message.                                                                 |

### Examples

- Change only access time:

```bash
touch -a file.txt
```

- Set a specific date and time:

```bash
touch -d "2025-07-13 14:30" file.txt
```

- Use times from a reference file:

```bash
touch -r ref.txt file.txt
```

- POSIX stamp format:

```bash
touch -t 2507131430 file.txt
```

- Obsolete usage (treated as POSIX stamp):

```bash
touch 2507131430 file.txt
```

For more details, run `--help` or see the GNU touch manual.

## Building from Source

Clone the repository and build:

```bash
git clone https://github.com/nicholas-fedor/touch.git
cd touch
go build -o touch ./cmd
```

Run locally:

```bash
./touch
```

## Requirements

- Go 1.22 or later (for module support).

## License

This project is licensed under the GNU Affero General Public License v3.0 — see the [LICENSE](LICENSE.md) file for details.

## Contributing

Contributions are welcome! Please submit issues or pull requests on GitHub.

Make sure to run tests before submitting:

```bash
go test ./...
```

## Logo

Special thanks to [Maria Letta](https://github.com/MariaLetta) for providing an awesome [collection](https://github.com/MariaLetta/free-gophers-pack) of Go gophers.
