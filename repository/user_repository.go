package repository

import (
	user "github.com/Santosh-Sohan/user-service/api"
	"github.com/Santosh-Sohan/user-service/db"
	"github.com/gocql/gocql"
	"github.com/google/uuid"
)

type UserRepository interface {
	CreateUser(*user.User) error
	UpdateUser(*user.UpdateUserRequest) error
	UpdateContact(*user.UpdateContactRequest) error
	BlockUser(id string) error
	UnblockUser(id string) error
	GetUserByEmailOrPhone(string, string) (*user.User, error)
}

type cassandraUserRepository struct{}

func NewUserRepository() UserRepository {
	return &cassandraUserRepository{}
}

func (r *cassandraUserRepository) CreateUser(u *user.User) error {
	u.Id = uuid.New().String()
	return db.Session.Query(`
		INSERT INTO users (id, first_name, last_name, gender, date_of_birth, phone_number, email, is_blocked)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		u.Id, u.FirstName, u.LastName, u.Gender, u.DateOfBirth, u.PhoneNumber, u.Email, u.IsBlocked,
	).Exec()
}

func (r *cassandraUserRepository) UpdateUser(req *user.UpdateUserRequest) error {
	return db.Session.Query(`
		UPDATE users SET first_name = ?, last_name = ?, gender = ?, date_of_birth = ?
		WHERE id = ?`,
		req.FirstName, req.LastName, req.Gender, req.DateOfBirth, req.Id,
	).Exec()
}

func (r *cassandraUserRepository) UpdateContact(req *user.UpdateContactRequest) error {
	return db.Session.Query(`
		UPDATE users SET phone_number = ?, email = ?
		WHERE id = ?`,
		req.PhoneNumber, req.Email, req.Id,
	).Exec()
}

func (r *cassandraUserRepository) BlockUser(id string) error {
	return db.Session.Query(`UPDATE users SET is_blocked = true WHERE id = ?`, id).Exec()
}

func (r *cassandraUserRepository) UnblockUser(id string) error {
	return db.Session.Query(`UPDATE users SET is_blocked = false WHERE id = ?`, id).Exec()
}

func (r *cassandraUserRepository) GetUserByEmailOrPhone(phone, email string) (*user.User, error) {
	var u user.User
	q := `SELECT id, first_name, last_name, gender, date_of_birth, phone_number, email, is_blocked
	      FROM users WHERE phone_number = ? OR email = ? LIMIT 1`
	err := db.Session.Query(q, phone, email).Consistency(gocql.One).Scan(
		&u.Id, &u.FirstName, &u.LastName, &u.Gender, &u.DateOfBirth, &u.PhoneNumber, &u.Email, &u.IsBlocked,
	)
	if err != nil {
		return nil, err
	}
	return &u, nil
}
