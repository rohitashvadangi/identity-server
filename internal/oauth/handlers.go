package oauth

import (
	"encoding/json"
	"net/http"
	"time"
	"crypto/rsa"
    "crypto/x509"
    "encoding/pem"
    "io/ioutil"
	"fmt"
    "github.com/golang-jwt/jwt/v5"

	"strings"

	"github.com/google/uuid"
	"github.com/rohitashvadangi/identity-server/internal/proto"
	"github.com/rohitashvadangi/identity-server/internal/stores"
)

// In-memory storage (demo)

type OauthHandler struct {
	authCodeStore stores.AuthCodeStore
	tokenStore    stores.TokenStore
}

var PrivateKey *rsa.PrivateKey
var Issuer = "https://my-idp-server"

func init() {
    keyBytes, err := ioutil.ReadFile("../keys/private.pem")
    if err != nil {
         fmt.Println("Failed to read private key: %v", err)
    }

    block, _ := pem.Decode(keyBytes)
    if block == nil {
        fmt.Println("Failed to decode PEM block containing private key")
    }


   parsedKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
    if err != nil {
        fmt.Println("Failed to parse PKCS#8 private key: %v", err)
    }

    var ok bool
    PrivateKey, ok = parsedKey.(*rsa.PrivateKey)
    if !ok {
        fmt.Println("Private key is not RSA")
    }


    fmt.Println("Private key loaded successfully")
}

func (oauth OauthHandler) AuthorizeHandler(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	clientID := q.Get("client_id")
	redirectURI := q.Get("redirect_uri")
	scope := q.Get("scope")
	state := q.Get("state")

	if clientID == "" || redirectURI == "" {
		http.Error(w, "Missing client_id or redirect_uri", http.StatusBadRequest)
		return
	}

	// Show login form
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, `<form method="POST" action="/login">
        <input type="hidden" name="client_id" value="%s">
        <input type="hidden" name="redirect_uri" value="%s">
        <input type="hidden" name="scope" value="%s">
        <input type="hidden" name="state" value="%s">
        Username: <input name="username"><br>
        Password: <input type="password" name="password"><br>
        <button type="submit">Login</button>
    </form>`, clientID, redirectURI, scope, state)
}

func (oauth OauthHandler) TokenHandler(w http.ResponseWriter, r *http.Request) {
	grantType := r.FormValue("grant_type")

	switch grantType {
	case "authorization_code":
		oauth.handleAuthCodeGrant(w, r)
	case "refresh_token":
		oauth.handleRefreshTokenGrant(w, r)
	default:
		http.Error(w, "unsupported grant_type", http.StatusBadRequest)
	}
}

func (oauth OauthHandler) handleAuthCodeGrant(w http.ResponseWriter, r *http.Request) {

	code := r.FormValue("code")
	clientInReq := r.FormValue("client_id")
	authCodeInfo, err := oauth.authCodeStore.Get(code)
	requestedScope := r.FormValue("scope") // space-separated
	if err != nil {
		fmt.Printf("Failed to get code from store ", err)
		http.Error(w, "invalid code!!!", http.StatusBadRequest)
		return
	}

	if authCodeInfo == nil {
		http.Error(w, "invalid code!!!", http.StatusBadRequest)
		return
	}

	if authCodeInfo.ClientID != clientInReq {
		fmt.Printf("Failed in client authCode %s and in request %s  ", authCodeInfo.ClientID, clientInReq)
		http.Error(w, "invalid client", http.StatusBadRequest)
		return
	}

	// Parse requested scopes and intersect with allowed scopes
	allowedScopes := map[string]bool{"openid": true, "profile": true}
	var scopes []string
	for _, s := range strings.Fields(requestedScope) {
		if allowedScopes[s] {
			scopes = append(scopes, s)
		}
	}

	oauth.authCodeStore.Delete(code)

	tokenId := "atk-" + uuid.NewString()
	refreshToken := "rtk-" + uuid.NewString()

	// Issue dummy access & ID token
	resp := map[string]interface{}{
		"access_token":  tokenId,
		"token_type":    "Bearer",
		"expires_in":    3600,
		"id_token":      "idtoken123",
		"refresh_token": refreshToken,
	}

	token := &proto.Token{ClientID: authCodeInfo.ClientID, TokenID: tokenId, UserID: authCodeInfo.UserID, Scope: scopes, ExpiresAt: time.Now().Add(30 * time.Minute), RevocationID: "revocation_" + tokenId}
	oauth.tokenStore.Save(token)
	//Refresh grant allowed
	if true {
		refreshToken := &proto.RefreshToken{ID: refreshToken, AccessToken: *token, ExpiresAt: time.Now().Add(100000 * time.Minute)}
		oauth.tokenStore.SaveRefresh(refreshToken)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (oauth OauthHandler) handleRefreshTokenGrant(w http.ResponseWriter, r *http.Request) {
	refresh := r.FormValue("refresh_token")
	oldToken, err := oauth.tokenStore.GetByRefresh(refresh)
	if err != nil {
		http.Error(w, "invalid refresh token", http.StatusBadRequest)
		return
	}

	
	// Generate new access token (reuse refresh token or issue new one)
	tokenId := "atk-" + uuid.NewString()
	newToken := &proto.Token{
		TokenID:      tokenId,
		UserID:       oldToken.AccessToken.UserID,
		ClientID:     oldToken.AccessToken.ClientID,
		Scope:        oldToken.AccessToken.Scope,
		ExpiresAt:    time.Now().Add(30 * time.Minute),
		RevocationID: oldToken.AccessToken.RevocationID,
	}

	oauth.tokenStore.Save(newToken)

	resp := map[string]interface{}{
		"access_token":  tokenId,
		"token_type":    "Bearer",
		"expires_in":    3600,
		"refresh_token": refresh,
		"scope":         strings.Join(newToken.Scope, " "),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (oauth OauthHandler) IntrospectHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	tokenStr := r.FormValue("token")
	if tokenStr == "" {
		http.Error(w, "missing token", http.StatusBadRequest)
		return
	}

	t, err := oauth.tokenStore.Get(tokenStr)
	if err != nil {
		json.NewEncoder(w).Encode(map[string]any{"active": false})
		return
	}

	if time.Now().After(t.ExpiresAt) {
		oauth.tokenStore.Delete(tokenStr)
		json.NewEncoder(w).Encode(map[string]any{"active": false})
		return
	}

	/**resp := map[string]any{
		"active":    true,
		"sub":       t.UserID,
		"client_id": t.ClientID,
		"scope":     strings.Join(t.Scope, " "),
		"exp":       t.ExpiresAt.Unix(),
	}
*/
	 claims := jwt.MapClaims{
        "sub": t.UserID,
        "exp": time.Now().Add(time.Hour * 1).Unix(), // 1 hour expiry
        "iat": time.Now().Unix(),
        "iss": Issuer,
		"client_id": t.ClientID,
		"scope":     strings.Join(t.Scope, " "),
		"active" : true,
    }

	resp,err:=oauth.generateJWT(PrivateKey,claims)
	if err!=nil {
		json.NewEncoder(w).Encode(err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (oauth OauthHandler) RevokeHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	tokenStr := r.FormValue("token")

	if tokenStr != "" {
		oauth.tokenStore.Delete(tokenStr)
	} else {
		http.Error(w, "token or revocation_id required", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func NewOauthHandler(authCodeStore stores.AuthCodeStore, tokenStore stores.TokenStore) *OauthHandler {
	return &OauthHandler{
		authCodeStore: authCodeStore, tokenStore: tokenStore,
	}
}

// GenerateJWT signs a token with RSA private key
func (oauth OauthHandler) generateJWT(privateKey *rsa.PrivateKey, claims jwt.MapClaims) (string, error) {

    token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
    return token.SignedString(privateKey)
}