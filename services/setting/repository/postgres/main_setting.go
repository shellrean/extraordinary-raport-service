package postgres

import (
	"database/sql"
	"context"

	"github.com/lib/pq"

	"github.com/shellrean/extraordinary-raport/domain"
)

type settingRepository struct {
	Conn *sql.DB
}

func NewPostgresSettingRepository(Conn *sql.DB) domain.SettingRepository {
	return &settingRepository{
		Conn,
	}
}

func (m *settingRepository) fetch(ctx context.Context, query string, args ...interface{}) (result []domain.Setting, err error) {
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

func (m *settingRepository) GetByName(ctx context.Context, name string) (res domain.Setting, err error) {
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

func (m *settingRepository) Fetch(ctx context.Context, names []string) (res []domain.Setting, err error) {
	query := `SELECT id,name,value,created_at,updated_at
		FROM settings WHERE name = ANY($1)`
    res, err = m.fetch(ctx, query, pq.Array(names))
    if err != nil {
        return nil, err
    }

    return
}