# jig

```
     _ _
    | (_) __ _
 _  | | |/ _` |
| |_| | | (_| |
 \___/|_|\__, |
          |___/

  Apple Dev CLI
```

A lightweight CLI tool that wraps `xcodebuild` and `swiftlint` commands using a simple `jig.toml` config file. Designed for clean, token-efficient output — ideal for both humans and AI agents.

---

## Features

- **Summary output by default** — only errors, warnings, and status lines are shown
- **Full log always written** to a temp file for post-hoc debugging
- **Multi-target support** — build, test, or clean iOS and macOS targets independently
- **TOML config** — drop a `jig.toml` in your project root, works from any subdirectory
- **Agent-friendly** — minimal noise, clear exit codes, structured output

---

## Installation

**Prerequisites:** Go 1.22+, [Homebrew](https://brew.sh)

```bash
git clone https://github.com/natebird/jig.git
cd jig
make install-local   # installs to ~/.local/bin/jig
```

Or system-wide (requires sudo):

```bash
sudo make install    # installs to /usr/local/bin/jig
```

Make sure the install directory is on your `$PATH`:

```bash
# ~/.zshrc or ~/.bashrc
export PATH="$HOME/.local/bin:$PATH"
```

---

## Configuration

Create a `jig.toml` in your Xcode project root:

```toml
[project]
name = "MyApp"
path = "MyApp.xcodeproj"  # optional; omit to let xcodebuild auto-detect

[targets.ios]
scheme = "iOS"
destination = "platform=iOS Simulator,name=iPhone 17"

[targets.macos]
scheme = "macOS"
destination = "platform=macOS"

[lint]
executable = "swiftlint"  # default; omit if using swiftlint
args = []                 # extra args passed to the lint executable
```

`jig` searches the current directory and all parent directories for `jig.toml`, similar to how `git` finds `.git`.

---

## Usage

```bash
jig build                    # build all targets
jig build --target ios       # build iOS only

jig test                     # test all targets
jig test --target ios        # iOS tests only
jig test --only-testing "Tests_iOS/SearchUITests"
jig test --test-plan MyPlan

jig clean                    # clean all targets
jig clean --target ios

jig lint                     # run swiftlint
```

Add `--verbose` / `-v` to any command to stream the full output:

```bash
jig build --target ios --verbose
```

---

## Output

Default (summary) mode keeps output minimal:

```
==> Building: ios (scheme=iOS)
** BUILD SUCCEEDED **
Elapsed: 13.7s
Log: /tmp/jig-xcodebuild-20260424-143201.log
```

On failure, only the relevant error lines are shown:

```
==> Building: ios (scheme=iOS)
/path/to/QuoteView.swift:42:10: error: use of undeclared identifier 'quoteText'
** BUILD FAILED **
Elapsed: 8.1s
Log: /tmp/jig-xcodebuild-20260424-143201.log
```

The full xcodebuild output is always saved to the log file.

---

## Dependencies

- [cobra](https://github.com/spf13/cobra) — CLI framework
- [BurntSushi/toml](https://github.com/BurntSushi/toml) — TOML parsing

---

## License

MIT
