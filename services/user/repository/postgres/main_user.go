package postgres

import (
    "context"
    "database/sql"

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

func (m *postgresUserRepository) Fetch(ctx context.Context, num int64) (res []domain.User, err error) {
    query := `SELECT id,name,email,password,created_at,updated_at
            FROM users ORDER BY created_at LIMIT $1`

    res, err = m.fetch(ctx, query, num)
    if err != nil {
        return nil, err
    }

    return
}

func (m *postgresUserRepository) GetByEmail(ctx context.Context, email string) (res domain.User, err error) {
    query := `SELECT id,name,email,password,created_at,updated_at
            FROM users WHERE email=$1`

    list, err := m.fetch(ctx, query, email)
    if err != nil {
        return domain.User{}, err
    }

    if len(list) > 1 {
        return domain.User{}, nil
    }

    res = list[0]

    return
}