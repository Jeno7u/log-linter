package analyzer

import (
	"strings"
	"unicode"
)

// checks that string is started with lowercase english letter. if empty returns false
func startsWithLowercaseLetter(s string) bool {
	for _, r := range s {
		return unicode.IsLetter(r) && unicode.IsLower(r)
	}

	return false
}

// returns true if the string contains letters that are not english
func isStringContainsNonEnglish(s string) bool {
	for _, r := range s {
		if unicode.IsLetter(r) && !isEnglishLetter(r) {
			return true
		}
	}

	return false
}

// returns true if string contains symbols distinct from numbers
func containsSpecialSymbolsOrEmojis(s string) bool {
	for _, r := range s {
		if r == ' ' {
			continue
		}

		if !unicode.IsLetter(r) && !isNumber(r) {
			return true
		}
	}

	return false
}

// containsSensitive returns true if any sensitive keyword appears in the message (case-insensitive).
func containsSensitive(s string) bool {
	low := strings.ToLower(s)
	for _, k := range sensitiveKeywords {
		// flag only when keyword is followed by ':' or '=' to avoid false positives
		if strings.Contains(low, k+":") || strings.Contains(low, k+"=") {
			return true
		}
	}
	return false
}

func literalsFirstLetterIsLowercase(lits []string) bool {
	for _, s := range lits {
		for _, r := range s {
			if unicode.IsLetter(r) {
				return unicode.IsLower(r)
			}
		}
	}
	return true
}

func literalsContainNonEnglish(lits []string) bool {
	for _, s := range lits {
		if isStringContainsNonEnglish(s) {
			return true
		}
	}
	return false
}

func literalsContainSpecial(lits []string) bool {
	for _, s := range lits {
		if containsSpecialSymbolsOrEmojis(s) {
			return true
		}
	}
	return false
}

func literalsContainSensitive(lits []string) bool {
	for _, s := range lits {
		if containsSensitive(s) {
			return true
		}
	}
	return false
}
