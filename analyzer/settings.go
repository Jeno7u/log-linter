package analyzer

import (
	"os"
	"strconv"
	"strings"

	"github.com/lpernett/godotenv"
)

const (
	envRuleLowercaseStart = "LOGLINT_RULE_LOWERCASE_START"
	envRuleEnglishOnly    = "LOGLINT_RULE_ENGLISH_ONLY"
	envRuleSpecialSymbols = "LOGLINT_RULE_SPECIAL_SYMBOLS"
	envRuleSensitiveData  = "LOGLINT_RULE_SENSITIVE_DATA"
	envSensitiveKeywords  = "LOGLINT_SENSITIVE_KEYWORDS"
)

var defaultSensitiveKeywords = []string{
	"password",
	"token",
	"api_key",
	"apikey",
	"secret",
}

type Settings struct {
	LowercaseStart bool
	EnglishOnly    bool
	SpecialSymbols bool
	SensitiveData  bool

	SensitiveKeywords []string
}

// load settings from .env
func LoadSettings() Settings {
	loadEnvFiles()

	settings := Settings{
		LowercaseStart:    getEnvAsBool(envRuleLowercaseStart, true),
		EnglishOnly:       getEnvAsBool(envRuleEnglishOnly, true),
		SpecialSymbols:    getEnvAsBool(envRuleSpecialSymbols, true),
		SensitiveData:     getEnvAsBool(envRuleSensitiveData, true),
		SensitiveKeywords: getEnvAsStrings(envRuleEnglishOnly, defaultSensitiveKeywords),
	}

	return settings
}

// loads .env variables (uses path .env and ../.env)
func loadEnvFiles() {
	for _, file := range []string{".env", "../.env"} {
		_ = godotenv.Load(file)
	}
}

// gets value by key from env. if not found returns fallback
func getEnvAsBool(key string, fallback bool) bool {
	value, ok := os.LookupEnv(key)
	if !ok || strings.TrimSpace(value) == "" {
		return fallback
	}

	parsed, err := strconv.ParseBool(value)
	if err != nil {
		return fallback
	}

	return parsed
}

// gets parsed array of strings by key from env. if not found returns fallback
func getEnvAsStrings(key string, fallback []string) []string {
	value, ok := os.LookupEnv(key)
	if !ok || strings.TrimSpace(value) == "" {
		return append([]string(nil), fallback...)
	}

	items := strings.Fields(value)
	if len(items) == 0 {
		return append([]string(nil), fallback...)
	}

	return items
}
