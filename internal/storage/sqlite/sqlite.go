package sqlite

import (
	"database/sql"

	"github.com/SandipKurmi/students-api/internal/config"
	"github.com/SandipKurmi/students-api/internal/types"
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

// GetStudentById

func (s *sqlite) GetStudentById(id int64) (types.Student, error) {
	// Get student by ID
	row := s.Db.QueryRow("SELECT id, name, email, age FROM students WHERE id = ?", id)

	// Create a new student
	var student types.Student

	// Scan the row into the student
	err := row.Scan(&student.ID, &student.Name, &student.Email, &student.Age)

	if err != nil {
		return types.Student{}, err
	}

	return student, nil
}

func (s *sqlite) GetStudents() ([]types.Student, error) {
	// Get all students
	rows, err := s.Db.Query("SELECT id, name, email, age FROM students")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Create a slice of students
	var students []types.Student

	// Loop through the rows and add each student to the slice
	for rows.Next() {
		var student types.Student
		if err := rows.Scan(&student.ID, &student.Name, &student.Email, &student.Age); err != nil {
			return nil, err
		}
		students = append(students, student)
	}

	return students, nil
}

func  (s *sqlite) UpdateStudent(id int64, name string, email string, age int) (bool, error) {
	// Update student by ID
	_, err := s.Db.Exec("UPDATE students SET name = ?, email = ?, age = ? WHERE id = ?", name, email, age, id)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (s *sqlite) DeleteStudent(id int64) (bool, error) {
	// Delete student by ID
	_, err := s.Db.Exec("DELETE FROM students WHERE id = ?", id)
	if err != nil {
		return false, err
	}
	return true, nil
}