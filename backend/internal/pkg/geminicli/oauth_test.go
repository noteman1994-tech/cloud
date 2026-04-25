package geminicli

import (
	"encoding/hex"
	"strings"
	"sync"
	"testing"
	"time"
)

// ---------------------------------------------------------------------------
// SessionStore 娴嬭瘯
// ---------------------------------------------------------------------------

func TestSessionStore_SetAndGet(t *testing.T) {
	store := NewSessionStore()
	defer store.Stop()

	session := &OAuthSession{
		State:     "test-state",
		OAuthType: "code_assist",
		CreatedAt: time.Now(),
	}
	store.Set("sid-1", session)

	got, ok := store.Get("sid-1")
	if !ok {
		t.Fatal("鏈熸湜 Get 杩斿洖 ok=true锛屽疄闄呰繑鍥?false")
	}
	if got.State != "test-state" {
		t.Errorf("鏈熸湜 State=%q锛屽疄闄?%q", "test-state", got.State)
	}
}

func TestSessionStore_GetNotFound(t *testing.T) {
	store := NewSessionStore()
	defer store.Stop()

	_, ok := store.Get("涓嶅瓨鍦ㄧ殑ID")
	if ok {
		t.Error("鏈熸湜涓嶅瓨鍦ㄧ殑 sessionID 杩斿洖 ok=false")
	}
}

func TestSessionStore_GetExpired(t *testing.T) {
	store := NewSessionStore()
	defer store.Stop()

	// Expired session should be treated as missing.
	session := &OAuthSession{
		State:     "expired-state",
		OAuthType: "code_assist",
		CreatedAt: time.Now().Add(-(SessionTTL + 1*time.Minute)),
	}
	store.Set("expired-sid", session)

	_, ok := store.Get("expired-sid")
	if ok {
		t.Error("鏈熸湜杩囨湡鐨?session 杩斿洖 ok=false")
	}
}

func TestSessionStore_Delete(t *testing.T) {
	store := NewSessionStore()
	defer store.Stop()

	session := &OAuthSession{
		State:     "to-delete",
		OAuthType: "code_assist",
		CreatedAt: time.Now(),
	}
	store.Set("del-sid", session)

	// Ensure the session exists before delete.
	if _, ok := store.Get("del-sid"); !ok {
		t.Fatal("鍒犻櫎鍓?session 搴旇瀛樺湪")
	}

	store.Delete("del-sid")

	if _, ok := store.Get("del-sid"); ok {
		t.Error("deleted session should not exist")
	}
}

func TestSessionStore_Stop_Idempotent(t *testing.T) {
	store := NewSessionStore()

	// 澶氭璋冪敤 Stop 涓嶅簲 panic
	store.Stop()
	store.Stop()
	store.Stop()
}

func TestSessionStore_ConcurrentAccess(t *testing.T) {
	store := NewSessionStore()
	defer store.Stop()

	const goroutines = 50
	var wg sync.WaitGroup
	wg.Add(goroutines * 3)

	// 骞跺彂鍐欏叆
	for i := 0; i < goroutines; i++ {
		go func(idx int) {
			defer wg.Done()
			sid := "concurrent-" + string(rune('A'+idx%26))
			store.Set(sid, &OAuthSession{
				State:     sid,
				OAuthType: "code_assist",
				CreatedAt: time.Now(),
			})
		}(i)
	}

	// 骞跺彂璇诲彇
	for i := 0; i < goroutines; i++ {
		go func(idx int) {
			defer wg.Done()
			sid := "concurrent-" + string(rune('A'+idx%26))
			store.Get(sid) // 鍙兘鎵惧埌涔熷彲鑳芥病鎵惧埌锛屽叧閿槸涓?panic
		}(i)
	}

	// 骞跺彂鍒犻櫎
	for i := 0; i < goroutines; i++ {
		go func(idx int) {
			defer wg.Done()
			sid := "concurrent-" + string(rune('A'+idx%26))
			store.Delete(sid)
		}(i)
	}

	wg.Wait()
}

// ---------------------------------------------------------------------------
// GenerateRandomBytes 娴嬭瘯
// ---------------------------------------------------------------------------

func TestGenerateRandomBytes(t *testing.T) {
	tests := []int{0, 1, 16, 32, 64}
	for _, n := range tests {
		b, err := GenerateRandomBytes(n)
		if err != nil {
			t.Errorf("GenerateRandomBytes(%d) 鍑洪敊: %v", n, err)
			continue
		}
		if len(b) != n {
			t.Errorf("GenerateRandomBytes(%d) 杩斿洖闀垮害=%d锛屾湡鏈?%d", n, len(b), n)
		}
	}
}

func TestGenerateRandomBytes_Uniqueness(t *testing.T) {
	// Two calls should almost certainly produce different results.
	a, _ := GenerateRandomBytes(32)
	b, _ := GenerateRandomBytes(32)
	if string(a) == string(b) {
		t.Error("涓ゆ GenerateRandomBytes(32) 杩斿洖浜嗙浉鍚岀粨鏋滐紝闅忔満鎬у彲鑳芥湁闂")
	}
}

// ---------------------------------------------------------------------------
// GenerateState 娴嬭瘯
// ---------------------------------------------------------------------------

func TestGenerateState(t *testing.T) {
	state, err := GenerateState()
	if err != nil {
		t.Fatalf("GenerateState() 鍑洪敊: %v", err)
	}
	if state == "" {
		t.Error("GenerateState() 杩斿洖绌哄瓧绗︿覆")
	}
	// base64url 缂栫爜涓嶅簲鍖呭惈 padding '='
	if strings.Contains(state, "=") {
		t.Errorf("GenerateState() 缁撴灉鍖呭惈 '=' padding: %s", state)
	}
	// base64url 涓嶅簲鍖呭惈 '+' 鎴?'/'
	if strings.ContainsAny(state, "+/") {
		t.Errorf("GenerateState() 缁撴灉鍖呭惈闈?base64url 瀛楃: %s", state)
	}
}

// ---------------------------------------------------------------------------
// GenerateSessionID 娴嬭瘯
// ---------------------------------------------------------------------------

func TestGenerateSessionID(t *testing.T) {
	sid, err := GenerateSessionID()
	if err != nil {
		t.Fatalf("GenerateSessionID() 鍑洪敊: %v", err)
	}
	// 16 瀛楄妭 -> 32 涓?hex 瀛楃
	if len(sid) != 32 {
		t.Errorf("GenerateSessionID() 闀垮害=%d锛屾湡鏈?32", len(sid))
	}
	// Session ID should be valid hex.
	if _, err := hex.DecodeString(sid); err != nil {
		t.Errorf("GenerateSessionID() 涓嶆槸鍚堟硶鐨?hex 瀛楃涓? %s, err=%v", sid, err)
	}
}

func TestGenerateSessionID_Uniqueness(t *testing.T) {
	a, _ := GenerateSessionID()
	b, _ := GenerateSessionID()
	if a == b {
		t.Error("GenerateSessionID() returned the same value twice")
	}
}

// ---------------------------------------------------------------------------
// GenerateCodeVerifier 娴嬭瘯
// ---------------------------------------------------------------------------

func TestGenerateCodeVerifier(t *testing.T) {
	verifier, err := GenerateCodeVerifier()
	if err != nil {
		t.Fatalf("GenerateCodeVerifier() 鍑洪敊: %v", err)
	}
	if verifier == "" {
		t.Error("GenerateCodeVerifier() 杩斿洖绌哄瓧绗︿覆")
	}
	// RFC 7636 requires a code_verifier length of at least 43.
	if len(verifier) < 43 {
		t.Errorf("GenerateCodeVerifier() 闀垮害=%d锛孯FC 7636 瑕佹眰鑷冲皯 43 瀛楃", len(verifier))
	}
	// base64url 缂栫爜涓嶅簲鍖呭惈 padding 鍜岄潪 URL 瀹夊叏瀛楃
	if strings.Contains(verifier, "=") {
		t.Errorf("GenerateCodeVerifier() 鍖呭惈 '=' padding: %s", verifier)
	}
	if strings.ContainsAny(verifier, "+/") {
		t.Errorf("GenerateCodeVerifier() 鍖呭惈闈?base64url 瀛楃: %s", verifier)
	}
}

// ---------------------------------------------------------------------------
// GenerateCodeChallenge 娴嬭瘯
// ---------------------------------------------------------------------------

func TestGenerateCodeChallenge(t *testing.T) {
	// 浣跨敤宸茬煡杈撳叆楠岃瘉杈撳嚭
	// RFC 7636 闄勫綍 B 绀轰緥: verifier = "dBjftJeZ4CVP-mB92K27uhbUJU1p1r_wW1gFWFOEjXk"
	// 棰勬湡 challenge = "E9Melhoa2OwvFrEMTJguCHaoeK1t8URWbuGJSstw-cM"
	verifier := "dBjftJeZ4CVP-mB92K27uhbUJU1p1r_wW1gFWFOEjXk"
	expected := "E9Melhoa2OwvFrEMTJguCHaoeK1t8URWbuGJSstw-cM"

	challenge := GenerateCodeChallenge(verifier)
	if challenge != expected {
		t.Errorf("GenerateCodeChallenge(%q) = %q锛屾湡鏈?%q", verifier, challenge, expected)
	}
}

func TestGenerateCodeChallenge_NoPadding(t *testing.T) {
	challenge := GenerateCodeChallenge("test-verifier-string")
	if strings.Contains(challenge, "=") {
		t.Errorf("GenerateCodeChallenge() 缁撴灉鍖呭惈 '=' padding: %s", challenge)
	}
}

// ---------------------------------------------------------------------------
// base64URLEncode 娴嬭瘯
// ---------------------------------------------------------------------------

func TestBase64URLEncode(t *testing.T) {
	tests := []struct {
		name  string
		input []byte
	}{
		{"empty", []byte{}},
		{"single-byte", []byte{0xff}},
		{"multi-byte", []byte{0x01, 0x02, 0x03, 0x04, 0x05}},
		{"鍏ㄩ浂", []byte{0x00, 0x00, 0x00}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := base64URLEncode(tt.input)
			// 涓嶅簲鍖呭惈 '=' padding
			if strings.Contains(result, "=") {
				t.Errorf("base64URLEncode(%v) 鍖呭惈 '=' padding: %s", tt.input, result)
			}
			// 涓嶅簲鍖呭惈鏍囧噯 base64 鐨?'+' 鎴?'/'
			if strings.ContainsAny(result, "+/") {
				t.Errorf("base64URLEncode(%v) 鍖呭惈闈?URL 瀹夊叏瀛楃: %s", tt.input, result)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// hasRestrictedScope 娴嬭瘯
// ---------------------------------------------------------------------------

func TestHasRestrictedScope(t *testing.T) {
	tests := []struct {
		scope    string
		expected bool
	}{
		// 鍙楅檺 scope
		{"https://www.googleapis.com/auth/generative-language", true},
		{"https://www.googleapis.com/auth/generative-language.retriever", true},
		{"https://www.googleapis.com/auth/generative-language.tuning", true},
		{"https://www.googleapis.com/auth/drive", true},
		{"https://www.googleapis.com/auth/drive.readonly", true},
		{"https://www.googleapis.com/auth/drive.file", true},
		// 闈炲彈闄?scope
		{"https://www.googleapis.com/auth/cloud-platform", false},
		{"https://www.googleapis.com/auth/userinfo.email", false},
		{"https://www.googleapis.com/auth/userinfo.profile", false},
		// 杈圭晫鎯呭喌
		{"", false},
		{"random-scope", false},
	}
	for _, tt := range tests {
		t.Run(tt.scope, func(t *testing.T) {
			got := hasRestrictedScope(tt.scope)
			if got != tt.expected {
				t.Errorf("hasRestrictedScope(%q) = %v锛屾湡鏈?%v", tt.scope, got, tt.expected)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// BuildAuthorizationURL 娴嬭瘯
// ---------------------------------------------------------------------------

func TestBuildAuthorizationURL(t *testing.T) {
	t.Setenv(GeminiCLIOAuthClientSecretEnv, "test-secret")

	authURL, err := BuildAuthorizationURL(
		OAuthConfig{},
		"test-state",
		"test-challenge",
		"https://example.com/callback",
		"",
		"code_assist",
	)
	if err != nil {
		t.Fatalf("BuildAuthorizationURL() 鍑洪敊: %v", err)
	}

	// Returned URL should include the expected query parameters.
	checks := []string{
		"response_type=code",
		"client_id=" + GeminiCLIOAuthClientID,
		"redirect_uri=",
		"state=test-state",
		"code_challenge=test-challenge",
		"code_challenge_method=S256",
		"access_type=offline",
		"prompt=consent",
		"include_granted_scopes=true",
	}
	for _, check := range checks {
		if !strings.Contains(authURL, check) {
			t.Errorf("BuildAuthorizationURL() URL 缂哄皯鍙傛暟 %q\nURL: %s", check, authURL)
		}
	}

	// Empty projectID should not add a project_id query parameter.
	if strings.Contains(authURL, "project_id=") {
		t.Errorf("BuildAuthorizationURL() 绌?projectID 鏃朵笉搴斿寘鍚?project_id 鍙傛暟")
	}

	// URL should begin with the OAuth authorize endpoint.
	if !strings.HasPrefix(authURL, AuthorizeURL+"?") {
		t.Errorf("BuildAuthorizationURL() URL 搴斾互 %s? 寮€澶达紝瀹為檯: %s", AuthorizeURL, authURL)
	}
}

func TestBuildAuthorizationURL_EmptyRedirectURI(t *testing.T) {
	t.Setenv(GeminiCLIOAuthClientSecretEnv, "test-secret")

	_, err := BuildAuthorizationURL(
		OAuthConfig{},
		"test-state",
		"test-challenge",
		"", // 绌?redirectURI
		"",
		"code_assist",
	)
	if err == nil {
		t.Error("BuildAuthorizationURL() 绌?redirectURI 搴旇鎶ラ敊")
	}
	if !strings.Contains(err.Error(), "redirect_uri") {
		t.Errorf("閿欒娑堟伅搴斿寘鍚?'redirect_uri'锛屽疄闄? %v", err)
	}
}

func TestBuildAuthorizationURL_WithProjectID(t *testing.T) {
	t.Setenv(GeminiCLIOAuthClientSecretEnv, "test-secret")

	authURL, err := BuildAuthorizationURL(
		OAuthConfig{},
		"test-state",
		"test-challenge",
		"https://example.com/callback",
		"my-project-123",
		"code_assist",
	)
	if err != nil {
		t.Fatalf("BuildAuthorizationURL() 鍑洪敊: %v", err)
	}
	if !strings.Contains(authURL, "project_id=my-project-123") {
		t.Errorf("BuildAuthorizationURL() 甯?projectID 鏃跺簲鍖呭惈 project_id 鍙傛暟\nURL: %s", authURL)
	}
}

func TestBuildAuthorizationURL_UsesBuiltinSecretFallback(t *testing.T) {
	t.Setenv(GeminiCLIOAuthClientSecretEnv, "")

	authURL, err := BuildAuthorizationURL(
		OAuthConfig{},
		"test-state",
		"test-challenge",
		"https://example.com/callback",
		"",
		"code_assist",
	)
	if err == nil {
		t.Fatalf("BuildAuthorizationURL() should fail when %s is unset, got URL: %s", GeminiCLIOAuthClientSecretEnv, authURL)
	}
	if !strings.Contains(err.Error(), GeminiCLIOAuthClientSecretEnv) {
		t.Fatalf("expected error to mention %s, got: %v", GeminiCLIOAuthClientSecretEnv, err)
	}
}

// ---------------------------------------------------------------------------
// EffectiveOAuthConfig 娴嬭瘯 - 鍘熸湁娴嬭瘯
// ---------------------------------------------------------------------------

func TestEffectiveOAuthConfig_GoogleOne(t *testing.T) {
	// Provide a test secret for builtin-client cases.
	t.Setenv(GeminiCLIOAuthClientSecretEnv, "test-built-in-secret")

	tests := []struct {
		name         string
		input        OAuthConfig
		oauthType    string
		wantClientID string
		wantScopes   string
		wantErr      bool
	}{
		{
			name:         "Google One 浣跨敤鍐呯疆瀹㈡埛绔紙绌洪厤缃級",
			input:        OAuthConfig{},
			oauthType:    "google_one",
			wantClientID: GeminiCLIOAuthClientID,
			wantScopes:   DefaultCodeAssistScopes,
			wantErr:      false,
		},
		{
			name: "Google One uses custom client when credentials are provided",
			input: OAuthConfig{
				ClientID:     "custom-client-id",
				ClientSecret: "custom-client-secret",
			},
			oauthType:    "google_one",
			wantClientID: "custom-client-id",
			wantScopes:   DefaultCodeAssistScopes,
			wantErr:      false,
		},
		{
			name: "Google One builtin client filters restricted scopes",
			input: OAuthConfig{
				Scopes: "https://www.googleapis.com/auth/cloud-platform https://www.googleapis.com/auth/generative-language.retriever https://www.googleapis.com/auth/drive.readonly",
			},
			oauthType:    "google_one",
			wantClientID: GeminiCLIOAuthClientID,
			wantScopes:   "https://www.googleapis.com/auth/cloud-platform",
			wantErr:      false,
		},
		{
			name: "Google One 鍐呯疆瀹㈡埛绔?+ 浠呭彈闄?scopes锛堝簲鍥為€€鍒伴粯璁わ級",
			input: OAuthConfig{
				Scopes: "https://www.googleapis.com/auth/generative-language.retriever https://www.googleapis.com/auth/drive.readonly",
			},
			oauthType:    "google_one",
			wantClientID: GeminiCLIOAuthClientID,
			wantScopes:   DefaultCodeAssistScopes,
			wantErr:      false,
		},
		{
			name:         "Code Assist uses builtin client",
			input:        OAuthConfig{},
			oauthType:    "code_assist",
			wantClientID: GeminiCLIOAuthClientID,
			wantScopes:   DefaultCodeAssistScopes,
			wantErr:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := EffectiveOAuthConfig(tt.input, tt.oauthType)
			if (err != nil) != tt.wantErr {
				t.Errorf("EffectiveOAuthConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil {
				return
			}
			if got.ClientID != tt.wantClientID {
				t.Errorf("EffectiveOAuthConfig() ClientID = %v, want %v", got.ClientID, tt.wantClientID)
			}
			if got.Scopes != tt.wantScopes {
				t.Errorf("EffectiveOAuthConfig() Scopes = %v, want %v", got.Scopes, tt.wantScopes)
			}
		})
	}
}

func TestEffectiveOAuthConfig_ScopeFiltering(t *testing.T) {
	t.Setenv(GeminiCLIOAuthClientSecretEnv, "test-built-in-secret")

	// 娴嬭瘯 Google One + 鍐呯疆瀹㈡埛绔繃婊ゅ彈闄?scopes
	cfg, err := EffectiveOAuthConfig(OAuthConfig{
		Scopes: "https://www.googleapis.com/auth/cloud-platform https://www.googleapis.com/auth/userinfo.email https://www.googleapis.com/auth/generative-language.retriever https://www.googleapis.com/auth/drive.readonly https://www.googleapis.com/auth/userinfo.profile",
	}, "google_one")

	if err != nil {
		t.Fatalf("EffectiveOAuthConfig() error = %v", err)
	}

	// 搴斾粎鍖呭惈 cloud-platform銆乽serinfo.email 鍜?userinfo.profile
	// 涓嶅簲鍖呭惈 generative-language 鎴?drive scopes
	if strings.Contains(cfg.Scopes, "generative-language") {
		t.Errorf("浣跨敤鍐呯疆瀹㈡埛绔椂 Scopes 涓嶅簲鍖呭惈 generative-language锛屽疄闄? %v", cfg.Scopes)
	}
	if strings.Contains(cfg.Scopes, "drive") {
		t.Errorf("浣跨敤鍐呯疆瀹㈡埛绔椂 Scopes 涓嶅簲鍖呭惈 drive锛屽疄闄? %v", cfg.Scopes)
	}
	if !strings.Contains(cfg.Scopes, "cloud-platform") {
		t.Errorf("Scopes 搴斿寘鍚?cloud-platform锛屽疄闄? %v", cfg.Scopes)
	}
	if !strings.Contains(cfg.Scopes, "userinfo.email") {
		t.Errorf("Scopes 搴斿寘鍚?userinfo.email锛屽疄闄? %v", cfg.Scopes)
	}
	if !strings.Contains(cfg.Scopes, "userinfo.profile") {
		t.Errorf("Scopes 搴斿寘鍚?userinfo.profile锛屽疄闄? %v", cfg.Scopes)
	}
}

// ---------------------------------------------------------------------------
// EffectiveOAuthConfig 娴嬭瘯 - 鏂板鍒嗘敮瑕嗙洊
// ---------------------------------------------------------------------------

func TestEffectiveOAuthConfig_OnlyClientID_NoSecret(t *testing.T) {
	// Supplying only client ID without secret should fail.
	_, err := EffectiveOAuthConfig(OAuthConfig{
		ClientID: "some-client-id",
	}, "code_assist")
	if err == nil {
		t.Error("鍙彁渚?ClientID 涓嶆彁渚?ClientSecret 搴旇鎶ラ敊")
	}
	if !strings.Contains(err.Error(), "client_id") || !strings.Contains(err.Error(), "client_secret") {
		t.Errorf("閿欒娑堟伅搴旀彁鍙?client_id 鍜?client_secret锛屽疄闄? %v", err)
	}
}

func TestEffectiveOAuthConfig_OnlyClientSecret_NoID(t *testing.T) {
	// Supplying only client secret without client ID should fail.
	_, err := EffectiveOAuthConfig(OAuthConfig{
		ClientSecret: "some-client-secret",
	}, "code_assist")
	if err == nil {
		t.Error("鍙彁渚?ClientSecret 涓嶆彁渚?ClientID 搴旇鎶ラ敊")
	}
	if !strings.Contains(err.Error(), "client_id") || !strings.Contains(err.Error(), "client_secret") {
		t.Errorf("閿欒娑堟伅搴旀彁鍙?client_id 鍜?client_secret锛屽疄闄? %v", err)
	}
}

func TestEffectiveOAuthConfig_AIStudio_DefaultScopes_BuiltinClient(t *testing.T) {
	t.Setenv(GeminiCLIOAuthClientSecretEnv, "test-built-in-secret")

	// Builtin ai_studio client with empty scopes should fall back to default code-assist scopes.
	cfg, err := EffectiveOAuthConfig(OAuthConfig{}, "ai_studio")
	if err != nil {
		t.Fatalf("EffectiveOAuthConfig() error = %v", err)
	}
	if cfg.Scopes != DefaultCodeAssistScopes {
		t.Errorf("ai_studio + 鍐呯疆瀹㈡埛绔簲浣跨敤 DefaultCodeAssistScopes锛屽疄闄? %q", cfg.Scopes)
	}
}

func TestEffectiveOAuthConfig_AIStudio_DefaultScopes_CustomClient(t *testing.T) {
	// ai_studio 绫诲瀷锛屼娇鐢ㄨ嚜瀹氫箟瀹㈡埛绔紝scopes 涓虹┖ -> 搴斾娇鐢?DefaultAIStudioScopes
	cfg, err := EffectiveOAuthConfig(OAuthConfig{
		ClientID:     "custom-id",
		ClientSecret: "custom-secret",
	}, "ai_studio")
	if err != nil {
		t.Fatalf("EffectiveOAuthConfig() error = %v", err)
	}
	if cfg.Scopes != DefaultAIStudioScopes {
		t.Errorf("ai_studio + 鑷畾涔夊鎴风搴斾娇鐢?DefaultAIStudioScopes锛屽疄闄? %q", cfg.Scopes)
	}
}

func TestEffectiveOAuthConfig_AIStudio_ScopeNormalization(t *testing.T) {
	// ai_studio 绫诲瀷锛屾棫鐨?generative-language scope 搴旇褰掍竴鍖栦负 generative-language.retriever
	cfg, err := EffectiveOAuthConfig(OAuthConfig{
		ClientID:     "custom-id",
		ClientSecret: "custom-secret",
		Scopes:       "https://www.googleapis.com/auth/generative-language https://www.googleapis.com/auth/cloud-platform",
	}, "ai_studio")
	if err != nil {
		t.Fatalf("EffectiveOAuthConfig() error = %v", err)
	}
	if strings.Contains(cfg.Scopes, "auth/generative-language ") || strings.HasSuffix(cfg.Scopes, "auth/generative-language") {
		// Ensure the unnormalized generative-language scope is not present.
		parts := strings.Fields(cfg.Scopes)
		for _, p := range parts {
			if p == "https://www.googleapis.com/auth/generative-language" {
				t.Errorf("ai_studio 搴斿皢 generative-language 褰掍竴鍖栦负 generative-language.retriever锛屽疄闄?scopes: %q", cfg.Scopes)
			}
		}
	}
	if !strings.Contains(cfg.Scopes, "generative-language.retriever") {
		t.Errorf("ai_studio 褰掍竴鍖栧悗搴斿寘鍚?generative-language.retriever锛屽疄闄? %q", cfg.Scopes)
	}
}

func TestEffectiveOAuthConfig_CommaSeparatedScopes(t *testing.T) {
	t.Setenv(GeminiCLIOAuthClientSecretEnv, "test-built-in-secret")

	// 閫楀彿鍒嗛殧鐨?scopes 搴旇褰掍竴鍖栦负绌烘牸鍒嗛殧
	cfg, err := EffectiveOAuthConfig(OAuthConfig{
		ClientID:     "custom-id",
		ClientSecret: "custom-secret",
		Scopes:       "https://www.googleapis.com/auth/cloud-platform,https://www.googleapis.com/auth/userinfo.email",
	}, "code_assist")
	if err != nil {
		t.Fatalf("EffectiveOAuthConfig() error = %v", err)
	}
	// 搴旇鐢ㄧ┖鏍煎垎闅旓紝鑰岄潪閫楀彿
	if strings.Contains(cfg.Scopes, ",") {
		t.Errorf("閫楀彿鍒嗛殧鐨?scopes 搴旇褰掍竴鍖栦负绌烘牸鍒嗛殧锛屽疄闄? %q", cfg.Scopes)
	}
	if !strings.Contains(cfg.Scopes, "cloud-platform") {
		t.Errorf("褰掍竴鍖栧悗搴斿寘鍚?cloud-platform锛屽疄闄? %q", cfg.Scopes)
	}
	if !strings.Contains(cfg.Scopes, "userinfo.email") {
		t.Errorf("褰掍竴鍖栧悗搴斿寘鍚?userinfo.email锛屽疄闄? %q", cfg.Scopes)
	}
}

func TestEffectiveOAuthConfig_MixedCommaAndSpaceScopes(t *testing.T) {
	// 娣峰悎閫楀彿鍜岀┖鏍煎垎闅旂殑 scopes
	cfg, err := EffectiveOAuthConfig(OAuthConfig{
		ClientID:     "custom-id",
		ClientSecret: "custom-secret",
		Scopes:       "https://www.googleapis.com/auth/cloud-platform, https://www.googleapis.com/auth/userinfo.email https://www.googleapis.com/auth/userinfo.profile",
	}, "code_assist")
	if err != nil {
		t.Fatalf("EffectiveOAuthConfig() error = %v", err)
	}
	parts := strings.Fields(cfg.Scopes)
	if len(parts) != 3 {
		t.Errorf("褰掍竴鍖栧悗搴旀湁 3 涓?scope锛屽疄闄? %d锛宻copes: %q", len(parts), cfg.Scopes)
	}
}

func TestEffectiveOAuthConfig_WhitespaceTriming(t *testing.T) {
	// 杈撳叆涓殑鍓嶅悗绌虹櫧搴旇娓呯悊
	cfg, err := EffectiveOAuthConfig(OAuthConfig{
		ClientID:     "  custom-id  ",
		ClientSecret: "  custom-secret  ",
		Scopes:       "  https://www.googleapis.com/auth/cloud-platform  ",
	}, "code_assist")
	if err != nil {
		t.Fatalf("EffectiveOAuthConfig() error = %v", err)
	}
	if cfg.ClientID != "custom-id" {
		t.Errorf("ClientID 搴斿幓闄ゅ墠鍚庣┖鐧斤紝瀹為檯: %q", cfg.ClientID)
	}
	if cfg.ClientSecret != "custom-secret" {
		t.Errorf("ClientSecret 搴斿幓闄ゅ墠鍚庣┖鐧斤紝瀹為檯: %q", cfg.ClientSecret)
	}
	if cfg.Scopes != "https://www.googleapis.com/auth/cloud-platform" {
		t.Errorf("Scopes 搴斿幓闄ゅ墠鍚庣┖鐧斤紝瀹為檯: %q", cfg.Scopes)
	}
}

func TestEffectiveOAuthConfig_NoEnvSecret(t *testing.T) {
	t.Setenv(GeminiCLIOAuthClientSecretEnv, "")

	cfg, err := EffectiveOAuthConfig(OAuthConfig{}, "code_assist")
	if err == nil {
		t.Fatalf("EffectiveOAuthConfig() should fail when %s is unset, got cfg=%+v", GeminiCLIOAuthClientSecretEnv, cfg)
	}
	if !strings.Contains(err.Error(), GeminiCLIOAuthClientSecretEnv) {
		t.Fatalf("expected error to mention %s, got: %v", GeminiCLIOAuthClientSecretEnv, err)
	}
}

func TestEffectiveOAuthConfig_AIStudio_BuiltinClient_CustomScopes(t *testing.T) {
	t.Setenv(GeminiCLIOAuthClientSecretEnv, "test-built-in-secret")

	// ai_studio + 鍐呯疆瀹㈡埛绔?+ 鑷畾涔?scopes -> 搴旇繃婊ゅ彈闄?scopes
	cfg, err := EffectiveOAuthConfig(OAuthConfig{
		Scopes: "https://www.googleapis.com/auth/cloud-platform https://www.googleapis.com/auth/generative-language.retriever",
	}, "ai_studio")
	if err != nil {
		t.Fatalf("EffectiveOAuthConfig() error = %v", err)
	}
	// 鍐呯疆瀹㈡埛绔簲杩囨护 generative-language.retriever
	if strings.Contains(cfg.Scopes, "generative-language") {
		t.Errorf("ai_studio + 鍐呯疆瀹㈡埛绔簲杩囨护鍙楅檺 scopes锛屽疄闄? %q", cfg.Scopes)
	}
	if !strings.Contains(cfg.Scopes, "cloud-platform") {
		t.Errorf("搴斾繚鐣?cloud-platform scope锛屽疄闄? %q", cfg.Scopes)
	}
}

func TestEffectiveOAuthConfig_UnknownOAuthType_DefaultScopes(t *testing.T) {
	t.Setenv(GeminiCLIOAuthClientSecretEnv, "test-built-in-secret")

	// 鏈煡鐨?oauthType 搴斿洖閫€鍒伴粯璁ょ殑 code_assist scopes
	cfg, err := EffectiveOAuthConfig(OAuthConfig{}, "unknown_type")
	if err != nil {
		t.Fatalf("EffectiveOAuthConfig() error = %v", err)
	}
	if cfg.Scopes != DefaultCodeAssistScopes {
		t.Errorf("鏈煡 oauthType 搴斾娇鐢?DefaultCodeAssistScopes锛屽疄闄? %q", cfg.Scopes)
	}
}

func TestEffectiveOAuthConfig_EmptyOAuthType_DefaultScopes(t *testing.T) {
	t.Setenv(GeminiCLIOAuthClientSecretEnv, "test-built-in-secret")

	// 绌虹殑 oauthType 搴旇蛋 default 鍒嗘敮锛屼娇鐢?DefaultCodeAssistScopes
	cfg, err := EffectiveOAuthConfig(OAuthConfig{}, "")
	if err != nil {
		t.Fatalf("EffectiveOAuthConfig() error = %v", err)
	}
	if cfg.Scopes != DefaultCodeAssistScopes {
		t.Errorf("绌?oauthType 搴斾娇鐢?DefaultCodeAssistScopes锛屽疄闄? %q", cfg.Scopes)
	}
}

func TestEffectiveOAuthConfig_CustomClient_NoScopeFiltering(t *testing.T) {
	// 鑷畾涔夊鎴风 + google_one + 鍖呭惈鍙楅檺 scopes -> 涓嶅簲琚繃婊わ紙鍥犱负涓嶆槸鍐呯疆瀹㈡埛绔級
	cfg, err := EffectiveOAuthConfig(OAuthConfig{
		ClientID:     "custom-id",
		ClientSecret: "custom-secret",
		Scopes:       "https://www.googleapis.com/auth/generative-language.retriever https://www.googleapis.com/auth/drive.readonly",
	}, "google_one")
	if err != nil {
		t.Fatalf("EffectiveOAuthConfig() error = %v", err)
	}
	// 鑷畾涔夊鎴风涓嶅簲杩囨护浠讳綍 scope
	if !strings.Contains(cfg.Scopes, "generative-language.retriever") {
		t.Errorf("鑷畾涔夊鎴风涓嶅簲杩囨护 generative-language.retriever锛屽疄闄? %q", cfg.Scopes)
	}
	if !strings.Contains(cfg.Scopes, "drive.readonly") {
		t.Errorf("鑷畾涔夊鎴风涓嶅簲杩囨护 drive.readonly锛屽疄闄? %q", cfg.Scopes)
	}
}

