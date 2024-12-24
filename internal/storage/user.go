package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/wavinamayola/user-management/internal/models"
)

const (
	createUserSql = `INSERT INTO users (
		username,
		first_name,
		last_name,
		email,
		age
	) VALUES (?, ?, ?, ?, ?)`
	getUserSQL    = `SELECT * FROM users WHERE id = ?`
	updateUserSQL = `UPDATE users SET username = ?, first_name = ?, last_name = ?, email = ?, age = ? WHERE id = ?`
	deleteUserSQL = `DELETE FROM users WHERE id = ?`
)

var (
	ErrNotFound       = errors.New("row not found")
	ErrNoRowsAffected = errors.New("no rows affected")
)

func (s *Storage) CreateUser(user models.UserRequest) (int, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return 0, fmt.Errorf("failed to start transaction: %+v", err)
	}

	defer func() {
		switch err {
		case nil:
			err = tx.Commit()
		default:
			tx.Rollback()
		}
	}()

	res, err := tx.Exec(createUserSql,
		user.Username,
		user.FirstName,
		user.LastName,
		user.Email,
		user.Age,
	)
	if err != nil {
		log.Printf("error creating user: %+v", err)
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		log.Printf("error retrieving last insert ID: %+v", err)
		return 0, err
	}

	return int(id), nil
}

func (s *Storage) GetUser(id int) (models.User, error) {
	var user models.User

	if err := s.db.QueryRow(getUserSQL, id).Scan(
		&user.ID,
		&user.Username,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Age,
		&user.Created,
		&user.Updated,
	); err != nil {
		if err == sql.ErrNoRows {
			return user, ErrNotFound
		} else {
			return user, fmt.Errorf("failed to retrieve user: %+v", err)
		}
	}
	return user, nil
}

func (s *Storage) UpdateUser(id int, user models.UserRequest) error {
	res, err := s.db.Exec(updateUserSQL,
		user.Username,
		user.FirstName,
		user.LastName,
		user.Email,
		user.Age,
		id,
	)
	if err != nil {
		log.Printf("error updating user: %+v", err)
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		log.Printf("error retrieving rows affected: %+v", err)
		return err
	}

	if rowsAffected == 0 {
		log.Printf("no rows afected")
		return ErrNoRowsAffected
	}

	return nil
}

func (s *Storage) DeleteUser(id int) error {
	res, err := s.db.Exec(deleteUserSQL, id)
	if err != nil {
		log.Printf("error deleting user: %+v", err)
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		log.Printf("error retrieving rows affected: %+v", err)
		return err
	}

	if rowsAffected == 0 {
		log.Printf("no rows afected")
		return ErrNoRowsAffected
	}

	return nil
}
