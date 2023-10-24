package migrate_test

import (
	"os"
	"testing"

	"github.com/katallaxie/pkg/migrate"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func withSqlite(t *testing.T, fn func(t *testing.T, db *sqlx.DB)) func(t *testing.T) {
	_ = os.Remove("__deleteme.db")

	db, err := sqlx.Connect("sqlite3", "__deleteme.db?mode=memory")
	if err != nil {
		t.Fatalf("Connect() err = %v; want nil", err)
	}

	t.Cleanup(func() {
		err := db.Close()
		if err != nil {
			t.Errorf("Close() err = %v; want nil", err)
		}
	})

	return func(t *testing.T) {
		fn(t, db)
	}
}

func TestMirgate(t *testing.T) {
	t.Run("simple migration", withSqlite(t, func(t *testing.T, db *sqlx.DB) {
		migrator := migrate.New()
		migrator.Add(migrate.NewMigration(createUsersTable, dropUsersTable))

		id, err := migrator.Migrate(db)
		require.NoError(t, err)
		assert.Equal(t, 1, id)
	}))
}

var (
	createUsersTable = `
CREATE TABLE users (
  id serial PRIMARY KEY,
  foo text
);`

	dropUsersTable = `DROP TABLE users;`
)
