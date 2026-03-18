package auth

import (
	"context"
	"net/http"
	"time"

	"github.com/alexedwards/scs/pgxstore"
	"github.com/alexedwards/scs/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

// SessionUserID is the key used to store the user ID in the session.
const SessionUserID = "user_id"

type contextKey string

const userIDContextKey contextKey = "user_id"

// ContextWithUserID stores a user ID in the context for downstream handlers.
func ContextWithUserID(ctx context.Context, id [16]byte) context.Context {
	return context.WithValue(ctx, userIDContextKey, id)
}

// UserIDFromContext retrieves the user ID stored by OptionalAuth or RequireAuth.
// Returns the zero value and false if not present.
func UserIDFromContext(ctx context.Context) ([16]byte, bool) {
	id, ok := ctx.Value(userIDContextKey).([16]byte)
	return id, ok
}

// NewSessionManager creates an SCS session manager backed by PostgreSQL via
// pgxstore. Session data is stored server-side; the client only receives an
// opaque token in a cookie — no gob registration needed, no data in the cookie.
//
// Requires a `sessions` table — see migrations/002_create_sessions.up.sql.
//
// SameSite is Lax (not Strict) so the session cookie is sent when an OAuth
// provider redirects back to your app (a cross-site top-level navigation).
// Strict would silently drop the cookie on that redirect, breaking auth state.
// The CSRF cookie is separately Strict (see internal/middleware/csrf.go).
func NewSessionManager(pool *pgxpool.Pool) *scs.SessionManager {
	sm := scs.New()
	sm.Store = pgxstore.New(pool)
	sm.Lifetime = 7 * 24 * time.Hour
	sm.Cookie.Name = "session"
	sm.Cookie.HttpOnly = true
	sm.Cookie.SameSite = http.SameSiteLaxMode
	sm.Cookie.Secure = false // set true in production (HTTPS)
	return sm
}
