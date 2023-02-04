package main

import (
	"context"

	"github.com/katallaxie/pkg/migrate"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

// Migrate ...
func Migrate() *migrate.SqlxMigrate {
	m := migrate.New()

	m.Add(createFlipperFeatures())
	m.Add(createFlipperGates())

	return m
}

func createFlipperFeatures() migrate.SqlxMigration {
	up := `
		CREATE TABLE flipper_features (
			id int(11) NOT NULL AUTO_INCREMENT,
			name varchar(255) NOT NULL,
			created_at datetime NOT NULL,
			updated_at datetime NOT NULL,
			description blob,
			PRIMARY KEY (` + "`id`" + `),
			UNIQUE KEY ` + "`index_flipper_features_on_name`" + ` (` + "`name`" + `)
		) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8
	`

	down := `DROP TABLE IF EXISTS flipper_features;`

	return migrate.NewMigration(up, down)
}

func createFlipperGates() migrate.SqlxMigration {
	up := `
		CREATE TABLE flipper_gates (
			id int(11) NOT NULL AUTO_INCREMENT,
			flipper_feature_id int(11) NOT NULL,
			name varchar(255) NOT NULL,
			value varchar(255) NOT NULL,
			created_at datetime NOT NULL,
			updated_at datetime NOT NULL,
			PRIMARY KEY (` + "`id`" + `),
			UNIQUE KEY ` + "`index_flipper_gates_on_flipper_feature_id_and_name_and_value`" + ` (` + "`flipper_feature_id`" + `, ` + "`name`" + `, ` + "`value`" + `)
		) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8
	`

	down := `DROP TABLE IF EXISTS flipper_gates;`

	return migrate.NewMigration(up, down)
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	db, err := sqlx.ConnectContext(ctx, "mysql", "root:example@(localhost:3306)/example")
	if err != nil {
		panic(err)
	}

	m := Migrate()

	_, err = m.Migrate(db)
	if err != nil {
		panic(err)
	}
}
