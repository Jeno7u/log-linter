package a
package a

import "log/slog"

func F() {
	// should warn: starts with uppercase
	// want "log message should start with a lowercase letter"
	slog.Info("Starting server on port 8080")

	// ok
	slog.Info("starting server on port 8080")

	// should warn: non-english (cyrillic)
	// want "log message must be in English"
	slog.Info("запуск сервера")

	// should warn: emoji
	// want "log message contains special symbols or emoji"
	slog.Info("server started 🚀")

	// should warn: sensitive data
	// want "log message contains potentially sensitive data"
	slog.Info("user password: secret")
}
