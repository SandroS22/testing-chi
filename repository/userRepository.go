package repository

import (
	"database/sql"
	"errors"

	"github.com/SandroS22/chi-test/db"
)

var (
	ErrUserNotFound = errors.New("user not found")
)

type UserRepository interface {
	Create(user *db.User) (int, error)
	GetByID(id int) (*db.User, error)
	GetAll() ([]db.User, error)
	Update(user *db.User) error
	Delete(id int) error
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(database *sql.DB) UserRepository {
	return &userRepository{db: database}
}

func (r *userRepository) Create(user *db.User) (int, error) {
	return db.InsertUser(r.db, user.Name, user.User, user.Password, user.Role)
}

func (r *userRepository) GetByID(id int) (*db.User, error) {
	user, err := db.GetUserByID(r.db, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return user, nil
}

func (r *userRepository) GetAll() ([]db.User, error) {
	users, err := db.GetAllUsers(r.db)
	if err != nil {
		return nil, err
	}
	if users == nil {
		users = []db.User{}
	}
	return users, nil
}

func (r *userRepository) Update(user *db.User) error {
	_, err := r.GetByID(user.ID)
	if err != nil {
		return err
	}
	return db.UpdateUser(r.db, *user)
}

func (r *userRepository) Delete(id int) error {
	_, err := r.GetByID(id)
	if err != nil {
		return err
	}
	return db.DeleteUser(r.db, id)
}
