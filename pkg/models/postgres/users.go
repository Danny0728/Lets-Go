package postgres

import (
	"database/sql"
	"github.com/yash/snippetbox/pkg/models"
	"golang.org/x/crypto/bcrypt"
)

type UserModel struct {
	DB *sql.DB
}

func (m *UserModel) Insert(name, email, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}
	sqlQuery := `INSERT INTO users (name, email, hashed_password, created) VALUES($1,$2,$3,LOCALTIMESTAMP)`
	_, err = m.DB.Exec(sqlQuery, name, email, string(hashedPassword))
	if err != nil {
		return models.ErrDuplicateEmail
	}
	return err
}

func (m *UserModel) Authenticate(email, password string) (int, error) {
	var id int
	var hashedPassword []byte

	row := m.DB.QueryRow("SELECT id, hashed_password FROM users WHERE email=$1", email)
	err := row.Scan(&id, &hashedPassword)
	if err == sql.ErrNoRows {
		return 0, models.ErrInvalidCredentials
	} else if err != nil {
		return 0, err
	}
	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return 0, models.ErrInvalidCredentials
	} else if err != nil {
		return 0, err
	}
	return id, nil
}
func (m *UserModel) Get(id int) (*models.User, error) {
	s := &models.User{}

	sqlQuery := `SELECT id,name,email, created FROM users WHERE id=$1`
	err := m.DB.QueryRow(sqlQuery, id).Scan(&s.ID, &s.Name, &s.Email, &s.Created)

	if err == sql.ErrNoRows {
		return nil, models.ErrNoRecord
	} else if err != nil {
		return nil, err
	}
	return s, nil
}
