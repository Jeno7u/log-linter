package analyzer

import "testing"

func TestFirstLetterIsLowercase(t *testing.T) {
	cases := []struct {
		in   string
		want bool
	}{
		{"starting server", true},
		{"Starting server", false},
		{"123 no letters", false},
		{"", false},
		{"привет", true},
	}

	for _, c := range cases {
		got := startsWithLowercaseLetter(c.in)
		if got != c.want {
			t.Fatalf("firstLetterIsLowercase(%q) = %v; want %v", c.in, got, c.want)
		}
	}
}

func TestIsStringContainsNonEnglish(t *testing.T) {
	cases := []struct {
		in   string
		want bool
	}{
		{"starting server", false},
		{"запуск сервера", true},
		{"mixed English и русский", true},
		{"12345", false},
	}
	for _, c := range cases {
		got := isStringContainsNonEnglish(c.in)
		if got != c.want {
			t.Fatalf("isStringContainsNonEnglish(%q) = %v; want %v", c.in, got, c.want)
		}
	}
}

func TestContainsSpecialSymbolsOrEmojis(t *testing.T) {
	cases := []struct {
		in   string
		want bool
	}{
		{"server started", false},
		{"server started 🚀", true},
		{"connection failed!!!", true},
		{"ok: value=42", true},
	}
	for _, c := range cases {
		got := containsSpecialSymbolsOrEmojis(c.in)
		if got != c.want {
			t.Fatalf("containsSpecialSymbolsOrEmojis(%q) = %v; want %v", c.in, got, c.want)
		}
	}
}

func TestContainsSensitive(t *testing.T) {
	keywords := []string{"password", "token", "api_key", "apikey", "secret"}
	cases := []struct {
		in   string
		want bool
	}{
		{"user password: secret", true},
		{"api_key=abcd", true},
		{"token validated", false},
		{"user authenticated successfully", false},
	}
	for _, c := range cases {
		got := containsSensitive(c.in, keywords)
		if got != c.want {
			t.Fatalf("containsSensitive(%q) = %v; want %v", c.in, got, c.want)
		}
	}
}

func TestContainsSensitiveWithCustomKeywords(t *testing.T) {
	got := containsSensitive("build-id: 123", []string{"build-id"})
	if !got {
		t.Fatalf("containsSensitive with custom keywords = %v; want true", got)
	}
}
