package store

import "testing"

func TestMigrationV19AddsPairingGrantExpiration(t *testing.T) {
	st, err := Open(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()
	rows, err := st.db.Query(`PRAGMA table_info(pairing_sessions)`)
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()
	found := false
	for rows.Next() {
		var cid int
		var name, kind string
		var notNull, primaryKey int
		var defaultValue any
		if err := rows.Scan(&cid, &name, &kind, &notNull, &defaultValue, &primaryKey); err != nil {
			t.Fatal(err)
		}
		found = found || name == "grant_expires_at"
	}
	if !found {
		t.Fatal("pairing_sessions.grant_expires_at is missing")
	}
}
