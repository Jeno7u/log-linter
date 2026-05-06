# Log-Linter (loglint)
A Go analyzer for checking log message rules.

---

## Requirements

- **Go 1.22+**
- **golangci-lint**

---

## Run Tests

Run analyzer tests:

```bash
go test ./analyzer -run TestAnalyzer
```

Run all tests:

```bash
go test ./...
```

---

## Use Cases

- Check log messages in Go services before they reach production.
- Keep `log/slog` and `go.uber.org/zap` messages consistent.
- Catch log messages that start with uppercase letters.
- Catch non-English text, emoji, special symbols, and sensitive keywords.

Example:

```go
import (
  "log/slog"

  "go.uber.org/zap"
)

func main() {
  slog.Info("Starting server")      // bad
  slog.Info("starting server")      // ok
  zap.L().Error("password: 123456") // bad
  zap.L().Info("server started")    // ok
}
```

---

## Build Custom `golangci-lint`

This project is used as a `golangci-lint` module plugin.

1. **Build the custom binary** from the repository root:

```bash
golangci-lint custom
```

2. **Check the generated binary**:

```bash
./bin/loglint version
```

---

## Use in Another Project

1. **Add `.golangci.yml`** to the target project:

```yaml
version: 2
linters:
  enable:
    - loglint
  settings:
    custom:
      loglint:
        type: module
        description: Log message rules for slog and zap
```

2. **Run `loglint`** from the target project root:

```bash
path/to/loglint run ./...
```

---

## Add `loglint` to `PATH` (optional)

If you are inside the repository root, add the local `bin` folder to `PATH`:

### macOS / Linux

```bash
export PATH="$(pwd)/bin:$PATH"
```

To keep it after terminal restart, add the same command to your shell config (`~/.zshrc`, `~/.bashrc`, etc.) and use the absolute path to this repository.

### Windows PowerShell

```powershell
$env:Path = "$(Get-Location)\bin;" + $env:Path
```

To make it permanent on Windows, add the `bin` folder to your user `PATH` in Environment Variables.
