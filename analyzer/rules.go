package analyzer

import (
	"strings"
	"unicode"
)

var sensitiveKeywords = []string{
	"password", "pwd", "token", "api_key", "apikey", "secret", "ssn", "creditcard", "card_number",
}

// return true if string is non empty
func isStringNonEmpty(s string) bool {
	return len(s) == 0
}

// checks that first letter in string is lowercase english. if empty returns false
func firstLetterIsLowercase(s string) bool {
	for _, r := range s {
		return unicode.IsLetter(r) && unicode.IsLower(r)
	}

	// if string is empty
	return true
}

// returns true if the string contains letters that are lowercase english
func isStringContainsNonEnglish(s string) bool {
	for _, r := range s {
		if unicode.IsLetter(r) && !((r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z')) {
			return false
		}
	}

	return true
}

// returns true if string contains symbols distinct from numbers
func containsSpecialSymbolsOrEmojis(s string) bool {
	for _, r := range s {
		if !unicode.IsLetter(r) && !(r >= '0' && r <= '9') {
			return false
		}
	}

	return true
}

// containsSensitive returns true if any sensitive keyword appears in the message (case-insensitive).
func containsSensitive(s string) bool {
	low := strings.ToLower(s)
	for _, k := range sensitiveKeywords {
		if strings.Contains(low, k) {
			return true
		}
	}
	return false
}
