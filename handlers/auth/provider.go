package auth

import (
	"net/http"
	"time"
)

// UserInfo contains information about an authenticated user
type UserInfo struct {
	Username string
	Email    string
	Groups   []string
}

// AuthProvider defines the interface for authentication providers
type AuthProvider interface {
	// Authenticate verifies if the request has valid authentication
	Authenticate(r *http.Request) (bool, *UserInfo, error)
	
	// StartSession begins an authenticated session
	StartSession(w http.ResponseWriter, r *http.Request)
	
	// ClearSession removes the authentication session
	ClearSession(w http.ResponseWriter)
	
	// GetLoginURL returns the URL to redirect for authentication (optional)
	GetLoginURL(returnURL string) string
}

// MultiProviderAuthenticator manages multiple authentication providers
type MultiProviderAuthenticator struct {
	providers       map[string]AuthProvider
	defaultProvider string
	currentUser     *UserInfo
	lastAuth        time.Time
}

// NewMultiProviderAuthenticator creates a new multi-provider authenticator
func NewMultiProviderAuthenticator(defaultProvider string) *MultiProviderAuthenticator {
	return &MultiProviderAuthenticator{
		providers:       make(map[string]AuthProvider),
		defaultProvider: defaultProvider,
	}
}

// RegisterProvider registers a new authentication provider
func (m *MultiProviderAuthenticator) RegisterProvider(name string, provider AuthProvider) {
	m.providers[name] = provider
}

// Authenticate implements the handlers.Authenticator interface
func (m *MultiProviderAuthenticator) Authenticate(r *http.Request) bool {
	provider := m.getProvider(r)
	if provider == nil {
		return false
	}
	
	authenticated, userInfo, err := provider.Authenticate(r)
	if err != nil || !authenticated {
		return false
	}
	
	// Store user info for later use
	m.currentUser = userInfo
	m.lastAuth = time.Now()
	
	return true
}

// StartSession implements the handlers.Authenticator interface
func (m *MultiProviderAuthenticator) StartSession(w http.ResponseWriter, r *http.Request) {
	provider := m.getProvider(r)
	if provider == nil {
		http.Error(w, "No authentication provider available", http.StatusInternalServerError)
		return
	}
	
	provider.StartSession(w, r)
}

// ClearSession implements the handlers.Authenticator interface
func (m *MultiProviderAuthenticator) ClearSession(w http.ResponseWriter) {
	// Clear session for all providers
	for _, provider := range m.providers {
		provider.ClearSession(w)
	}
	
	// Clear internal state
	m.currentUser = nil
	m.lastAuth = time.Time{}
}

// GetCurrentUser returns information about the currently authenticated user
func (m *MultiProviderAuthenticator) GetCurrentUser() *UserInfo {
	return m.currentUser
}

// GetLoginURL returns the login URL for the current provider
func (m *MultiProviderAuthenticator) GetLoginURL(returnURL string) string {
	provider := m.getProvider(nil)
	if provider == nil {
		return "/login"
	}
	return provider.GetLoginURL(returnURL)
}

// getProvider returns the appropriate provider for the request
func (m *MultiProviderAuthenticator) getProvider(r *http.Request) AuthProvider {
	// For now, use the default provider
	// In the future, we could determine the provider based on headers, path, etc.
	if provider, exists := m.providers[m.defaultProvider]; exists {
		return provider
	}
	
	// Fallback to any available provider
	for _, provider := range m.providers {
		return provider
	}
	
	return nil
}