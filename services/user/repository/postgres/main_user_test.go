package postgres_test

import (
    "database/sql"
    "log"
    "testing"
    "time"
    "context"
    
    sqlmock "github.com/DATA-DOG/go-sqlmock"
    "github.com/stretchr/testify/assert"

    _userRepo "github.com/shellrean/extraordinary-raport/services/user/repository/postgres"
    "github.com/shellrean/extraordinary-raport/domain"
)

func NewMock() (*sql.DB, sqlmock.Sqlmock) {
    db, mock, err := sqlmock.New()
    if err != nil {
        log.Fatalf("an error '%s' was not expected when opening a stub connection", err)
    }

    return db, mock
}

func TestGetByEmail(t *testing.T) {
    db, mock := NewMock()
    
    u := domain.User{
        ID:         1,
        Name:       "user1",
        Email:      "user1@shellrean.com",
        Password:   "password",
        CreatedAt:  time.Now(),
        UpdatedAt:  time.Now(),
    }

    rows := sqlmock.NewRows([]string{"id","name","email","password","created_at","updated_at"}).
        AddRow(u.ID, u.Name, u.Email, u.Password, u.CreatedAt, u.UpdatedAt)
    
    query := "SELECT id,name,email,password,created_at,updated_at FROM users WHERE email=\\$1"

    mock.ExpectQuery(query).WillReturnRows(rows)

    repo := _userRepo.NewPostgresUserRepository(db)
    user, err := repo.GetByEmail(context.TODO(), u.Email)

    assert.NotNil(t, user)
    assert.NoError(t, err)
}

func TestGetByEmailError(t *testing.T) {
    db, mock := NewMock()

    u := domain.User{
        ID:         1,
        Name:       "user1",
        Email:      "user1@shellrean.com",
        Password:   "password",
        CreatedAt:  time.Now(),
        UpdatedAt:  time.Now(),
    }

    rows := sqlmock.NewRows([]string{"id","name","email","password","created_at","updated_at"}).
        AddRow(u.ID, u.Name, u.Email, u.Password, u.CreatedAt, u.UpdatedAt)

    query := "SELECT id,name,email,password,created_at,updated_at FROM user WHERE email=\\$1"

    mock.ExpectQuery(query).WillReturnRows(rows)

    repo := _userRepo.NewPostgresUserRepository(db)
    user, err := repo.GetByEmail(context.TODO(), u.Email)

    assert.Empty(t, user)
    assert.Error(t, err)
}