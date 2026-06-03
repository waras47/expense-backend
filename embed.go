package expense_backend

import (
	"embed"
	"io/fs"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

var EmbedMigrations fs.FS

func init() {
	var err error
	// fs.Sub akan membuat folder "migrations" menjadi root (.) bagi variabel EmbedMigrations
	EmbedMigrations, err = fs.Sub(embedMigrations, "migrations")
	if err != nil {
		panic(err)
	}
}