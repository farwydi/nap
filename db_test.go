package nap

import (
	"database/sql"
	"testing"
	"testing/quick"

	_ "github.com/mattn/go-sqlite3"
)

func TestOpen(t *testing.T) {
	// https://www.sqlite.org/inmemorydb.html
	db, err := Open("sqlite3", ":memory:;:memory:;:memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	if err = db.Ping(); err != nil {
		t.Error(err)
	}

	if want, got := 3, len(db.pdbs); want != got {
		t.Errorf("Unexpected number of physical dbs. Got: %d, Want: %d", got, want)
	}
}

func TestClose(t *testing.T) {
	db, err := Open("sqlite3", ":memory:;:memory:;:memory:")
	if err != nil {
		t.Fatal(err)
	}

	if err = db.Close(); err != nil {
		t.Fatal(err)
	}

	if err = db.Ping(); err.Error() != "sql: database is closed" {
		t.Errorf("Physical dbs were not closed correctly. Got: %s", err)
	}
}

func TestWrap(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	wdb, err := Wrap(db)
	if err != nil {
		t.Fatal(err)
	}

	if err := wdb.Ping(); err != nil {
		t.Error(err)
	}

	if want, got := 1, len(wdb.pdbs); want != got {
		t.Errorf("Unexpected number of physical dbs. Got: %d, Want: %d", got, want)
	}
}

func TestSlave(t *testing.T) {
	db := &DB{}
	last := -1

	err := quick.Check(func(n int) bool {
		index := db.slave(n)
		if n <= 1 {
			return index == 0
		}

		result := index > 0 && index < n && index != last
		last = index

		return result
	}, nil)

	if err != nil {
		t.Error(err)
	}
}
