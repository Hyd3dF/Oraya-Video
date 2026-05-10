package repository

import "github.com/oroya/backend/internal/supabase"

// eq builds a PostgREST equality filter, e.g. eq("id", uuid) -> "id=eq.<uuid>".
func eq(col, val string) (string, string) {
	return col, "eq." + val
}

// idFilter is shorthand for {id=eq.<id>}.
func idFilter(id string) supabase.Filters {
	f := supabase.NewFilters()
	f.Set("id", "eq."+id)
	return f
}
