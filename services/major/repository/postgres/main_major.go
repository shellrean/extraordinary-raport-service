package postgres

import (
	"database/sql"
	"context"

	"github.com/shellrean/extraordinary-raport/domain"
)

type repository struct {
	Conn *sql.DB
}

func New(Conn *sql.DB) domain.MajorRepository{
	return &repository{
		Conn,
	}
}

func (m *repository) fetch(ctx context.Context, query string, args ...interface{}) (result []domain.Major, err error) {
	rows, err := m.Conn.QueryContext(ctx, query, args...)
	if err != nil {
		return
	}

	defer func() {
		rows.Close()
	}()

	result = []domain.Major{}
	for rows.Next() {
		t := domain.Major{}
		err = rows.Scan(
			&t.ID,
			&t.Name,
			&t.CreatedAt,
			&t.UpdatedAt,
		)

		if err != nil {
			return 
		}

		result = append(result, t)
	}
	return
}

func (m *repository) Fetch(ctx context.Context) (res []domain.Major, err error) {
	query := `SELECT id, name, created_at, updated_at
		FROM majors`
	
	res, err = m.fetch(ctx, query)
	if err != nil {
		return
	}
	return
}

func (m *repository) GetByID(ctx context.Context, id int64) (res domain.Major, err error) {
	query := `SELECT id,name,created_at,updated_at FROM majors
		WHERE id = $1`
	list, err := m.fetch(ctx, query, id)
	if err != nil {
		return
	}
	if len(list) < 1 {
		return
	}
	res = list[0]
	return
}