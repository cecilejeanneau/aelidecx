package main

import (
	"strings"
	"testing"

	"github.com/microcosm-cc/bluemonday"
)

func TestSanitizeTitleTrimAndLimit(t *testing.T) {
	input := "   " + strings.Repeat("x", 300) + "   "
	out := sanitizeTitle(input)
	if len(out) != 255 {
		t.Fatalf("expected length 255, got %d", len(out))
	}
	if strings.HasPrefix(out, " ") || strings.HasSuffix(out, " ") {
		t.Fatalf("expected trimmed title")
	}
}

func TestValidateFolder(t *testing.T) {
	if !validateFolder("inbox") {
		t.Fatalf("expected inbox to be valid")
	}
	if validateFolder("../../etc") {
		t.Fatalf("expected traversal-like folder to be invalid")
	}
}

func TestSanitizeContentRemovesScript(t *testing.T) {
	htmlPolicy = bluemonday.UGCPolicy()
	input := `<p>hello</p><script>alert('xss')</script>`
	out := sanitizeContent(input)
	if strings.Contains(strings.ToLower(out), "<script") {
		t.Fatalf("expected script tags to be removed, got %s", out)
	}
}

func TestIsAllowedImageMime(t *testing.T) {
	if !isAllowedImageMime(".png", "image/png") {
		t.Fatalf("expected png mime to be allowed")
	}
	if isAllowedImageMime(".png", "image/jpeg") {
		t.Fatalf("expected mismatched mime to be rejected")
	}
	if isAllowedImageMime(".svg", "image/svg+xml") {
		t.Fatalf("expected svg to be rejected")
	}
}

func TestDeriveKeyDeterministic(t *testing.T) {
	salt := []byte("1234567890abcdef")
	k1 := deriveKey("passphrase-123", salt, 1000)
	k2 := deriveKey("passphrase-123", salt, 1000)
	if len(k1) != 32 {
		t.Fatalf("expected key length 32, got %d", len(k1))
	}
	if string(k1) != string(k2) {
		t.Fatalf("expected deterministic key derivation")
	}
}

func TestGetConfiguredAPIKeysPriority(t *testing.T) {
	t.Setenv("API_KEY", "single_key_abcdefghijklmnopqrstuvwxyz")
	t.Setenv("API_KEYS", "keyA_abcdefghijklmnopqrstuvwxyz,keyB_abcdefghijklmnopqrstuvwxyz")

	keys := getConfiguredAPIKeys()
	if len(keys) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(keys))
	}
	if keys[0] != "keyA_abcdefghijklmnopqrstuvwxyz" {
		t.Fatalf("expected API_KEYS to have priority")
	}
}

func TestGetConfiguredAPIKeysEmpty(t *testing.T) {
	t.Setenv("API_KEY", "")
	t.Setenv("API_KEYS", "")
	keys := getConfiguredAPIKeys()
	if len(keys) != 0 {
		t.Fatalf("expected 0 keys when env is empty, got %v", keys)
	}
}

func TestGetConfiguredAPIKeysSingleKey(t *testing.T) {
	t.Setenv("API_KEYS", "")
	t.Setenv("API_KEY", "only_key_abcdefghijklmnopqrstuvwxyz")
	keys := getConfiguredAPIKeys()
	if len(keys) != 1 || keys[0] != "only_key_abcdefghijklmnopqrstuvwxyz" {
		t.Fatalf("expected single key via API_KEY, got %v", keys)
	}
}

func TestValidateFolderAllValues(t *testing.T) {
	valid := []string{"inbox", "concepts", "practices", "ui-ux", "postmortems", "references"}
	for _, f := range valid {
		if !validateFolder(f) {
			t.Errorf("expected %q to be valid", f)
		}
	}
	invalid := []string{"", "admin", "../../etc", "INBOX", "concepts/evil", " inbox", "references "}
	for _, f := range invalid {
		if validateFolder(f) {
			t.Errorf("expected %q to be invalid", f)
		}
	}
}

func TestSanitizeTitleEdgeCases(t *testing.T) {
	if got := sanitizeTitle(""); got != "" {
		t.Fatalf("empty: expected \"\", got %q", got)
	}
	if got := sanitizeTitle("   "); got != "" {
		t.Fatalf("spaces only: expected \"\", got %q", got)
	}
	s255 := strings.Repeat("x", 255)
	if got := sanitizeTitle(s255); got != s255 {
		t.Fatalf("exactly 255 chars: expected unchanged, got len=%d", len(got))
	}
	input := "   " + strings.Repeat("a", 300)
	got := sanitizeTitle(input)
	if len(got) != 255 {
		t.Fatalf("long with spaces: expected len 255, got %d", len(got))
	}
	if strings.HasPrefix(got, " ") {
		t.Fatal("long with spaces: expected no leading spaces after sanitize")
	}
}
