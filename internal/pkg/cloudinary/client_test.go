package cloudinary

import (
	"crypto/sha1"
	"encoding/hex"
	"strings"
	"testing"
)

func TestSignParams_MatchesCloudinaryAlgorithm(t *testing.T) {
	// Example from Cloudinary docs: alphabetical params + secret, SHA-1 hex.
	params := map[string]string{
		"folder":    "lumina/books",
		"timestamp": "1315060510",
	}
	got := signParams(params, "abcd")
	payload := "folder=lumina/books&timestamp=1315060510abcd"
	sum := sha1.Sum([]byte(payload))
	want := hex.EncodeToString(sum[:])
	if got != want {
		t.Fatalf("signature=%s want=%s", got, want)
	}
}

func TestSanitizeFilename(t *testing.T) {
	if got := sanitizeFilename("../../evil.png"); got != "evil.png" {
		t.Fatalf("got %q", got)
	}
	if got := sanitizeFilename(""); !strings.Contains(got, "upload") {
		t.Fatalf("got %q", got)
	}
}

func TestAllowedFolders(t *testing.T) {
	for _, f := range []string{"books", "students", "badges", "users"} {
		if _, ok := AllowedFolders[f]; !ok {
			t.Fatalf("missing folder %s", f)
		}
	}
}
