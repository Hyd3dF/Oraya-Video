// Package repository is the only layer that talks to the database. With the
// Supabase-only architecture, every repo holds a *supabase.Client and goes
// through PostgREST (/rest/v1/) — no direct PostgreSQL connection.
package repository

import (
	"context"
	"errors"

	"github.com/oroya/backend/internal/supabase"
)

// ErrNotFound is returned when a single-row lookup matches no rows.
var ErrNotFound = errors.New("not found")

// DB is a thin alias around the Supabase client so existing code that takes
// *DB keeps compiling. All real I/O goes through the Supabase client below.
type DB struct {
	SB *supabase.Client
}

func New(sb *supabase.Client) *DB {
	return &DB{SB: sb}
}

func (d *DB) Ping(ctx context.Context) error {
	if d == nil || d.SB == nil {
		return errors.New("db not initialized")
	}
	return d.SB.PingREST(ctx)
}

func (d *DB) Close() {} // no-op for HTTP client

// translateNotFound maps the supabase ErrNotFound to the repository's.
func translateNotFound(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, supabase.ErrNotFound) {
		return ErrNotFound
	}
	return err
}
