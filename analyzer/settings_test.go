package analyzer

import (
	"os"
	"testing"
)

// helper function to isolate tests from local .env
func withTempWorkdir(t *testing.T, fn func()) {
	t.Helper()

	oldwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	tmp := t.TempDir()
	if err := os.Chdir(tmp); err != nil {
		t.Fatal(err)
	}

	defer func() {
		_ = os.Chdir(oldwd)
	}()

	fn()
}

func TestLoadSettingsDefaults(t *testing.T) {
	withTempWorkdir(t, func() {
		t.Setenv(envRuleLowercaseStart, "")
		t.Setenv(envRuleEnglishOnly, "")
		t.Setenv(envRuleSpecialSymbols, "")
		t.Setenv(envRuleSensitiveData, "")
		t.Setenv(envSensitiveKeywords, "")

		settings := LoadSettings()

		if !settings.LowercaseStart || !settings.EnglishOnly || !settings.SpecialSymbols || !settings.SensitiveData {
			t.Fatal("defaults should enable all rules")
		}

		if len(settings.SensitiveKeywords) == 0 {
			t.Fatal("default sensitive keywords must not be empty")
		}
	})
}

func TestLoadSettingsFromEnv(t *testing.T) {
	withTempWorkdir(t, func() {
		t.Setenv(envRuleLowercaseStart, "false")
		t.Setenv(envRuleEnglishOnly, "true")
		t.Setenv(envRuleSpecialSymbols, "0")
		t.Setenv(envRuleSensitiveData, "false")
		t.Setenv(envSensitiveKeywords, "password token build-id")

		settings := LoadSettings()

		if settings.LowercaseStart {
			t.Fatal("lowercase_start should be disabled")
		}
		if !settings.EnglishOnly {
			t.Fatal("english_only should stay enabled")
		}
		if settings.SpecialSymbols {
			t.Fatal("special_symbols should be disabled")
		}
		if settings.SensitiveData {
			t.Fatal("sensitive_data should be disabled")
		}

		want := []string{"password", "token", "build-id"}
		if len(settings.SensitiveKeywords) != len(want) {
			t.Fatalf("keywords = %#v; want %#v", settings.SensitiveKeywords, want)
		}
		for i := range want {
			if settings.SensitiveKeywords[i] != want[i] {
				t.Fatalf("keywords = %#v; want %#v", settings.SensitiveKeywords, want)
			}
		}
	})
}
