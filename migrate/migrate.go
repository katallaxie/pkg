package migrate

import (
	"database/sql"
	"errors"
	"fmt"
	"math"
	"strconv"
	"sync"

	"github.com/jmoiron/sqlx"
)

// SqlxMigrate ...
type SqlxMigrate struct {
	Migrations []SqlxMigration

	rw sync.Mutex
}

// New ...
func New() *SqlxMigrate {
	return &SqlxMigrate{}
}

// Add ...
func (s *SqlxMigrate) Add(m SqlxMigration) {
	m.id = len(s.Migrations) + 1
	s.Migrations = append(s.Migrations, m)
}

// Step applies 1 migration (up)
func (s *SqlxMigrate) Step(db *sqlx.DB) (int, error) {
	return s.Run(db, 1)
}

// Migrate ...
func (s *SqlxMigrate) Migrate(db *sqlx.DB) (int, error) {
	return s.Run(db, math.MaxInt32)
}

// Run ...
func (s *SqlxMigrate) Run(db *sqlx.DB, steps int) (int, error) {
	s.rw.Lock()
	defer s.rw.Unlock()

	version := -1
	if err := s.createMigrationTable(db); err != nil {
		return -1, err
	}

	for _, m := range s.Migrations {
		found, err := s.selectVersion(db, m.id)
		if err != nil {
			return version, err
		}

		if found {
			continue
		}

		err = s.migrate(db, m)
		if err != nil {
			return version, err
		}

		version = m.id
	}

	return version, nil
}

// Rollback ...
func (s *SqlxMigrate) Rollback(db *sqlx.DB) error {
	s.rw.Lock()
	defer s.rw.Unlock()

	err := s.createMigrationTable(db)
	if err != nil {
		return err
	}

	for _, m := range s.Migrations {
		if m.Down == nil {
			continue
		}

		found, err := s.selectVersion(db, m.id)
		if err != nil {
			return err
		}

		if !found && err != nil {
			continue
		}

		err = s.rollback(db, m)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *SqlxMigrate) rollback(db *sqlx.DB, m SqlxMigration) error {
	errorf := func(err error) error { return fmt.Errorf("running rollback: %w", err) }

	tx, err := db.Beginx()
	if err != nil {
		return errorf(err)
	}

	err = m.Down(tx)
	if err != nil {
		_ = tx.Rollback()
		return errorf(err)
	}

	err = tx.Commit()
	if err != nil {
		return errorf(err)
	}

	return nil
}

func (s *SqlxMigrate) migrate(db *sqlx.DB, m SqlxMigration) error {
	errorf := func(err error) error { return fmt.Errorf("running migration: %w", err) }

	tx, err := db.Beginx()
	if err != nil {
		return errorf(err)
	}

	err = s.insertVersion(db, m.id)
	if err != nil {
		_ = tx.Rollback()
		return errorf(err)
	}

	err = m.Up(tx)
	if err != nil {
		_ = tx.Rollback()
		return errorf(err)
	}

	err = tx.Commit()
	if err != nil {
		return errorf(err)
	}

	return nil
}

// TableName ...
func (s *SqlxMigrate) TableName() string {
	return "migrations"
}

// ColumnName ...
func (s *SqlxMigrate) ColumnName() string {
	return "version"
}

// SqlxVersion ...
type SqlxVersion struct {
	Version int
}

func (s *SqlxMigrate) createMigrationTable(db *sqlx.DB) error {
	_, err := db.Exec(fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s ( %s INTEGER )`, s.TableName(), s.ColumnName()))
	if err != nil {
		return fmt.Errorf("creating migrations table: %w", err)
	}

	return nil
}

func (s *SqlxMigrate) selectVersion(db *sqlx.DB, version int) (bool, error) {
	var row struct {
		Version int
	}

	err := db.Get(&row, fmt.Sprintf(`SELECT %s FROM %s WHERE %s=%s`, s.ColumnName(), s.TableName(), s.ColumnName(), strconv.Itoa(version)))
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}

	if err != nil {
		return false, fmt.Errorf("select current version: %w", err)
	}

	return true, nil
}

func (s *SqlxMigrate) insertVersion(db *sqlx.DB, version int) error {
	_, err := db.Exec(fmt.Sprintf(`INSERT INTO %s ( %s ) VALUES ( %s )`, s.TableName(), s.ColumnName(), strconv.Itoa(version)))
	if err != nil {
		return fmt.Errorf("creating migrations table: %w", err)
	}

	return nil
}

// SqlxMigration ...
type SqlxMigration struct {
	id   int
	Up   func(tx *sqlx.Tx) error
	Down func(tx *sqlx.Tx) error
}

// NewMigration ...
func NewMigration(up, down string) SqlxMigration {
	queryFn := func(query string) func(tx *sqlx.Tx) error {
		if query == "" {
			return nil
		}

		return func(tx *sqlx.Tx) error {
			_, err := tx.Exec(query)
			return err
		}
	}

	m := SqlxMigration{
		Up:   queryFn(up),
		Down: queryFn(down),
	}

	return m
}
