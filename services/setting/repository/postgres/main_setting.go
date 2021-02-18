package postgres

import (
	"database/sql"
    "context"
    "fmt"

	"github.com/lib/pq"

	"github.com/shellrean/extraordinary-raport/domain"
)

type repository struct {
	Conn *sql.DB
}

func New(Conn *sql.DB) domain.SettingRepository {
	return &repository{
		Conn,
	}
}

func (m *repository) fetch(ctx context.Context, query string, args ...interface{}) (result []domain.Setting, err error) {
    rows, err := m.Conn.QueryContext(ctx, query, args...)
    if err != nil {
        return nil, err
    }

    defer func() {
        rows.Close()
    }()

    for rows.Next() {
        t := domain.Setting{}
        err = rows.Scan(
            &t.ID,
            &t.Name,
            &t.Value,
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

func (m *repository) GetByName(ctx context.Context, name string) (res domain.Setting, err error) {
	query := `SELECT id,name,value,created_at,updated_at 
		FROM settings WHERE name=$1`
	err = m.Conn.QueryRowContext(ctx, query, name).
		Scan(
			&res.ID,
			&res.Name,
			&res.Value,
			&res.CreatedAt,
			&res.UpdatedAt,
		)
	if err != nil {
		return
	}
	return
}

func (m *repository) Fetch(ctx context.Context, names []string) (res []domain.Setting, err error) {
	query := `SELECT id,name,value,created_at,updated_at
		FROM settings WHERE name = ANY($1)`
    res, err = m.fetch(ctx, query, pq.Array(names))
    if err != nil {
        return nil, err
    }

    return
}

func (m *repository) Update(ctx context.Context, s *domain.Setting) (err error) {
    query := `UPDATE settings SET value=$1, updated_at=$2
            WHERE name=$3`

    result, err := m.Conn.ExecContext(ctx, query, s.Value, s.UpdatedAt, s.Name)
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