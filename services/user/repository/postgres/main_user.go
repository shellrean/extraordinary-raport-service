package postgres

import (
    "context"
    "database/sql"
    "fmt"
    "github.com/lib/pq"
    "strings"

    "github.com/shellrean/extraordinary-raport/domain"
)

type repository struct {
    Conn *sql.DB
}

func New(Conn *sql.DB) domain.UserRepository {
    return &repository{
        Conn,
    }
}

func (m *repository) fetch(ctx context.Context, query string, args ...interface{}) (result []domain.User, err error) {
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
            &t.Role,
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

func (m *repository) Fetch(ctx context.Context, q string, cursor int64, num int64) (res []domain.User, err error) {
    query := `SELECT id,name,email,password,role,created_at,updated_at
            FROM users WHERE (LOWER(name) LIKE '%' || $1 || '%' OR email LIKE '%' || $2 || '%') AND role=$3 AND id > $4 ORDER BY id LIMIT $5`

    res, err = m.fetch(ctx, query, q, q, domain.RoleTeacher, cursor, num)
    if err != nil {
        return nil, err
    }

    return
}

func (m *repository) GetByID(ctx context.Context, id int64) (res domain.User, err error) {
    query := `SELECT id,name,email,password,role,created_at,updated_at
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

func (m *repository) GetByEmail(ctx context.Context, email string) (res domain.User, err error) {
    query := `SELECT id,name,email,password,role,created_at,updated_at
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

func (m *repository) Store(ctx context.Context, u *domain.User) (err error) {
    query := `INSERT INTO users (name, email, password, role,created_at, updated_at)
            VALUES ($1,$2,$3,$4,$5,$6) RETURNING id`
        
    err = m.Conn.QueryRowContext(ctx, query, u.Name, u.Email, u.Password, u.Role,u.CreatedAt, u.UpdatedAt).Scan(&u.ID)
    if err != nil {
        return err
    }
    
    return
}

func (m *repository) StoreMultiple(ctx context.Context, us []domain.User) (err error) {
    var valueStrings []string
    var valueArgs []interface{}

    for i, item := range us {
        valueStrings = append(valueStrings, fmt.Sprintf(
            "($%d, $%d, $%d, $%d, $%d, $%d)", 
            (6*i+1), (6*i+2), (6*i+3), (6*i+4), (6*i+5), (6*i+6),
        ))

        valueArgs = append(valueArgs, item.Name)
        valueArgs = append(valueArgs, item.Email)
        valueArgs = append(valueArgs, item.Password)
        valueArgs = append(valueArgs, item.Role)
        valueArgs = append(valueArgs, item.CreatedAt)
        valueArgs = append(valueArgs, item.UpdatedAt)
    }
    
    query := `INSERT INTO users (name, email, password, role,created_at, updated_at) VALUES %s`
    query = fmt.Sprintf(query, strings.Join(valueStrings, ","))

    tx, err := m.Conn.Begin()
    if err != nil {
        return
    }

    _, err = tx.ExecContext(ctx, query, valueArgs...)
    if err != nil {
        _ = tx.Rollback()
        return
    }
    
    err = tx.Commit()
    
    return
}

func (m *repository) Update(ctx context.Context, u *domain.User) (err error) {
    query := `UPDATE users SET name=$1, email=$2, password=$3, role=$4, updated_at=$5
            WHERE id=$6`

    result, err := m.Conn.ExecContext(ctx, query, u.Name, u.Email, u.Password, u.Role, u.UpdatedAt, u.ID)
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

func (m *repository) Delete(ctx context.Context, id int64) (err error) {
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

func (m *repository) DeleteMultiple(ctx context.Context, ids []string) (err error) {
    query := `DELETE FROM users WHERE id = ANY($1)`

    _, err = m.Conn.ExecContext(ctx, query, pq.Array(ids))
    if err != nil {
        return err
    }

    return
}