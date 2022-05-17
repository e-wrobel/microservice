package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

type Storage struct {
	DB *sql.DB
}

func New(dbFile string) (*Storage, error) {
	var db *sql.DB
	// Check if database exists otherwise initialize it
	if _, err := os.Stat(dbFile); errors.Is(err, os.ErrNotExist) {
		db, err = initDatabase(dbFile)
		if err != nil {
			return nil, ErrInitDatabase
		}
		if err = createTables(db); err != nil {
			return nil, ErrCreatingTables
		}
	} else {
		db, err = sql.Open("sqlite3", dbFile)
		if err != nil {
			return nil, err
		}
	}

	return &Storage{
		DB: db,
	}, nil
}

func createTables(db *sql.DB) error {
	createStudentTableSQL := `CREATE TABLE files (
		"identifier" TEXT NOT NULL PRIMARY KEY,		
		"details" TEXT NOT NULL
	  );`

	// Prepare SQL Statement
	statement, err := db.Prepare(createStudentTableSQL)
	if err != nil {
		return err
	}
	_, err = statement.Exec()
	if err != nil {
		return err
	}

	return nil
}

// Shutdown is intended to close database
func (s Storage) Shutdown() error {
	return s.DB.Close()
}

// Create method puts data into database
func (s Storage) Create(identifier string, details []byte) error {
	updateStatement := `INSERT INTO files(identifier, details) VALUES (?, ?)`
	statement, err := s.DB.Prepare(updateStatement)

	if err != nil {
		return fmt.Errorf("unable to prepare statement: %v", err)
	}

	_, err = statement.Exec(identifier, details)
	if err != nil {
		return fmt.Errorf("unable to exec statement: %v", err)
	}

	return nil
}

func (s Storage) Read(obj *PortEntity) error {
	//TODO implement me
	panic("implement me")
}

// Update method puts data into database
func (s Storage) Update(identifier string, details []byte) error {
	insertStatement := `UPDATE files SET identifier = ?, details = ? WHERE identifier = ?)`
	statement, err := s.DB.Prepare(insertStatement)

	if err != nil {
		return fmt.Errorf("unable to prepare statement: %v", err)
	}

	_, err = statement.Exec(identifier, details)
	if err != nil {
		return fmt.Errorf("unable to exec statement: %v", err)
	}

	return nil
}

func (s Storage) Delete(obj *PortEntity) error {
	//TODO implement me
	panic("implement me")
}

type DBImpl interface {
	ShutDown() error
	Create(identifier string, details []byte) error
	Read(obj *PortEntity) error
	Update(identifier string, details []byte) error
	Delete(obj *PortEntity) error
}

// initDatabase creates sqlite database
func initDatabase(dbFile string) (*sql.DB, error) {
	file, err := os.Create(dbFile)
	if err != nil {
		return nil, fmt.Errorf("unable to create database: %v", err)
	}
	file.Close()

	sqliteDatabase, err := sql.Open("sqlite3", dbFile)

	return sqliteDatabase, err
}
