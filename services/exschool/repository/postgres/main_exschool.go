package postgres

import (
	"context"
	"database/sql"

	"github.com/shellrean/extraordinary-raport/domain"
)

type repository struct{
	Conn *sql.DB
}

func New(Conn *sql.DB) domain.ExschoolRepository{
	return &repository{
		Conn,
	}
}

func (m *repository) fetch(ctx context.Context, query string, args ...interface{}) (res []domain.Exschool, err error) {
	rows, err := m.Conn.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	defer func() {
		rows.Close()
	}()

	for rows.Next() {
		t := domain.Exschool{}
		err = rows.Scan(
			&t.ID,
			&t.Name,
			&t.CreatedAt,
			&t.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		res = append(res, t)
	}

	return
}

func (m *repository) Fetch(ctx context.Context) (res []domain.Exschool, err error) {
	query := `SELECT id, name, created_at, updated_at
		FROM exschools`
	res, err = m.fetch(ctx, query)
	if err != nil {
		return nil, err
	}
	return
}

func (m *repository) GetByID(ctx context.Context, id int64) (res domain.Exschool, err error) {
	query := `SELECT id, name, created_at, updated_at
		FROM exschools WHERE id=$1`

	list, err := m.fetch(ctx, query, id)
    if err != nil {
        return domain.Exschool{}, err
    }
    if len(list) < 1 {
        return domain.Exschool{}, err
    }
	res = list[0]
	return
}

func (m *repository) Store(ctx context.Context, ex *domain.Exschool) (err error) {
	return
}