package postgres

import (
    "context"
    "database/sql"
    "fmt"

    "github.com/shellrean/extraordinary-raport/domain"
)

type postgresUserRepository struct {
    Conn *sql.DB
}

func NewPostgresUserRepository(Conn *sql.DB) domain.UserRepository {
    return &postgresUserRepository{
        Conn,
    }
}

func (m *postgresUserRepository) fetch(ctx context.Context, query string, args ...interface{}) (result []domain.User, err error) {
    rows, err := m.Conn.QueryContext(ctx, query, args...)
    if err != nil {
        return nil, err
    }

    defer func() {
        rows.Close()
    }()

    result = make([]domain.User, 0)
    for rows.Next() {
        t := domain.User{}
        err = rows.Scan(
            &t.ID,
            &t.Name,
            &t.Email,
            &t.Password,
            &t.CreatedAt,
            &t.UpdatedAt,
        )

        if err != nil {
            return nil, err
        }

        result = append(result, t)
    }

    return
}

func (m *postgresUserRepository) Fetch(ctx context.Context, cursor int64, num int64) (res []domain.User, err error) {
    query := `SELECT id,name,email,password,created_at,updated_at
            FROM users WHERE id > $1 ORDER BY id LIMIT $2`

    res, err = m.fetch(ctx, query, cursor, num)
    if err != nil {
        return nil, err
    }

    return
}

func (m *postgresUserRepository) GetByID(ctx context.Context, id int64) (res domain.User, err error) {
    query := `SELECT id,name,email,password,created_at,updated_at
            FROM users WHERE id=$1`

    list, err := m.fetch(ctx, query, id)
    if err != nil {
        return domain.User{}, err
    }
    if len(list) < 1 {
        return domain.User{}, err
    }
    res = list[0]
    return
}

func (m *postgresUserRepository) GetByEmail(ctx context.Context, email string) (res domain.User, err error) {
    query := `SELECT id,name,email,password,created_at,updated_at
            FROM users WHERE email=$1`

    list, err := m.fetch(ctx, query, email)
    if err != nil {
        return domain.User{}, err
    }

    if len(list) < 1 {
        return domain.User{}, nil
    }

    res = list[0]

    return
}

func (m *postgresUserRepository) Store(ctx context.Context, u *domain.User) (err error) {
    query := `INSERT INTO users (name, email, password, created_at, updated_at)
            VALUES ($1,$2,$3,$4,$5)`
        
    result, err := m.Conn.ExecContext(ctx, query, u.Name, u.Email, u.Password, u.CreatedAt, u.UpdatedAt)
    if err != nil {
        return err
    }
    rows, err := result.RowsAffected()
    if err != nil {
        return err
    }
    if rows != 1 {
        return fmt.Errorf("expected single row affected, got %d rows affected", rows)
    }

    userId, err := result.LastInsertId()
    if err != nil {
        return err
    }
    u.ID = userId
    return
}

func (m *postgresUserRepository) Update(ctx context.Context, u *domain.User) (err error) {
    query := `UPDATE users SET name=$1, email=$2, password=$3, updated_at=$4
            WHERE id=$5`

    result, err := m.Conn.ExecContext(ctx, query, u.Name, u.Email, u.Password, u.UpdatedAt, u.ID)
    if err != nil {
        return err
    }
    rows, err := result.RowsAffected()
    if err != nil {
        return err
    }
    if rows != 1 {
        return fmt.Errorf("expected single row affected, got %d rows affected", rows)
    }
    return
}

func (m *postgresUserRepository) Delete(ctx context.Context, id int64) (err error) {
    query := `DELETE FROM users WHERE id=$1`
    result, err := m.Conn.ExecContext(ctx, query, id)
    if err != nil {
        return err
    }
    rows, err := result.RowsAffected()
    if err != nil {
        return err
    }
    if rows != 1 {
        return fmt.Errorf("expected single row affected, got %d rows affected", rows)
    }
    return
}