package middleware

import (
	"net/http"

	"github.com/alexedwards/scs/v2"
	"github.com/justestif/go-web-template/internal/auth"
)

// OptionalAuth loads the session user ID into the request context if a valid
// session exists, but does not block unauthenticated requests. Use this on
// pages that render differently for logged-in vs anonymous visitors (e.g. the
// landing page navbar showing "Sign in" vs an avatar).
//
// SCS's LoadAndSave middleware must wrap the router for session data to be
// available. OptionalAuth reads from the already-loaded session context.
//
// Retrieve the user ID in a handler with:
//
//	id, ok := auth.UserIDFromContext(r.Context())
func OptionalAuth(sm *scs.SessionManager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			idBytes := sm.GetBytes(r.Context(), auth.SessionUserID)
			if len(idBytes) != 16 {
				next.ServeHTTP(w, r)
				return
			}

			var id [16]byte
			copy(id[:], idBytes)
			ctx := auth.ContextWithUserID(r.Context(), id)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// RequireAuth rejects unauthenticated requests with a redirect to /login.
// Place this on any route group that requires a logged-in user.
//
// SCS's LoadAndSave middleware must wrap the router for session data to be
// available. RequireAuth reads from the already-loaded session context.
func RequireAuth(sm *scs.SessionManager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			idBytes := sm.GetBytes(r.Context(), auth.SessionUserID)
			if len(idBytes) != 16 {
				http.Redirect(w, r, "/login", http.StatusSeeOther)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
