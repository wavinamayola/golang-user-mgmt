package storage

import (
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/wavinamayola/user-management/internal/models"
)

const (
	testUsername  = "testusername"
	testFirstName = "testfirstname"
	testLastName  = "testlastname"
	testEmail     = "testemail@gmail.com"
	testAge       = 25
)

func Test_CreateUser_Success(t *testing.T) {
	const expectedID = 1
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	storage := &Storage{db: db}
	user := models.UserRequest{
		Username:  testUsername,
		FirstName: testFirstName,
		LastName:  testLastName,
		Email:     testEmail,
		Age:       testAge,
	}

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO users").
		WithArgs(user.Username, user.FirstName, user.LastName, user.Email, user.Age).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	id, err := storage.CreateUser(user)
	assert.NoError(t, err)
	assert.Equal(t, expectedID, id)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func Test_GetUser_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	storage := &Storage{db: db}

	userID := 1
	expectedUser := models.User{
		ID:        userID,
		Username:  testUsername,
		FirstName: testFirstName,
		LastName:  testLastName,
		Email:     testEmail,
		Age:       testAge,
		Created:   time.Now(),
		Updated:   time.Now(),
	}

	mock.ExpectQuery("SELECT \\* FROM users WHERE id = \\?").
		WithArgs(userID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "username", "first_name", "last_name", "email", "age", "created", "updated"}).
			AddRow(expectedUser.ID, expectedUser.Username, expectedUser.FirstName, expectedUser.LastName, expectedUser.Email, expectedUser.Age, expectedUser.Created, expectedUser.Updated))

	user, err := storage.GetUser(userID)
	assert.NoError(t, err)
	assert.Equal(t, expectedUser, user)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func Test_GetUser_NoRowsFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	storage := &Storage{db: db}

	userID := 1
	mock.ExpectQuery("SELECT \\* FROM users WHERE id = \\?").
		WithArgs(userID).
		WillReturnError(sql.ErrNoRows)

	user, err := storage.GetUser(userID)
	assert.Error(t, err)
	assert.Equal(t, ErrNotFound, err)
	assert.Equal(t, models.User{}, user)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func Test_UpdateUser_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	storage := &Storage{db: db}

	id := 1
	user := models.UserRequest{
		Username:  testUsername,
		FirstName: testFirstName,
		LastName:  testLastName,
		Email:     testEmail,
		Age:       testAge,
	}

	mock.ExpectExec("UPDATE users").
		WithArgs(user.Username, user.FirstName, user.LastName, user.Email, user.Age, id).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = storage.UpdateUser(id, user)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func Test_UpdateUser_NoRowsAffected(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	storage := &Storage{db: db}

	id := 1
	user := models.UserRequest{
		Username:  testUsername,
		FirstName: testFirstName,
		LastName:  testLastName,
		Email:     testEmail,
		Age:       testAge,
	}

	mock.ExpectExec("UPDATE users").
		WithArgs(user.Username, user.FirstName, user.LastName, user.Email, user.Age, id).
		WillReturnResult(sqlmock.NewResult(0, 0))

	err = storage.UpdateUser(id, user)

	assert.Error(t, err)
	assert.Equal(t, ErrNoRowsAffected, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func Test_DeleteUser_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	storage := &Storage{db: db}

	id := 1
	mock.ExpectExec("DELETE FROM users WHERE id = ?").
		WithArgs(id).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = storage.DeleteUser(id)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func Test_DeleteUser_NoRowsAffected(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	storage := &Storage{db: db}

	id := 1

	mock.ExpectExec("DELETE FROM users WHERE id = ?").
		WithArgs(id).
		WillReturnResult(sqlmock.NewResult(0, 0))

	err = storage.DeleteUser(id)

	assert.Error(t, err)
	assert.Equal(t, ErrNoRowsAffected, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
