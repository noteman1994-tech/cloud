//go:build unit

package antigravity

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"net/url"
	"strings"
	"testing"
	"time"
)

func TestGetClientSecretReturnsTrimmedDefault(t *testing.T) {
	old := defaultClientSecret
	defaultClientSecret = "  valid-secret  "
	t.Cleanup(func() { defaultClientSecret = old })

	secret, err := getClientSecret()
	if err != nil {
		t.Fatalf("getClientSecret returned error: %v", err)
	}
	if secret != "valid-secret" {
		t.Fatalf("getClientSecret = %q, want %q", secret, "valid-secret")
	}
}

func TestGetClientSecretRequiresConfiguredValue(t *testing.T) {
	old := defaultClientSecret
	defaultClientSecret = ""
	t.Cleanup(func() { defaultClientSecret = old })

	secret, err := getClientSecret()
	if err == nil {
		t.Fatal("getClientSecret returned nil error for empty secret")
	}
	if secret != "" {
		t.Fatalf("secret = %q, want empty", secret)
	}
	if !strings.Contains(err.Error(), AntigravityOAuthClientSecretEnv) {
		t.Fatalf("error %q does not mention %s", err.Error(), AntigravityOAuthClientSecretEnv)
	}
}

func TestForwardBaseURLsPrioritizesDailyWithoutMutatingBaseURLs(t *testing.T) {
	original := append([]string(nil), BaseURLs...)

	urls := ForwardBaseURLs()

	if len(urls) != len(BaseURLs) {
		t.Fatalf("len(urls) = %d, want %d", len(urls), len(BaseURLs))
	}
	if len(urls) > 0 && urls[0] != antigravityDailyBaseURL {
		t.Fatalf("urls[0] = %q, want %q", urls[0], antigravityDailyBaseURL)
	}
	for i := range BaseURLs {
		if BaseURLs[i] != original[i] {
			t.Fatalf("BaseURLs mutated at %d: got %q want %q", i, BaseURLs[i], original[i])
		}
	}
}

func TestURLAvailabilityLifecycle(t *testing.T) {
	ua := NewURLAvailability(20 * time.Millisecond)
	testURL := "https://example.com"

	if !ua.IsAvailable(testURL) {
		t.Fatal("new URL should be available")
	}

	ua.MarkUnavailable(testURL)
	if ua.IsAvailable(testURL) {
		t.Fatal("URL should be unavailable after MarkUnavailable")
	}

	time.Sleep(30 * time.Millisecond)
	if !ua.IsAvailable(testURL) {
		t.Fatal("URL should recover after TTL")
	}

	ua.MarkUnavailable(testURL)
	ua.MarkSuccess(testURL)
	if !ua.IsAvailable(testURL) {
		t.Fatal("URL should be available after MarkSuccess")
	}
	if ua.lastSuccess != testURL {
		t.Fatalf("lastSuccess = %q, want %q", ua.lastSuccess, testURL)
	}
}

func TestURLAvailabilityGetAvailableURLsWithBase(t *testing.T) {
	ua := NewURLAvailability(5 * time.Minute)
	base := []string{"https://prod.example.com", "https://daily.example.com", "https://backup.example.com"}

	ua.MarkUnavailable(base[0])
	ua.MarkSuccess(base[2])

	got := ua.GetAvailableURLsWithBase(base)
	want := []string{base[2], base[1]}

	if len(got) != len(want) {
		t.Fatalf("len(got) = %d, want %d; got=%v", len(got), len(want), got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("got[%d] = %q, want %q; got=%v", i, got[i], want[i], got)
		}
	}
}

func TestSessionStoreGetDeleteAndExpiry(t *testing.T) {
	store := NewSessionStore()
	t.Cleanup(store.Stop)

	fresh := &OAuthSession{
		State:        "fresh",
		CodeVerifier: "verifier",
		CreatedAt:    time.Now(),
	}
	store.Set("fresh", fresh)

	got, ok := store.Get("fresh")
	if !ok || got == nil {
		t.Fatal("expected fresh session to be returned")
	}
	if got.State != fresh.State {
		t.Fatalf("got state %q, want %q", got.State, fresh.State)
	}

	store.Delete("fresh")
	if _, ok := store.Get("fresh"); ok {
		t.Fatal("expected deleted session to be missing")
	}

	expired := &OAuthSession{
		State:        "expired",
		CodeVerifier: "verifier",
		CreatedAt:    time.Now().Add(-SessionTTL - time.Minute),
	}
	store.Set("expired", expired)
	if _, ok := store.Get("expired"); ok {
		t.Fatal("expected expired session to be treated as missing")
	}
}

func TestGenerateHelpersProduceExpectedFormats(t *testing.T) {
	state, err := GenerateState()
	if err != nil {
		t.Fatalf("GenerateState error: %v", err)
	}
	if len(state) != 43 {
		t.Fatalf("len(state) = %d, want 43", len(state))
	}
	if strings.ContainsAny(state, "+/=") {
		t.Fatalf("state %q is not base64url-safe", state)
	}

	sessionID, err := GenerateSessionID()
	if err != nil {
		t.Fatalf("GenerateSessionID error: %v", err)
	}
	if len(sessionID) != 32 {
		t.Fatalf("len(sessionID) = %d, want 32", len(sessionID))
	}
	if _, err := hex.DecodeString(sessionID); err != nil {
		t.Fatalf("sessionID %q is not hex: %v", sessionID, err)
	}

	verifier, err := GenerateCodeVerifier()
	if err != nil {
		t.Fatalf("GenerateCodeVerifier error: %v", err)
	}
	if len(verifier) != 43 {
		t.Fatalf("len(verifier) = %d, want 43", len(verifier))
	}
	if strings.ContainsAny(verifier, "+/=") {
		t.Fatalf("verifier %q is not base64url-safe", verifier)
	}
}

func TestGenerateCodeChallengeMatchesSHA256Base64URL(t *testing.T) {
	verifier := "verifier-123"
	sum := sha256.Sum256([]byte(verifier))
	want := strings.TrimRight(base64.URLEncoding.EncodeToString(sum[:]), "=")

	got := GenerateCodeChallenge(verifier)

	if got != want {
		t.Fatalf("GenerateCodeChallenge = %q, want %q", got, want)
	}
}

func TestBuildAuthorizationURLIncludesExpectedParameters(t *testing.T) {
	state := "state+with/special=chars"
	challenge := "challenge-value"

	authURL := BuildAuthorizationURL(state, challenge)
	if !strings.HasPrefix(authURL, AuthorizeURL+"?") {
		t.Fatalf("auth URL %q does not start with %q", authURL, AuthorizeURL+"?")
	}

	parsed, err := url.Parse(authURL)
	if err != nil {
		t.Fatalf("url.Parse error: %v", err)
	}

	params := parsed.Query()
	expected := map[string]string{
		"client_id":              ClientID,
		"redirect_uri":           RedirectURI,
		"response_type":          "code",
		"scope":                  Scopes,
		"state":                  state,
		"code_challenge":         challenge,
		"code_challenge_method":  "S256",
		"access_type":            "offline",
		"prompt":                 "consent",
		"include_granted_scopes": "true",
	}

	if len(params) != len(expected) {
		t.Fatalf("len(params) = %d, want %d", len(params), len(expected))
	}
	for key, want := range expected {
		if got := params.Get(key); got != want {
			t.Fatalf("param %s = %q, want %q", key, got, want)
		}
	}
}

func TestConstantsAndScopes(t *testing.T) {
	if AuthorizeURL != "https://accounts.google.com/o/oauth2/v2/auth" {
		t.Fatalf("AuthorizeURL = %q", AuthorizeURL)
	}
	if TokenURL != "https://oauth2.googleapis.com/token" {
		t.Fatalf("TokenURL = %q", TokenURL)
	}
	if UserInfoURL != "https://www.googleapis.com/oauth2/v2/userinfo" {
		t.Fatalf("UserInfoURL = %q", UserInfoURL)
	}
	if ClientID != "ANTIGRAVITY_OAUTH_CLIENT_ID" {
		t.Fatalf("ClientID = %q", ClientID)
	}
	if GetUserAgent() != "antigravity/1.21.9 windows/amd64" {
		t.Fatalf("GetUserAgent() = %q", GetUserAgent())
	}
	if SessionTTL != 30*time.Minute {
		t.Fatalf("SessionTTL = %v", SessionTTL)
	}
	if URLAvailabilityTTL != 5*time.Minute {
		t.Fatalf("URLAvailabilityTTL = %v", URLAvailabilityTTL)
	}

	requiredScopes := []string{
		"https://www.googleapis.com/auth/cloud-platform",
		"https://www.googleapis.com/auth/userinfo.email",
		"https://www.googleapis.com/auth/userinfo.profile",
		"https://www.googleapis.com/auth/cclog",
		"https://www.googleapis.com/auth/experimentsandconfigs",
	}
	for _, scope := range requiredScopes {
		if !strings.Contains(Scopes, scope) {
			t.Fatalf("Scopes does not contain %q", scope)
		}
	}
}
