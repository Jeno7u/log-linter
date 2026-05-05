package b

import "go.uber.org/zap"

func F() {
	// -----------------------------------
	// RULE 1 - starts with lowercase letter
	// -----------------------------------

	zap.L().Info("Starting server on port 8080") // want "log message should start with a lowercase letter"

	zap.L().Error("Failed to connect to database") // want "log message should start with a lowercase letter"

	// ok
	zap.L().Info("starting server on port 8080")
	zap.L().Error("failed to connect to database")

	// -----------------------------------
	// RULE 2 - contains non english letters
	// -----------------------------------

	zap.L().Info("запуск сервера") // want "log message must be in English"

	zap.L().Error("ошибка подключения к базе данных") // want "log message must be in English"

	// ok
	zap.L().Info("starting server")
	zap.L().Error("failed to connect to database")

	// -----------------------------------
	// RULE 3 - contains special symbols or emoji
	// -----------------------------------

	zap.L().Info("server started 🚀") // want "log message contains special symbols or emoji"

	zap.L().Error("connection failed!!!") // want "log message contains special symbols or emoji"

	zap.L().Warn("warning: something went wrong...") // want "log message contains special symbols or emoji"

	// ok
	zap.L().Info("server started")
	zap.L().Error("connection failed")
	zap.L().Warn("something went wrong")

	// -----------------------------------
	// RULE 4 - contains potentially sensitive data
	// -----------------------------------

	password := "password"
	apiKey := "api-key"
	token := "token"

	zap.L().Info("user password:" + password) // want "log message contains potentially sensitive data"

	zap.L().Debug("api_key=" + apiKey) // want "log message contains potentially sensitive data"

	zap.L().Info("token:" + token) // want "log message contains potentially sensitive data"

	// ok
	zap.L().Info("user authenticated successfully")
	zap.L().Debug("api request completed")
	zap.L().Info("token validated")
}
