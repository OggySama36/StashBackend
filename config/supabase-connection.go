package config

import (
	supa "github.com/supabase-community/supabase-go"
)

var SupabaseAdmin *supa.Client
var SupabaseEmail *supa.Client
var SupabasePassword *supa.Client

func Init() {
	SupabaseAdmin, _ = supa.NewClient(
		App.SupabaseURL,
		App.SupabaseKey,
		nil,
	)

	SupabaseEmail, _ = supa.NewClient(
		App.SupabaseURL,
		App.SupabaseKey,
		nil,
	)

	SupabasePassword, _ = supa.NewClient(
		App.SupabaseURL,
		App.SupabaseKey,
		nil,
	)

}
