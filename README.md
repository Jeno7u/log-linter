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

## Use in Projects

1. **Build `loglint`** from this repository root:

```bash
golangci-lint custom
```

2. **Add `.golangci.yml`** to the target project:

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

3. **Add/Copy `.env`** to the target project root (optional):

The generated `loglint` binary reads `.env` file at runtime. 

.env file file should be in the same folder or in parent folder.

```bash
LOGLINT_RULE_LOWERCASE_START=true
LOGLINT_RULE_ENGLISH_ONLY=true
LOGLINT_RULE_SPECIAL_SYMBOLS=true
LOGLINT_RULE_SENSITIVE_DATA=true
LOGLINT_SENSITIVE_KEYWORDS=password token api_key apikey secret
```

All rules are enabled by default. Set any rule to `false` if you want to disable it. Sensitive keywords are split by spaces.

4. **Run `loglint`** from the target project root:

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

---

## CI/CD

- **GitHub Actions**: `.github/workflows/ci.yml`

It is doing next things:

1. run `go test ./...`
2. build the custom `loglint` binary with `golangci-lint custom`
3. Check that `loglint` binary present

---

## Notes

1. Could be analyzer folder files broken into multiple folders but there is no much files, so it should be ok

2. 