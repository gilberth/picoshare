package handlers

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/mtlynch/picoshare/v2/handlers/auth"
)

var contextKeyIsAuthenticated = &contextKey{"is-authenticated"}

func (s Server) authGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("Debug: /api/auth GET called from %s\n", r.RemoteAddr)
		s.authenticator.StartSession(w, r)
	}
}

func (s Server) authPost() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.authenticator.StartSession(w, r)
	}
}

func (s Server) authDelete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.authenticator.ClearSession(w)
	}
}

func (s Server) authCallback() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.authenticator.StartSession(w, r)
	}
}

func (s Server) checkAuthentication(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), contextKeyIsAuthenticated, s.authenticator.Authenticate((r)))
		h.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (s Server) requireAuthentication(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !isAuthenticated((r.Context())) {
			s.authenticator.ClearSession(w)
			
			// Check if this is a multi-provider authenticator that supports login URLs
			if multiAuth, ok := s.authenticator.(*auth.MultiProviderAuthenticator); ok {
				if loginURL := multiAuth.GetLoginURL(r.URL.String()); loginURL != "" && loginURL != "/auth" {
					// For Authentik OAuth2, redirect to the authorization URL
					http.Redirect(w, r, loginURL, http.StatusTemporaryRedirect)
					return
				}
			}
			
			// Check if request accepts HTML (browser request)
			if strings.Contains(r.Header.Get("Accept"), "text/html") {
				// Redirect to login page for browser requests
				http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
				return
			}
			
			// For API requests, return 401
			http.Error(w, "Authentication required", http.StatusUnauthorized)
			return
		}
		h.ServeHTTP(w, r)
	})
}

func isAuthenticated(ctx context.Context) bool {
	val, ok := ctx.Value(contextKeyIsAuthenticated).(bool)
	if !ok {
		return false
	}
	return val
}
