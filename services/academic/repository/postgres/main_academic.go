package repository

import (
	"database/sql"
	"context"
	"fmt"

	"github.com/shellrean/extraordinary-raport/domain"
)

type postgresAcademicRepository struct {
	Conn *sql.DB
}

func NewPostgresAcademicRepository(Conn *sql.DB) domain.AcademicRepository {
	return &postgresAcademicRepository{
		Conn,
	}
}

func (m *postgresAcademicRepository) fetch(ctx context.Context, query string, args ...interface{}) (result []domain.Academic, err error) {
	rows, err := m.Conn.QueryContext(ctx, query, args...)
	if err != nil {
		return
	}

	defer func() {
		rows.Close()
	}()

	result = []domain.Academic{}
	for rows.Next() {
		t := domain.Academic{}
		err = rows.Scan(
			&t.ID,
			&t.Name,
			&t.Semester,
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

func (m *postgresAcademicRepository) Fetch(ctx context.Context) (res []domain.Academic, err error) {
	query := `SELECT id,name,semester,created_at,updated_at
		FROM academics`

	res, err = m.fetch(ctx, query)
	if err != nil {
		return
	}

	return
}

func (m *postgresAcademicRepository) GetByID(ctx context.Context, id int64) (res domain.Academic, err error) {
	query := `SELECT id,name,semester,created_at,updated_at
		FROM academics WHERE id=$1`
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

func (m *postgresAcademicRepository) GetByYearAndSemester(ctx context.Context, year int, sem int) (res domain.Academic, err error) {
	query := `SELECT id,name,semester,created_at,updated_at
		FROM academics WHERE name=$1 AND semester=$2`
	list, err := m.fetch(ctx, query, year, sem)
	if err != nil {
		return
	}
	if len(list) < 1 {
		return
	}
	res = list[0]
	return
}

func (m *postgresAcademicRepository) Store(ctx context.Context, ac *domain.Academic) (err error) {
	query := `INSERT INTO academics (name,semester,created_at,updated_at)
		VALUES($1,$2,$3,$4) RETURNING id`
	err = m.Conn.QueryRowContext(ctx, query, ac.Name, ac.Semester, ac.CreatedAt, ac.UpdatedAt).Scan(&ac.ID)
	if err != nil {
		return
	}
	return
}

func (m *postgresAcademicRepository) Delete(ctx context.Context, id int64) (err error) {
	query := `DELETE FROM academics WHERE id=$1`
	result, err := m.Conn.ExecContext(ctx, query, id)
	if err != nil {
		return
	}

	rows, err :=result.RowsAffected()
	if err != nil {
		return
	}

	if rows != 1 {
		return fmt.Errorf("expected single row affected, got %d rows affected", rows)
	}
	return
}