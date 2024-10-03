package sqlite

import (
	"database/sql"

	"github.com/SandipKurmi/students-api/internal/config"
	_ "github.com/mattn/go-sqlite3" // SQLite driver
)


type sqlite struct {
	Db *sql.DB
}

func New(cfg *config.Config) (*sqlite, error) {
	// Connect to the SQLite database
	db, err := sql.Open("sqlite3", cfg.StoragePath)
	if err != nil {
		return nil, err
	}


	_ , err = db.Exec(`CREATE TABLE IF NOT EXISTS students (
	id INTEGER PRIMARY KEY AUTOINCREMENT, 
	name TEXT, 
	email TEXT,
	age INTEGER
	)`)

	if err != nil {
		return nil, err
	}

	return &sqlite{
		Db: db,
	}, nil

}  

func (s *sqlite) CreateStudent(name string, email string, age int) (int64, error) {

	stmt, err := s.Db.Prepare("INSERT INTO students(name, email, age) VALUES (?, ?, ?)")
	if err != nil {
		return 0, err
	}

	defer stmt.Close()

	res, err := stmt.Exec(name, email, age)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()

	if err != nil {
		return 0, err
	}

	return id, nil

}