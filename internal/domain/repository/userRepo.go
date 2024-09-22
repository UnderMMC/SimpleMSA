package repository

import (
	"database/sql"
	"github.com/SimpleMSA/internal/domain/entity"
)

type PostgresUserRepository struct {
	db *sql.DB
}

func NewPostgresUserRepository(db *sql.DB) *PostgresUserRepository {
	return &PostgresUserRepository{db: db}
}

func (r *PostgresUserRepository) UserRegistration(user entity.User) error {
	_, err := r.db.Exec("INSERT INTO logdata (login, password) VALUES ($1, $2)", user.Login, user.Password)
	if err != nil {
		return err
	}
	return err
}

func (r *PostgresUserRepository) GetUserHashedPass(user entity.User) (string, error) {
	var storedPassword string
	err := r.db.QueryRow("SELECT password FROM logdata WHERE login=$1", user.Login).Scan(&storedPassword)

	return storedPassword, err
}

func (r *PostgresUserRepository) GetUserID(user entity.User) (int, error) {
	var ID int
	err := r.db.QueryRow("SELECT user_id FROM logdata WHERE login=$1", user.Login).Scan(&ID)
	if err != nil {
		return 0, err
	}
	return ID, nil
}
