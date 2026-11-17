package provision

import (
	"context"
	"database/sql"
	"errors"
)

// SQLLookup resolves ids against auth.users with the same filter that
// playground_v2.ApplyOnTop uses (email exact, deleted_at IS NULL).
type SQLLookup struct{ DB *sql.DB }

// UserIDByEmail implements UserIDLookup.
func (l SQLLookup) UserIDByEmail(ctx context.Context, email string) (string, bool, error) {
	var id string
	err := l.DB.QueryRowContext(ctx,
		`SELECT id::text FROM auth.users WHERE email = $1 AND deleted_at IS NULL`, email).Scan(&id)
	if errors.Is(err, sql.ErrNoRows) {
		return "", false, nil
	}
	if err != nil {
		return "", false, err
	}
	return id, true, nil
}
