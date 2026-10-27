package db

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

// User represents a user record in the database.
type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	User     string `json:"user"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

// InitDB opens a SQLite database at the specified path and creates necessary tables.
func InitDB(filepath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", filepath)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	if err := CreateTables(db); err != nil {
		return nil, err
	}

	return db, nil
}

// CreateTables sets up the database schema.
func CreateTables(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		user TEXT NOT NULL,
		password TEXT NOT NULL,
		role TEXT NOT NULL
	);`

	_, err := db.Exec(query)
	return err
}

// InsertUser adds a new user to the database and returns the generated ID.
func InsertUser(db *sql.DB, name, user, password, role string) (int, error) {
	query := `INSERT INTO users (name, user, password, role) VALUES (?, ?, ?, ?)`
	res, err := db.Exec(query, name, user, password, role)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

// GetUserByID retrieves a single user by their ID.
func GetUserByID(db *sql.DB, id int) (*User, error) {
	query := `SELECT id, name, user, role FROM users WHERE id = ?`
	row := db.QueryRow(query, id)

	var u User
	err := row.Scan(&u.ID, &u.Name, &u.User, &u.Role)
	if err != nil {
		return nil, err
	}

	return &u, nil
}

// GetAllUsers retrieves all users from the database.
func GetAllUsers(db *sql.DB) ([]User, error) {
	query := `SELECT id, name, user, password, role FROM users ORDER BY id ASC`
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.ID, &u.Name, &u.User, &u.Password, &u.Role); err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	return users, rows.Err()
}

// UpdateUser updates an existing user record.
func UpdateUser(db *sql.DB, u User) error {
	query := `UPDATE users SET name = ?, user = ?, role = ? WHERE id = ?`
	_, err := db.Exec(query, u.Name, u.User, u.Role, u.ID)
	return err
}

// DeleteUser deletes a user by their ID.
func DeleteUser(db *sql.DB, id int) error {
	query := `DELETE FROM users WHERE id = ?`
	_, err := db.Exec(query, id)
	return err
}
