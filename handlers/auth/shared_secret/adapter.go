package shared_secret

import (
	"net/http"

	"github.com/mtlynch/picoshare/v2/handlers/auth"
)

// Adapter wraps SharedSecretAuthenticator to implement auth.AuthProvider
type Adapter struct {
	authenticator SharedSecretAuthenticator
}

// NewAdapter creates a new adapter for SharedSecretAuthenticator
func NewAdapter(sharedSecretKey string) (*Adapter, error) {
	authenticator, err := New(sharedSecretKey)
	if err != nil {
		return nil, err
	}

	return &Adapter{
		authenticator: authenticator,
	}, nil
}

// Authenticate implements auth.AuthProvider
func (a *Adapter) Authenticate(r *http.Request) (bool, *auth.UserInfo, error) {
	authenticated := a.authenticator.Authenticate(r)
	if !authenticated {
		return false, nil, nil
	}

	// For shared secret auth, we don't have user-specific information
	userInfo := &auth.UserInfo{
		Username: "admin", // Default username for shared secret
		Email:    "",      // No email available
		Groups:   []string{"admin"}, // Default admin group
	}

	return true, userInfo, nil
}

// StartSession implements auth.AuthProvider
func (a *Adapter) StartSession(w http.ResponseWriter, r *http.Request) {
	a.authenticator.StartSession(w, r)
}

// ClearSession implements auth.AuthProvider
func (a *Adapter) ClearSession(w http.ResponseWriter) {
	a.authenticator.ClearSession(w)
}

// GetLoginURL implements auth.AuthProvider
func (a *Adapter) GetLoginURL(returnURL string) string {
	// For shared secret auth, return the login endpoint
	return "/auth"
}