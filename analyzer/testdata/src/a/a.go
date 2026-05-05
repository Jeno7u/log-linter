package a

import "log/slog"

func F() {
	// -----------------------------------
	// RULE 1 - starts with lowercase letter
	// -----------------------------------

	// want "log message should start with a lowercase letter"
	slog.Info("Starting server on port 8080")

	// want "log message should start with a lowercase letter"
	slog.Error("Failed to connect to database")

	// ok
	slog.Info("starting server on port 8080")
	slog.Error("failed to connect to database")

	// -----------------------------------
	// RULE 2 - contains non english letters
	// -----------------------------------

	// want "log message must be in English"
	slog.Info("запуск сервера")

	// want "log message must be in English"
	slog.Error("ошибка подключения к базе данных")

	// ok
	slog.Info("starting server")
	slog.Error("failed to connect to database")

	// -----------------------------------
	// RULE 3 - contains special symbols or emoji
	// -----------------------------------

	// want "log message contains special symbols or emoji"
	slog.Info("server started 🚀")

	// want "log message contains special symbols or emoji"
	slog.Error("connection failed!!!")

	// want "log message contains special symbols or emoji"
	slog.Warn("warning: something went wrong...")

	// ok
	slog.Info("server started")
	slog.Error("connection failed")
	slog.Warn("something went wrong")

	// -----------------------------------
	// RULE 4 - contains potentially sensitive data
	// -----------------------------------

	password := "password"
	apiKey := "api-key"
	token := "token"

	// want "log message contains potentially sensitive data"
	slog.Info("user password:" + password)

	// want "log message contains potentially sensitive data"
	slog.Debug("api_key=" + apiKey)

	// want "log message contains potentially sensitive data"
	slog.Info("token:" + token)

	// ok
	slog.Info("user authenticated successfully")
	slog.Debug("api request completed")
	slog.Info("token validated")
}
