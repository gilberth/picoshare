package authentik

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/mtlynch/picoshare/v2/handlers/auth"
)

const (
	stateCookieName = "authentik_state"
	tokenCookieName = "authentik_token"
)

// Config contains configuration for Authentik authentication
type Config struct {
	AuthentikURL   string
	ClientID       string
	ClientSecret   string
	RedirectURI    string
	TrustedProxies []string // For reverse proxy mode
	ReverseProxy   bool     // Use reverse proxy mode instead of OAuth2
}

// Provider implements AuthProvider for Authentik
type Provider struct {
	config Config
}

// New creates a new Authentik authentication provider
func New(config Config) *Provider {
	return &Provider{
		config: config,
	}
}

// Authenticate implements auth.AuthProvider
func (p *Provider) Authenticate(r *http.Request) (bool, *auth.UserInfo, error) {
	if p.config.ReverseProxy {
		return p.authenticateReverseProxy(r)
	}
	return p.authenticateOAuth2(r)
}

// authenticateReverseProxy validates authentication headers from reverse proxy
func (p *Provider) authenticateReverseProxy(r *http.Request) (bool, *auth.UserInfo, error) {
	// Check if request comes from trusted proxy
	if !p.isTrustedProxy(r) {
		return false, nil, fmt.Errorf("request not from trusted proxy")
	}

	// Get user info from headers
	username := r.Header.Get("X-Authentik-Username")
	email := r.Header.Get("X-Authentik-Email")
	groupsHeader := r.Header.Get("X-Authentik-Groups")

	if username == "" {
		return false, nil, nil
	}

	// Parse groups
	var groups []string
	if groupsHeader != "" {
		groups = strings.Split(groupsHeader, ",")
		for i, group := range groups {
			groups[i] = strings.TrimSpace(group)
		}
	}

	userInfo := &auth.UserInfo{
		Username: username,
		Email:    email,
		Groups:   groups,
	}

	return true, userInfo, nil
}

// authenticateOAuth2 validates OAuth2 tokens
func (p *Provider) authenticateOAuth2(r *http.Request) (bool, *auth.UserInfo, error) {
	// Get token from cookie
	tokenCookie, err := r.Cookie(tokenCookieName)
	if err != nil {
		fmt.Printf("Debug: No token cookie found: %v\n", err)
		return false, nil, nil
	}
	
	fmt.Printf("Debug: Found token cookie, validating...\n")

	// Validate token with Authentik
	userInfo, err := p.validateToken(tokenCookie.Value)
	if err != nil {
		fmt.Printf("Debug: Token validation failed: %v\n", err)
		return false, nil, err
	}
	
	fmt.Printf("Debug: Token validation successful for user: %s\n", userInfo.Username)
	return true, userInfo, nil
}

// StartSession implements auth.AuthProvider
func (p *Provider) StartSession(w http.ResponseWriter, r *http.Request) {
	if p.config.ReverseProxy {
		// In reverse proxy mode, authentication is handled by the proxy
		http.Error(w, "Authentication handled by reverse proxy", http.StatusBadRequest)
		return
	}

	// Handle OAuth2 callback
	if r.URL.Query().Get("code") != "" {
		p.handleCallback(w, r)
		return
	}

	// Generate state parameter
	state, err := p.generateState()
	if err != nil {
		http.Error(w, "Failed to generate state", http.StatusInternalServerError)
		return
	}

	// Store state in cookie
	fmt.Printf("Debug: Setting state cookie: %s\n", state)
	http.SetCookie(w, &http.Cookie{
		Name:     stateCookieName,
		Value:    state,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode, // Changed from Strict to Lax for OAuth callback
		MaxAge:   600, // 10 minutes
	})

	// Build auth URL with the generated state
	params := url.Values{
		"response_type": {"code"},
		"client_id":     {p.config.ClientID},
		"redirect_uri":  {p.config.RedirectURI},
		"scope":         {"openid profile email"},
		"state":         {state},
	}

	authURL := fmt.Sprintf("%s/application/o/authorize/?%s", 
		strings.TrimRight(p.config.AuthentikURL, "/"), 
		params.Encode())

	// Redirect to Authentik
	http.Redirect(w, r, authURL, http.StatusTemporaryRedirect)
}

// ClearSession implements auth.AuthProvider
func (p *Provider) ClearSession(w http.ResponseWriter) {
	// Clear authentication cookies
	http.SetCookie(w, &http.Cookie{
		Name:     tokenCookieName,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   -1,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     stateCookieName,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   -1,
	})
}

// GetLoginURL implements auth.AuthProvider
func (p *Provider) GetLoginURL(returnURL string) string {
	if p.config.ReverseProxy {
		return "/auth"
	}

	// Generate state parameter for security
	state, err := p.generateState()
	if err != nil {
		// Fallback to return URL as state
		state = returnURL
	}

	params := url.Values{
		"response_type": {"code"},
		"client_id":     {p.config.ClientID},
		"redirect_uri":  {p.config.RedirectURI},
		"scope":         {"openid profile email"},
		"state":         {state},
	}

	return fmt.Sprintf("%s/application/o/authorize/?%s", 
		strings.TrimRight(p.config.AuthentikURL, "/"), 
		params.Encode())
}

// handleCallback processes OAuth2 callback
func (p *Provider) handleCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	state := r.URL.Query().Get("state")

	// Debug: List all cookies
	fmt.Printf("Debug: All cookies in request:\n")
	for _, cookie := range r.Cookies() {
		fmt.Printf("  Cookie: %s = %s\n", cookie.Name, cookie.Value)
	}
	
	// Verify state parameter
	stateCookie, err := r.Cookie(stateCookieName)
	if err != nil {
		fmt.Printf("Debug: No state cookie found - error: %v, looking for: %s\n", err, stateCookieName)
		http.Error(w, "Invalid state parameter - no cookie", http.StatusBadRequest)
		return
	} 
	
	if stateCookie.Value != state {
		fmt.Printf("Debug: State mismatch - cookie: '%s', param: '%s'\n", stateCookie.Value, state)
		http.Error(w, "Invalid state parameter - mismatch", http.StatusBadRequest)
		return
	}
	
	fmt.Printf("Debug: State validation successful: %s\n", state)

	// Exchange code for token
	token, err := p.exchangeCodeForToken(code)
	if err != nil {
		fmt.Printf("Debug: Token exchange failed: %v\n", err)
		http.Error(w, "Failed to exchange code for token", http.StatusInternalServerError)
		return
	}
	fmt.Printf("Debug: Token exchange successful\n")

	// Store token in cookie
	fmt.Printf("Debug: Setting token cookie: %s\n", tokenCookieName)
	http.SetCookie(w, &http.Cookie{
		Name:     tokenCookieName,
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode, // Changed from Strict to Lax for OAuth callback compatibility
		MaxAge:   3600, // 1 hour
	})

	// Redirect to home page after successful authentication
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}

// generateState generates a random state parameter
func (p *Provider) generateState() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// isTrustedProxy checks if the request comes from a trusted proxy
func (p *Provider) isTrustedProxy(r *http.Request) bool {
	if len(p.config.TrustedProxies) == 0 {
		return true // If no trusted proxies configured, trust all
	}

	clientIP := r.RemoteAddr
	if forwardedFor := r.Header.Get("X-Forwarded-For"); forwardedFor != "" {
		clientIP = strings.Split(forwardedFor, ",")[0]
		clientIP = strings.TrimSpace(clientIP)
	}

	for _, trustedIP := range p.config.TrustedProxies {
		if clientIP == trustedIP {
			return true
		}
	}

	return false
}

// exchangeCodeForToken exchanges authorization code for access token
func (p *Provider) exchangeCodeForToken(code string) (string, error) {
	data := url.Values{
		"grant_type":   {"authorization_code"},
		"code":         {code},
		"redirect_uri": {p.config.RedirectURI},
		"client_id":    {p.config.ClientID},
		"client_secret": {p.config.ClientSecret},
	}

	tokenURL := fmt.Sprintf("%s/application/o/token/", 
		strings.TrimRight(p.config.AuthentikURL, "/"))

	fmt.Printf("Debug: Requesting token from: %s\n", tokenURL)
	resp, err := http.PostForm(tokenURL, data)
	if err != nil {
		fmt.Printf("Debug: HTTP request failed: %v\n", err)
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		fmt.Printf("Debug: Token exchange failed with status %d, body: %s\n", resp.StatusCode, string(body))
		return "", fmt.Errorf("token exchange failed with status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var tokenResponse struct {
		AccessToken string `json:"access_token"`
	}

	if err := json.Unmarshal(body, &tokenResponse); err != nil {
		return "", err
	}

	return tokenResponse.AccessToken, nil
}

// validateToken validates access token and returns user info
func (p *Provider) validateToken(token string) (*auth.UserInfo, error) {
	userInfoURL := fmt.Sprintf("%s/application/o/userinfo/", 
		strings.TrimRight(p.config.AuthentikURL, "/"))

	req, err := http.NewRequest("GET", userInfoURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("token validation failed with status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var userInfoResponse struct {
		Username string   `json:"preferred_username"`
		Email    string   `json:"email"`
		Groups   []string `json:"groups"`
	}

	if err := json.Unmarshal(body, &userInfoResponse); err != nil {
		return nil, err
	}

	return &auth.UserInfo{
		Username: userInfoResponse.Username,
		Email:    userInfoResponse.Email,
		Groups:   userInfoResponse.Groups,
	}, nil
}