//go:build unit

package antigravity

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"
)

// ---------------------------------------------------------------------------
// getClientSecret
// ---------------------------------------------------------------------------

func TestGetClientSecret_闂佺粯绮犻崹浼淬€傞妸鈺佺煑婵せ鍋撻柛锝堟閹峰宕滆閺?t *testing.T) {
	old := defaultClientSecret
	defaultClientSecret = ""
	t.Cleanup(func() { defaultClientSecret = old })
	t.Setenv(AntigravityOAuthClientSecretEnv, "my-secret-value")

	// 闂傚倸娲犻崑鎾绘偡閺囨氨绐旈柛锝嗘倐瀵剟骞忕仦婵囨礋瀹?init 闂備緡鍋呭Σ鎺旀椤愶附鏅慨姗嗗幖椤や線鏌涢弬璇插缂侇喚濞€閹娊顢涘鍏兼闂佸憡鐟﹂敃銏ゅ闯閾忚瀚氶悹鍥ㄥ絻缁?
	defaultClientSecret = os.Getenv(AntigravityOAuthClientSecretEnv)

	secret, err := getClientSecret()
	if err != nil {
		t.Fatalf("闂佸吋鍎抽崲鑼躲亹?client_secret 婵犮垺鍎肩划鍓ф喆? %v", err)
	}
	if secret != "my-secret-value" {
		t.Errorf("client_secret 婵炴垶鎸哥粔瀵镐焊椤曗偓閺? got %s, want my-secret-value", secret)
	}
}

func TestGetClientSecret_闂佺粯绮犻崹浼淬€傞妸鈺佺煑婵せ鍋撻柛锝囧劋缁嬪鎼归崗澶嬫櫈(t *testing.T) {
	old := defaultClientSecret
	defaultClientSecret = ""
	t.Cleanup(func() { defaultClientSecret = old })

	_, err := getClientSecret()
	if err == nil {
		t.Fatal("defaultClientSecret 婵炴垶鎹佸▍锝夊煘閺嶎厼绫嶉悹楦挎鐎瑰寮堕埡鍌涚叆婵炲弶鐗犻弻銊モ枎閹烘繂娈?)
	}
	if !strings.Contains(err.Error(), AntigravityOAuthClientSecretEnv) {
		t.Errorf("闂備焦瀵ч悷銊╊敋閵堝棛鈹嶉柍鈺佸暕缁辨牠骞栫€涙ɑ鐓ラ悗鍨耿瀹曘儵顢曢妶鍛暰婵犫拃鍐ㄦ殭鐟滄澘顦甸弻宀冪疀閹炬潙鈧? got %s", err.Error())
	}
}

func TestGetClientSecret_闂佺粯绮犻崹浼淬€傞妸鈺佺煑婵せ鍋撻柛锝嗘そ瀵敻顢旈崶璺烘櫗缂?t *testing.T) {
	old := defaultClientSecret
	defaultClientSecret = ""
	t.Cleanup(func() { defaultClientSecret = old })

	_, err := getClientSecret()
	if err == nil {
		t.Fatal("defaultClientSecret 婵炴垶鎹佸▍锝夊煘閺嶎厼绫嶉悹楦挎鐎瑰寮堕埡鍌涚叆婵炲弶鐗犻弻銊モ枎閹烘繂娈?)
	}
}

func TestGetClientSecret_闂佺粯绮犻崹浼淬€傞妸鈺佺煑婵せ鍋撻柛锝嗘そ瀹曘儵顢曢妶鍌涙櫈闂?t *testing.T) {
	old := defaultClientSecret
	defaultClientSecret = "   "
	t.Cleanup(func() { defaultClientSecret = old })

	_, err := getClientSecret()
	if err == nil {
		t.Fatal("defaultClientSecret 婵炲濮撮幊搴ㄥ箚閸垻鐭氶柣鎴烆焽婢规劙鏌￠崘顓熺【缂併劌寮跺濠氬棘閹稿海顦ラ梻浣瑰閻熴劑顢?)
	}
}

func TestGetClientSecret_闂佺粯绮犻崹浼淬€傞妸鈺佺煑婵せ鍋撻柛锝嗘そ瀵灚寰勬繝鍌ゆ瀫闂佸憡鑹惧ù鐑藉煘閺嶎厼鍐€?t *testing.T) {
	old := defaultClientSecret
	defaultClientSecret = "  valid-secret  "
	t.Cleanup(func() { defaultClientSecret = old })

	secret, err := getClientSecret()
	if err != nil {
		t.Fatalf("闂佸吋鍎抽崲鑼躲亹?client_secret 婵犮垺鍎肩划鍓ф喆? %v", err)
	}
	if secret != "valid-secret" {
		t.Errorf("闁圭厧鐡ㄩ弻銊╃嵁閹剧粯鈷旈柕鍫濇噹椤ゅ懘鏌涘顒佹崳闁宠鐗犲? got %q, want %q", secret, "valid-secret")
	}
}

// ---------------------------------------------------------------------------
// ForwardBaseURLs
// ---------------------------------------------------------------------------

func TestForwardBaseURLs_Daily婵炴潙鍚嬮敋闁?t *testing.T) {
	urls := ForwardBaseURLs()
	if len(urls) == 0 {
		t.Fatal("ForwardBaseURLs 闁哄鏅滈弻銊ッ洪弽顐ょ煔闁告繂瀚悘娆撴偠?)
	}

	// daily URL 闁圭厧鐡ㄥ鍦暜閹捐鎹堕柕濞у苯鍓虫繛鎴炴尨閸嬫挸霉?	if urls[0] != antigravityDailyBaseURL {
		t.Errorf("缂備焦顨忛崗娑氱博鐎涙鈻?URL 闁圭厧鐡ㄩ弻褏鎷?daily: got %s, want %s", urls[0], antigravityDailyBaseURL)
	}

	// 闁圭厧鐡ㄩ弻銊р偓鍨耿瀹曘儵顢曢敐鍜佹殹闂?URL
	if len(urls) != len(BaseURLs) {
		t.Errorf("URL 闂佽桨妞掗崡鎶藉闯閻戞鈻旂€广儱鎳庨悥閬嶆⒑? got %d, want %d", len(urls), len(BaseURLs))
	}

	// 婵°倗濮撮惌渚€鎯?prod URL 婵炴垶姊婚崰搴★耿椤忓牆绀嗘俊銈呭閳ь剙鍟粙?	found := false
	for _, u := range urls {
		if u == antigravityProdBaseURL {
			found = true
			break
		}
	}
	if !found {
		t.Error("ForwardBaseURLs 婵炴垶鎼╅崢鎯ь啅鏉堚晙鐒?prod URL")
	}
}

func TestForwardBaseURLs_婵炴垶鎸哥粔鎶藉箞閵娾晛缁╅柣鐔告緲閺傃囨煕閹烘垶宸濆?t *testing.T) {
	originalFirst := BaseURLs[0]
	_ = ForwardBaseURLs()
	// 缂佺虎鍙庨崰鏇犳崲濮樿泛鍌ㄩ柣鏂款殠濞?BaseURLs 闂佸搫鐗滄禍锝夛綖閿旂晫鈹嶆い鏃囧Г閺?
	if BaseURLs[0] != originalFirst {
		t.Errorf("ForwardBaseURLs 婵炴垶鎸哥粔瀵歌姳閸欏鈹嶆い鏃囧Г閺嗩參鏌涘Ο鐓庢瀻妞?BaseURLs: got %s, want %s", BaseURLs[0], originalFirst)
	}
}

// ---------------------------------------------------------------------------
// URLAvailability
// ---------------------------------------------------------------------------

func TestNewURLAvailability(t *testing.T) {
	ua := NewURLAvailability(5 * time.Minute)
	if ua == nil {
		t.Fatal("NewURLAvailability 闁哄鏅滈弻銊ッ?nil")
	}
	if ua.ttl != 5*time.Minute {
		t.Errorf("TTL 婵炴垶鎸哥粔瀵镐焊椤曗偓閺? got %v, want 5m", ua.ttl)
	}
	if ua.unavailable == nil {
		t.Error("unavailable map 婵炴垶鎸哥粔瀵歌姳閸欏鈻?nil")
	}
}

func TestURLAvailability_MarkUnavailable(t *testing.T) {
	ua := NewURLAvailability(5 * time.Minute)
	testURL := "https://example.com"

	ua.MarkUnavailable(testURL)

	if ua.IsAvailable(testURL) {
		t.Error("闂佸搫绉村ú鈺咁敊閸ャ劎鈻旈弶鐐村閻熸繈鏌涘▎妯虹仯闁轰降鍊濆畷?IsAvailable 闁圭厧鐡ㄥΛ浣烘崲閹达箑鐐?false")
	}
}

func TestURLAvailability_MarkSuccess(t *testing.T) {
	ua := NewURLAvailability(5 * time.Minute)
	testURL := "https://example.com"

	// 闂佺绻愰悧濠囨偉閿濆洦濯煎Λ棰佹祰缁€瀣槈閹惧磭孝鐟滅増鐓￠幃?	ua.MarkUnavailable(testURL)
	if ua.IsAvailable(testURL) {
		t.Error("闂佸搫绉村ú鈺咁敊閸ャ劎鈻旈弶鐐村閻熸繈鏌涘▎妯虹仯闁轰降鍊濆畷銉︽償閿涘嫬鐣ㄦ繛鎴炴尭缁夌銇愰弻銉﹀仺?)
	}

	// 闂佸搫绉村ú鈺咁敊閸ヮ剙绠ｉ柟閭﹀墮椤娀鏌涘顒勵€楃紒銊︾叀楠炰線濮€閻欌偓濡插鏌涘▎妯虹仯闁?
	ua.MarkSuccess(testURL)
	if !ua.IsAvailable(testURL) {
		t.Error("MarkSuccess 闂佸憡鑹炬鍝ヨ姳閺屻儱绠掗柕蹇曞濡插鏌涘▎妯虹仯闁?)
	}

	// 婵°倗濮撮惌渚€鎯?lastSuccess 闁荤偞鍑归崑澶愵敊閺囩姷纾?	ua.mu.RLock()
	if ua.lastSuccess != testURL {
		t.Errorf("lastSuccess 婵炴垶鎸哥粔瀵镐焊椤曗偓閺? got %s, want %s", ua.lastSuccess, testURL)
	}
	ua.mu.RUnlock()
}

func TestURLAvailability_IsAvailable_TTL闁哄鏅涘ú锕€锕?t *testing.T) {
	// 婵炶揪缍€濞夋洟寮妶澶婂嚑濞达綀娅ｉ崣姘舵煟?TTL
	ua := NewURLAvailability(1 * time.Millisecond)
	testURL := "https://example.com"

	ua.MarkUnavailable(testURL)
	// 缂備焦绋戦ˇ顖滄?TTL 闁哄鏅涘ú锕€锕?
	time.Sleep(5 * time.Millisecond)

	if !ua.IsAvailable(testURL) {
		t.Error("TTL 闁哄鏅涘ú锕€锕㈤敓鐘茶Е?URL 闁圭厧鐡ㄥ瑙勭椤旇棄绶炵€广儱鎳庣拋鏌ユ煟?)
	}
}

func TestURLAvailability_IsAvailable_闂佸搫鐗滄禍婵嬫偉閿濆洦濯奸柟娈垮枟閻ｇ洩RL(t *testing.T) {
	ua := NewURLAvailability(5 * time.Minute)
	if !ua.IsAvailable("https://never-marked.com") {
		t.Error("闂佸搫鐗滄禍婵嬫偉閿濆洦濯奸柟娈垮枟閻?URL 闁圭厧鐡ㄥΛ鍐垝椤栨粍濯奸柕鍫濇噹鐠佹煡鏌?)
	}
}

func TestURLAvailability_GetAvailableURLs(t *testing.T) {
	ua := NewURLAvailability(10 * time.Minute)

	// 婵帗绋掗…鍫ヮ敇婵犳艾绠ラ柍褜鍓熷?URL 闂備緡鍠涘Λ鍕亹閺屻儲鍋?	urls := ua.GetAvailableURLs()
	if len(urls) != len(BaseURLs) {
		t.Errorf("闂佸憡鐟崹鎶藉极?URL 闂佽桨妞掗崡鎶藉闯閻戞鈻旂€广儱鎳庨悥閬嶆⒑? got %d, want %d", len(urls), len(BaseURLs))
	}
}

func TestURLAvailability_GetAvailableURLs_闂佸搫绉村ú鈺咁敊閸ャ劎鈻旈柍褜鍓氱粙澶愵敂閸愵亞鎲归梺鍛婄懐閸ㄦ娊寮?t *testing.T) {
	ua := NewURLAvailability(10 * time.Minute)

	if len(BaseURLs) < 2 {
		t.Skip("BaseURLs 闁诲繐绻戦崹宕囪姳?2 婵炴垶鎼╂禍椋庢濠靛牊宕夐悗鍦Х缁犳牗鎱ㄥ┑鍕偓鏍矈鐎靛憡瀚?)
	}

	ua.MarkUnavailable(BaseURLs[0])
	urls := ua.GetAvailableURLs()

	// 闂佸搫绉村ú鈺咁敊閸ヮ剚鍎?URL 婵炴垶鎸哥粔瀵歌姳閺屻儱绀勯柧蹇曟嚀缁犳盯鏌涢敂鍝勫鐟滅増鐓￠幃浠嬪Ω閵夈儳浠氶柣鐐寸◤閸斿鎳?
	for _, u := range urls {
		if u == BaseURLs[0] {
			t.Errorf("闁荤偞鍑归崑鍛存偉閿濆洦濯煎Λ棰佽兌閻熸繈鏌涘▎妯虹仯闁轰降鍊濋幆?URL 婵炴垶鎸哥粔瀵歌姳閺屻儱绀勯柧蹇曟嚀缁犳盯鏌涢敂鍝勫鐟滅増鐓￠幃浠嬪Ω閵夈儳浠氶柣鐐寸◤閸斿鎳? %s", BaseURLs[0])
		}
	}
}

func TestURLAvailability_GetAvailableURLsWithBase(t *testing.T) {
	ua := NewURLAvailability(10 * time.Minute)
	customURLs := []string{"https://a.com", "https://b.com", "https://c.com"}

	urls := ua.GetAvailableURLsWithBase(customURLs)
	if len(urls) != 3 {
		t.Errorf("闂佸憡鐟崹鎶藉极?URL 闂佽桨妞掗崡鎶藉闯閻戞鈻旂€广儱鎳庨悥閬嶆⒑? got %d, want 3", len(urls))
	}
}

func TestURLAvailability_GetAvailableURLsWithBase_LastSuccess婵炴潙鍚嬮敋闁?t *testing.T) {
	ua := NewURLAvailability(10 * time.Minute)
	customURLs := []string{"https://a.com", "https://b.com", "https://c.com"}

	ua.MarkSuccess("https://c.com")

	urls := ua.GetAvailableURLsWithBase(customURLs)
	if len(urls) != 3 {
		t.Fatalf("闂佸憡鐟崹鎶藉极?URL 闂佽桨妞掗崡鎶藉闯閻戞鈻旂€广儱鎳庨悥閬嶆⒑? got %d, want 3", len(urls))
	}
	// c.com 闁圭厧鐡ㄥ鍦暜閹捐鎹堕柕濞у苯鍓虫繛鎴炴尨閸嬫挸霉?	if urls[0] != "https://c.com" {
		t.Errorf("lastSuccess 闁圭厧鐡ㄥ鍦暜閹捐鎹堕柕濞у苯鍓虫繛鎴炴尨閸嬫挸霉? got %s, want https://c.com", urls[0])
	}
	// 闂佺绻戝﹢鍦礊閹达箑绠板璺猴工閺傃冣攽椤旂⒈鍎旈柕鍡楀閹?	if urls[1] != "https://a.com" {
		t.Errorf("缂備焦顨忛崗娑氳姳閳哄倻鈻旀い蹇撳鐎瑰鈽?a.com: got %s", urls[1])
	}
	if urls[2] != "https://b.com" {
		t.Errorf("缂備焦顨忛崗娑氱箔娴ｅ湱鈻旀い蹇撳鐎瑰鈽?b.com: got %s", urls[2])
	}
}

func TestURLAvailability_GetAvailableURLsWithBase_LastSuccess婵炴垶鎸哥粔纾嬨亹閺屻儲鍋?t *testing.T) {
	ua := NewURLAvailability(10 * time.Minute)
	customURLs := []string{"https://a.com", "https://b.com"}

	ua.MarkSuccess("https://b.com")
	ua.MarkUnavailable("https://b.com")

	urls := ua.GetAvailableURLsWithBase(customURLs)
	// b.com 闁荤偞鍑归崑鍛存偉閿濆洦濯煎Λ棰佽兌閻熸繈鏌涘▎妯虹仯闁轰降鍊濋弫宥囦沪閼测晝鎲归柟鐓庣摠閺屻劑宕甸銏″仢?	if len(urls) != 1 {
		t.Fatalf("闂佸憡鐟崹鎶藉极?URL 闂佽桨妞掗崡鎶藉闯閻戞鈻旂€广儱鎳庨悥閬嶆⒑? got %d, want 1", len(urls))
	}
	if urls[0] != "https://a.com" {
		t.Errorf("婵?a.com 闁圭厧鐡ㄩ弻銊ㄣ亹閺屻儲鍋? got %s", urls[0])
	}
}

func TestURLAvailability_GetAvailableURLsWithBase_LastSuccess婵炴垶鎸哥粔鏉戯耿椤忓牆绀嗘俊銈呭閳ь剙鍟粙?t *testing.T) {
	ua := NewURLAvailability(10 * time.Minute)
	customURLs := []string{"https://a.com", "https://b.com"}

	ua.MarkSuccess("https://not-in-list.com")

	urls := ua.GetAvailableURLsWithBase(customURLs)
	// lastSuccess 婵炴垶鎸哥粔鏉戯耿椤忓牊鍤婃い蹇撳閺嗘澘鈽夐弬娆炬Ц闁割煈浜為幃浼村Ω閵堝牆骞€闂佹寧绋戞總鏃傜箔婢跺鍎熼柡鍐ㄦ祩濞艰埖绻涢敐鍫殭濠?
	if len(urls) != 2 {
		t.Fatalf("闂佸憡鐟崹鎶藉极?URL 闂佽桨妞掗崡鎶藉闯閻戞鈻旂€广儱鎳庨悥閬嶆⒑? got %d, want 2", len(urls))
	}
}

// ---------------------------------------------------------------------------
// SessionStore
// ---------------------------------------------------------------------------

func TestNewSessionStore(t *testing.T) {
	store := NewSessionStore()
	defer store.Stop()

	if store == nil {
		t.Fatal("NewSessionStore 闁哄鏅滈弻銊ッ?nil")
	}
	if store.sessions == nil {
		t.Error("sessions map 婵炴垶鎸哥粔瀵歌姳閸欏鈻?nil")
	}
}

func TestSessionStore_SetAndGet(t *testing.T) {
	store := NewSessionStore()
	defer store.Stop()

	session := &OAuthSession{
		State:        "test-state",
		CodeVerifier: "test-verifier",
		ProxyURL:     "http://proxy.example.com",
		CreatedAt:    time.Now(),
	}

	store.Set("session-1", session)

	got, ok := store.Get("session-1")
	if !ok {
		t.Fatal("Get 闁圭厧鐡ㄥΛ浣烘崲閹达箑鐐?true")
	}
	if got.State != "test-state" {
		t.Errorf("State 婵炴垶鎸哥粔瀵镐焊椤曗偓閺? got %s", got.State)
	}
	if got.CodeVerifier != "test-verifier" {
		t.Errorf("CodeVerifier 婵炴垶鎸哥粔瀵镐焊椤曗偓閺? got %s", got.CodeVerifier)
	}
	if got.ProxyURL != "http://proxy.example.com" {
		t.Errorf("ProxyURL 婵炴垶鎸哥粔瀵镐焊椤曗偓閺? got %s", got.ProxyURL)
	}
}

func TestSessionStore_Get_婵炴垶鎸哥粔鎾偤閵娾晛鎹?t *testing.T) {
	store := NewSessionStore()
	defer store.Stop()

	_, ok := store.Get("nonexistent")
	if ok {
		t.Error("婵炴垶鎸哥粔鎾偤閵娾晛鎹堕柕濞у嫮鏆?session 闁圭厧鐡ㄥΛ浣烘崲閹达箑鐐?false")
	}
}

func TestSessionStore_Get_闁哄鏅涘ú锕€锕?t *testing.T) {
	store := NewSessionStore()
	defer store.Stop()

	session := &OAuthSession{
		State:     "expired-state",
		CreatedAt: time.Now().Add(-SessionTTL - time.Minute), // 閻庤鐡曠亸顏嗘崲閸愵喖瀚?	}

	store.Set("expired-session", session)

	_, ok := store.Get("expired-session")
	if ok {
		t.Error("闁哄鏅涘ú锕€锕㈤敓鐘冲剭?session 闁圭厧鐡ㄥΛ浣烘崲閹达箑鐐?false")
	}
}

func TestSessionStore_Delete(t *testing.T) {
	store := NewSessionStore()
	defer store.Stop()

	session := &OAuthSession{
		State:     "to-delete",
		CreatedAt: time.Now(),
	}

	store.Set("del-session", session)
	store.Delete("del-session")

	_, ok := store.Get("del-session")
	if ok {
		t.Error("闂佸憡甯炴繛鈧繛鍛叄瀹?Get 闁圭厧鐡ㄥΛ浣烘崲閹达箑鐐?false")
	}
}

func TestSessionStore_Delete_婵炴垶鎸哥粔鎾偤閵娾晛鎹?t *testing.T) {
	store := NewSessionStore()
	defer store.Stop()

	// 闂佸憡甯炴繛鈧繛鍛缁嬪顓奸崨顖涙喖闂侀潻璐熼崝搴♀枔?session 婵炴垶鎸哥粔瀵歌姳?panic
	store.Delete("nonexistent")
}

func TestSessionStore_Stop(t *testing.T) {
	store := NewSessionStore()
	store.Stop()

	// 婵犮垼鍩栫喊宥夘敃?Stop 婵炴垶鎸哥粔瀵歌姳?panic
	store.Stop()
}

func TestSessionStore_婵犮垼鍩栭惌顕€鏌屽ǎ鐚爏sion(t *testing.T) {
	store := NewSessionStore()
	defer store.Stop()

	for i := 0; i < 10; i++ {
		session := &OAuthSession{
			State:     "state-" + string(rune('0'+i)),
			CreatedAt: time.Now(),
		}
		store.Set("session-"+string(rune('0'+i)), session)
	}

	// 婵°倗濮撮惌渚€鎯佹径鎰劸闁稿﹦鍠庨崢鎾煕濞嗘劗澧柛?
	for i := 0; i < 10; i++ {
		_, ok := store.Get("session-" + string(rune('0'+i)))
		if !ok {
			t.Errorf("session-%d 闁圭厧鐡ㄩ弻銊╂偤閵娾晛鎹?, i)
		}
	}
}

// ---------------------------------------------------------------------------
// GenerateRandomBytes
// ---------------------------------------------------------------------------

func TestGenerateRandomBytes_闂傚倵鍋撻柛顭戝枛椤斿﹥鎱ㄥ┑鎾舵偧闁?t *testing.T) {
	sizes := []int{0, 1, 16, 32, 64, 128}
	for _, size := range sizes {
		b, err := GenerateRandomBytes(size)
		if err != nil {
			t.Fatalf("GenerateRandomBytes(%d) 婵犮垺鍎肩划鍓ф喆? %v", size, err)
		}
		if len(b) != size {
			t.Errorf("闂傚倵鍋撻柛顭戝枛椤斿﹤鈽夐幘宕囆㈤悘蹇ｅ櫍閺? got %d, want %d", len(b), size)
		}
	}
}

func TestGenerateRandomBytes_婵炴垶鎸哥粔鎾箖閹捐埖瀚柛鎰典簼閺嗗繐霉濠х姴鎳忛弲绋库槈閹惧磭孝闁诡喗鎸剧槐鎺楀箻鐎甸晲鍑?t *testing.T) {
	b1, err := GenerateRandomBytes(32)
	if err != nil {
		t.Fatalf("缂備焦顨忛崗娑氱博閺夋垟鏋庨梽鍥儍閻斿吋鍋ㄩ柕濞垮劘娴滃ジ鎮? %v", err)
	}
	b2, err := GenerateRandomBytes(32)
	if err != nil {
		t.Fatalf("缂備焦顨忛崗娑氳姳閳轰讲鏋庨梽鍥儍閻斿吋鍋ㄩ柕濞垮劘娴滃ジ鎮? %v", err)
	}
	// 婵炴垶鎸堕崐鏍敃婵傚憡鍋ㄩ柣鏃傤焾閻忓洭鏌ｉ妸銉ヮ伃婵炲牄鍨藉鐢稿传閸曨厽鎲婚梺鐓庢惈閸婂摜鑺遍懠顒佸珰闁靛鍊楅悷婵嬫煕濮橆剛鐒跨紒杈ㄧ懃椤宕掑鍐У婵炴垶鎸搁敃銈夊吹閹寸偟鈻曢幖瀛樼箘閻熸繈鏌涘▎妯虹仴闁稿繑锕㈤幆鍕矙閹稿骸鈧亶鏌?	if string(b1) == string(b2) {
		t.Error("婵炴垶鎸堕崐鏍敃婵傚憡鍋ㄩ柣鏃傤焾閻忓洭鏌ｉ妸銉ヮ伃婵炲牄鍨藉鐢稿传閸曨厽鎲婚梺鐓庢惈閸婅霉婢舵劕瑙︾€光偓鐎ｎ剛顦┑顔藉笒閸婅顔忔總绋垮嚑濞达綁顥撶粔鐢告煥濞戞瀚扮憸鐗堢叀閹虫鎳為妷锔剧暢闂傚倸鍋嗛崳锝夈€?)
	}
}

// ---------------------------------------------------------------------------
// GenerateState
// ---------------------------------------------------------------------------

func TestGenerateState_闁哄鏅滈弻銊ッ洪弽顓炵９闁稿繗鍋愭竟鎰偓?t *testing.T) {
	state, err := GenerateState()
	if err != nil {
		t.Fatalf("GenerateState 婵犮垺鍎肩划鍓ф喆? %v", err)
	}
	if state == "" {
		t.Error("GenerateState 闁哄鏅滈弻銊ッ洪弽顐ょ煔闁告繂瀚幗鐔虹磼濡ゅ绱伴悷?)
	}
	// base64url 缂傚倸鍊归悧婊堟偉濠婂嫮鈻旂€广儱鎳愮€瑰鏌涢弽褎鍣归柟?+, /, =
	if strings.ContainsAny(state, "+/=") {
		t.Errorf("GenerateState 闁哄鏅滈弻銊ッ洪弽顓炵９闁绘挸楠搁惁鍫曟煕濮樼厧浜版繛?base64url 闁诲孩绋掗〃鍫ヮ敄? %s", state)
	}
	// 32 闁诲孩绋掗〃澶嬩繆椤撱垺鍎?base64url 缂傚倸鍊归悧婊堟偉濠婂牊鈷愰柛顭戝枛椤斿﹪骞栫€涙ɑ鐓ｉ悹?43闂佹寧绋戦悧鍡涚嵁閹捐绠冲璺侯槺閸熷﹪鎮樿箛娑氱暫闁?= 婵犻潧顦介崑鍕储閺嶎厽鏅?	if len(state) != 43 {
		t.Errorf("GenerateState 闁哄鏅滈弻銊ッ洪弽顓炵９濠靛倸鎲″В鎰板箹鏉堝墽绱扮紒妤€顦靛畷鐘恒亹閹烘垵璧? got %d, want 43", len(state))
	}
}

func TestGenerateState_闂佸摜鍎愰崹顖滅博閹绢喖绠?t *testing.T) {
	s1, _ := GenerateState()
	s2, _ := GenerateState()
	if s1 == s2 {
		t.Error("婵炴垶鎸堕崐鏍敃?GenerateState 缂傚倷鐒﹂幐濠氭倶婢舵劖鍎庣紒瀣閸?)
	}
}

// ---------------------------------------------------------------------------
// GenerateSessionID
// ---------------------------------------------------------------------------

func TestGenerateSessionID_闁哄鏅滈弻銊ッ洪弽顓炵９闁稿繗鍋愭竟鎰偓?t *testing.T) {
	id, err := GenerateSessionID()
	if err != nil {
		t.Fatalf("GenerateSessionID 婵犮垺鍎肩划鍓ф喆? %v", err)
	}
	if id == "" {
		t.Error("GenerateSessionID 闁哄鏅滈弻銊ッ洪弽顐ょ煔闁告繂瀚幗鐔虹磼濡ゅ绱伴悷?)
	}
	// 16 闁诲孩绋掗〃澶嬩繆椤撱垺鍎?hex 缂傚倸鍊归悧婊堟偉濠婂牊鈷愰柛顭戝枛椤斿﹪骞栫€涙ɑ鐓ｉ悹?32
	if len(id) != 32 {
		t.Errorf("GenerateSessionID 闁哄鏅滈弻銊ッ洪弽顓炵９濠靛倸鎲″В鎰板箹鏉堝墽绱扮紒妤€顦靛畷鐘恒亹閹烘垵璧? got %d, want 32", len(id))
	}
	// 婵°倗濮撮惌渚€鎯佹径鎰強妞ゆ牗纰嶉崐銈嗙箾婢跺绀堟繛?hex 闁诲孩绋掗〃鍫ヮ敄娴ｅ湱鈻?	if _, err := hex.DecodeString(id); err != nil {
		t.Errorf("GenerateSessionID 闁哄鏅滈弻銊ッ洪弽顓炵９闂傚倸顕悷婵嬫煛閸曢潧鐏犻柟顔奸閳绘棃寮撮悩宕囨殸 hex 闁诲孩绋掗〃鍫ヮ敄娴ｅ湱鈻? %s, err: %v", id, err)
	}
}

func TestGenerateSessionID_闂佸摜鍎愰崹顖滅博閹绢喖绠?t *testing.T) {
	id1, _ := GenerateSessionID()
	id2, _ := GenerateSessionID()
	if id1 == id2 {
		t.Error("婵炴垶鎸堕崐鏍敃?GenerateSessionID 缂傚倷鐒﹂幐濠氭倶婢舵劖鍎庣紒瀣閸?)
	}
}

// ---------------------------------------------------------------------------
// GenerateCodeVerifier
// ---------------------------------------------------------------------------

func TestGenerateCodeVerifier_闁哄鏅滈弻銊ッ洪弽顓炵９闁稿繗鍋愭竟鎰偓?t *testing.T) {
	verifier, err := GenerateCodeVerifier()
	if err != nil {
		t.Fatalf("GenerateCodeVerifier 婵犮垺鍎肩划鍓ф喆? %v", err)
	}
	if verifier == "" {
		t.Error("GenerateCodeVerifier 闁哄鏅滈弻銊ッ洪弽顐ょ煔闁告繂瀚幗鐔虹磼濡ゅ绱伴悷?)
	}
	// base64url 缂傚倸鍊归悧婊堟偉濠婂嫮鈻旂€广儱鎳愮€瑰鏌涢弽褎鍣归柟?+, /, =
	if strings.ContainsAny(verifier, "+/=") {
		t.Errorf("GenerateCodeVerifier 闁哄鏅滈弻銊ッ洪弽顓炵９闁绘挸楠搁惁鍫曟煕濮樼厧浜版繛?base64url 闁诲孩绋掗〃鍫ヮ敄? %s", verifier)
	}
	// 32 闁诲孩绋掗〃澶嬩繆椤撱垺鍎?base64url 缂傚倸鍊归悧婊堟偉濠婂牊鈷愰柛顭戝枛椤斿﹪骞栫€涙ɑ鐓ｉ悹?43
	if len(verifier) != 43 {
		t.Errorf("GenerateCodeVerifier 闁哄鏅滈弻銊ッ洪弽顓炵９濠靛倸鎲″В鎰板箹鏉堝墽绱扮紒妤€顦靛畷鐘恒亹閹烘垵璧? got %d, want 43", len(verifier))
	}
}

func TestGenerateCodeVerifier_闂佸摜鍎愰崹顖滅博閹绢喖绠?t *testing.T) {
	v1, _ := GenerateCodeVerifier()
	v2, _ := GenerateCodeVerifier()
	if v1 == v2 {
		t.Error("婵炴垶鎸堕崐鏍敃?GenerateCodeVerifier 缂傚倷鐒﹂幐濠氭倶婢舵劖鍎庣紒瀣閸?)
	}
}

// ---------------------------------------------------------------------------
// GenerateCodeChallenge
// ---------------------------------------------------------------------------

func TestGenerateCodeChallenge_SHA256_Base64URL(t *testing.T) {
	verifier := "dBjftJeZ4CVP-mB92K27uhbUJU1p1r_wW1gFWFOEjXk"

	challenge := GenerateCodeChallenge(verifier)

	// 闂佸綊娼ч鍛叏閳哄啯濯奸柨娑樺閺嗩剙螞閺夊灝顏繝鈧敓鐘茬９?	hash := sha256.Sum256([]byte(verifier))
	expected := strings.TrimRight(base64.URLEncoding.EncodeToString(hash[:]), "=")

	if challenge != expected {
		t.Errorf("CodeChallenge 婵炴垶鎸哥粔瀵镐焊椤曗偓閺? got %s, want %s", challenge, expected)
	}
}

func TestGenerateCodeChallenge_婵炴垶鎸哥粔鎾箚閸喓绠欐い鎰╁灩鐢娀鎮楀☉娆樻畼妞?t *testing.T) {
	challenge := GenerateCodeChallenge("test-verifier")
	if strings.Contains(challenge, "=") {
		t.Errorf("CodeChallenge 婵炴垶鎸哥粔瀵歌姳閺屻儱绀岄柛娑卞幗閸?= 婵犻潧顦介崑鍕储閺嶎偀鍋撳☉娆樻畼妞? %s", challenge)
	}
}

func TestGenerateCodeChallenge_婵炴垶鎸哥粔鎾箚閸儲顥堟繝鐢电ゲL闁诲海鎳撻ˇ顖炲矗韫囨洍鍋撳☉娆樻畼妞?t *testing.T) {
	challenge := GenerateCodeChallenge("another-verifier")
	if strings.ContainsAny(challenge, "+/") {
		t.Errorf("CodeChallenge 婵炴垶鎸哥粔瀵歌姳閺屻儱绀岄柛娑卞幗閸?+ 闂?/ 闁诲孩绋掗〃鍫ヮ敄? %s", challenge)
	}
}

func TestGenerateCodeChallenge_闂佺儵鏅濋…鍫ュ箖閹惧瓨缍囬柟鎯у暱瀵娊鏌ｉ埡鍐剧劸闁诡喗鎸诲蹇涘箻閸愬弶鐦?t *testing.T) {
	c1 := GenerateCodeChallenge("same-verifier")
	c2 := GenerateCodeChallenge("same-verifier")
	if c1 != c2 {
		t.Errorf("闂佺儵鏅濋…鍫ュ箖閹惧瓨缍囬柟鎯у暱瀵娊骞栫€涙ɑ鐓ｅΔ鐘叉喘閹粙鎮㈤悜妯笺偛闂佸憡鑹鹃惌浣烘椤撱垹绀? got %s and %s", c1, c2)
	}
}

func TestGenerateCodeChallenge_婵炴垶鎸哥粔鎾箖閹惧瓨缍囬柟鎯у暱瀵啿鈽夐幘宕囆㈤柟顔芥尰濞煎繘骞橀崘鍙夌様(t *testing.T) {
	c1 := GenerateCodeChallenge("verifier-1")
	c2 := GenerateCodeChallenge("verifier-2")
	if c1 == c2 {
		t.Error("婵炴垶鎸哥粔鎾箖閹惧瓨缍囬柟鎯у暱瀵娊骞栫€涙ɑ鐓ｅΔ鐘叉喘閹粙鎮㈤崜浣烘喒闂佸憡鑹鹃惌浣烘椤撱垹绀?)
	}
}

// ---------------------------------------------------------------------------
// BuildAuthorizationURL
// ---------------------------------------------------------------------------

func TestBuildAuthorizationURL_闂佸憡鐟ラ崐褰掑汲閻斿壊娈界€光偓閸愵亝顫?t *testing.T) {
	state := "test-state-123"
	codeChallenge := "test-challenge-abc"

	authURL := BuildAuthorizationURL(state, codeChallenge)

	// 婵°倗濮撮惌渚€鎯佹径瀣浄?AuthorizeURL 閻庢鍠掗崑鎾愁熆?	if !strings.HasPrefix(authURL, AuthorizeURL+"?") {
		t.Errorf("URL 闁圭厧鐡ㄩ弻褎绂?%s? 閻庢鍠掗崑鎾愁熆? got %s", AuthorizeURL, authURL)
	}

	// 闁荤喐鐟辩徊楣冩倵?URL 濡ょ姷鍋涢悥濂告偘濞嗘垶瀚氬ù锝囶焾濡﹢鏌?	parsed, err := url.Parse(authURL)
	if err != nil {
		t.Fatalf("闁荤喐鐟辩徊楣冩倵?URL 婵犮垺鍎肩划鍓ф喆? %v", err)
	}

	params := parsed.Query()

	expectedParams := map[string]string{
		"client_id":              ClientID,
		"redirect_uri":           RedirectURI,
		"response_type":          "code",
		"scope":                  Scopes,
		"state":                  state,
		"code_challenge":         codeChallenge,
		"code_challenge_method":  "S256",
		"access_type":            "offline",
		"prompt":                 "consent",
		"include_granted_scopes": "true",
	}

	for key, want := range expectedParams {
		got := params.Get(key)
		if got != want {
			t.Errorf("闂佸憡鐟ラ崐褰掑汲?%s 婵炴垶鎸哥粔瀵镐焊椤曗偓閺? got %q, want %q", key, got, want)
		}
	}
}

func TestBuildAuthorizationURL_闂佸憡鐟ラ崐褰掑汲閻旂厧鏋佸ù鍏兼綑濞?t *testing.T) {
	authURL := BuildAuthorizationURL("s", "c")
	parsed, err := url.Parse(authURL)
	if err != nil {
		t.Fatalf("闁荤喐鐟辩徊楣冩倵?URL 婵犮垺鍎肩划鍓ф喆? %v", err)
	}

	params := parsed.Query()
	// 闁圭厧鐡ㄩ弻銊р偓鍨耿瀹?10 婵炴垶鎼╂禍婊冾嚕椤掑嫬鏋?	expectedCount := 10
	if len(params) != expectedCount {
		t.Errorf("闂佸憡鐟ラ崐褰掑汲閻旂厧鏋佸ù鍏兼綑濞呫倕鈽夐幘宕囆㈤悘蹇ｅ櫍閺? got %d, want %d", len(params), expectedCount)
	}
}

func TestBuildAuthorizationURL_闂佺粯顨夐～澶愭偩閳哄啠鍋撳☉娆樻畼妞ゆ垳鑳剁槐鎾诲冀椤愩倕鐏?t *testing.T) {
	state := "state+with/special=chars"
	codeChallenge := "challenge+value"

	authURL := BuildAuthorizationURL(state, codeChallenge)

	parsed, err := url.Parse(authURL)
	if err != nil {
		t.Fatalf("闁荤喐鐟辩徊楣冩倵?URL 婵犮垺鍎肩划鍓ф喆? %v", err)
	}

	// 闁荤喐鐟辩徊楣冩倵娴犲瑙﹂幖杈剧磿鐎瑰鎱ㄥ┑鎾舵偧闁炽儲蓱濞艰螣閸濆嫭鏋鹃梺缁橆殙椤鎮鹃埡鍐ｅ亾濞戞瑯娈樻い?
	if got := parsed.Query().Get("state"); got != state {
		t.Errorf("state 闂佸憡鐟ラ崐褰掑汲閻旇櫣纾介柡宥庡亞閸?闁荤喐鐟辩徊鍧楁偉濠婂嫮鈻旂€广儱鎳庨悥閬嶆⒑? got %q, want %q", got, state)
	}
}

// ---------------------------------------------------------------------------
// 闁汇埄鍨遍幃鍌炲闯濞差亜纾瑰┑鍌氭憸瀹曪綁鎮?// ---------------------------------------------------------------------------

func TestConstants_闂佺锕ら崥瀣敆濠婂懏鍏?t *testing.T) {
	if AuthorizeURL != "https://accounts.google.com/o/oauth2/v2/auth" {
		t.Errorf("AuthorizeURL 婵炴垶鎸哥粔瀵镐焊椤曗偓閺? got %s", AuthorizeURL)
	}
	if TokenURL != "https://oauth2.googleapis.com/token" {
		t.Errorf("TokenURL 婵炴垶鎸哥粔瀵镐焊椤曗偓閺? got %s", TokenURL)
	}
	if UserInfoURL != "https://www.googleapis.com/oauth2/v2/userinfo" {
		t.Errorf("UserInfoURL 婵炴垶鎸哥粔瀵镐焊椤曗偓閺? got %s", UserInfoURL)
	}
	if ClientID != "ANTIGRAVITY_OAUTH_CLIENT_ID" {
		t.Errorf("ClientID 婵炴垶鎸哥粔瀵镐焊椤曗偓閺? got %s", ClientID)
	}
    secret, err := getClientSecret()
    if err == nil {
        t.Fatalf("getClientSecret should return an error when no default client secret is configured")
    }
    if secret != "" {
        t.Errorf("default client_secret should be empty, got %s", secret)
    }
    if GetUserAgent() != "antigravity/1.21.9 windows/amd64" {
		t.Errorf("UserAgent 婵炴垶鎸哥粔瀵镐焊椤曗偓閺? got %s", GetUserAgent())
	}
	if SessionTTL != 30*time.Minute {
		t.Errorf("SessionTTL 婵炴垶鎸哥粔瀵镐焊椤曗偓閺? got %v", SessionTTL)
	}
	if URLAvailabilityTTL != 5*time.Minute {
		t.Errorf("URLAvailabilityTTL 婵炴垶鎸哥粔瀵镐焊椤曗偓閺? got %v", URLAvailabilityTTL)
	}
}

func TestScopes_闂佸憡鐗曢幊搴ㄥ箚閸絿鐤€闁告盯鍋婂ú锝夋煠閻撳骸鏆欐繛?t *testing.T) {
	expectedScopes := []string{
		"https://www.googleapis.com/auth/cloud-platform",
		"https://www.googleapis.com/auth/userinfo.email",
		"https://www.googleapis.com/auth/userinfo.profile",
		"https://www.googleapis.com/auth/cclog",
		"https://www.googleapis.com/auth/experimentsandconfigs",
	}

	for _, scope := range expectedScopes {
		if !strings.Contains(Scopes, scope) {
			t.Errorf("Scopes 缂傚倸鍊搁幖顐︽儍?%s", scope)
		}
	}
}

