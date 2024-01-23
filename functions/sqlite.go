package functions

import (
	"database/sql"
	"log"
)

func CreateTable(db *sql.DB) {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS users (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	username TEXT NOT NULL,
	email TEXT NOT NULL,
	password_hash TEXT NOT NULL,
	created_at TIMESTAMP NOT NULL,
	profile_picture BLOB,
	is_admin INTEGER DEFAULT 0 CHECK (is_admin IN (0, 1)),
	is_banned BOOLEAN NOT NULL DEFAULT FALSE
);`)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`
CREATE TABLE IF NOT EXISTS topics (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	username TEXT NOT NULL,
	title TEXT NOT NULL,
	content TEXT NOT NULL,
	created_at TIMESTAMP NOT NULL,
	category INTEGER NOT NULL,
	upvotes TEXT NOT NULL DEFAULT '0',
	image BLOB
);
`)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS messages (
	id TEXT PRIMARY KEY,
	username TEXT NOT NULL,
	contenu TEXT,
	created_at TIMESTAMP NOT NULL,
	valeur TEXT
);`)
	if err != nil {
		log.Fatal(err)
	}
}
