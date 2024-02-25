package store

import (
	"context"
	"log"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewMockDB() (*gorm.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("error %s opening a stub database connection", err)
	}

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})

	if err != nil {
		log.Fatalf("error %s opening gorm database", err)
	}

	return gormDB, mock
}

func TestInitDB(t *testing.T) {
	mockDB, mock := NewMockDB()

	mock.ExpectExec("CREATE TABLE \"users\"").WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectExec("CREATE TABLE \"transactions\"").WillReturnResult(sqlmock.NewResult(0, 0))

	err := InitDB(mockDB)
	assert.NoError(t, err)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestGetOrCreateUser(t *testing.T) {
	testCases := []struct {
		name          string
		inputUser     User
		expectedUser  User
		expectedError error
	}{
		{
			name:          "New user",
			inputUser:     User{Name: "Jose Rivera"},
			expectedUser:  User{Name: "Jose Rivera"},
			expectedError: nil,
		},
		{
			name: "Existing user",
			inputUser: User{
				Name: "Jose Rivera",
			},
			expectedUser:  User{Name: "Jose Rivera"},
			expectedError: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockDB, mock := NewMockDB()
			rows := sqlmock.NewRows([]string{"name"}).AddRow(tc.inputUser.Name)
			mock.ExpectQuery("SELECT").WillReturnRows(rows)
			user, err := GetOrCreateUser(mockDB, context.Background(), &tc.inputUser)
			assert.Equal(t, tc.expectedUser, *user)
			assert.Equal(t, tc.expectedError, err)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
