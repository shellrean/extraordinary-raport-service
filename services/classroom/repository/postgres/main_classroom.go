package postgres

import (
	"database/sql"
	"context"
	"fmt"

	"github.com/shellrean/extraordinary-raport/domain"
)

type repository struct {
	Conn *sql.DB
}

func New(Conn *sql.DB) domain.ClassroomRepository {
	return &repository{
		Conn,
	}
}

func (m *repository) fetch(ctx context.Context, query string, args ...interface{})(result []domain.Classroom, err error) {
	rows, err := m.Conn.QueryContext(ctx, query, args...)
	if err != nil {
		return
	}

	defer func() {
		rows.Close()
	}()

	result = []domain.Classroom{}
	for rows.Next() {
		t := domain.Classroom{}
		var majorId int64
		err = rows.Scan(
			&t.ID,
			&t.Name,
			&t.Grade,
			&majorId,
			&t.CreatedAt,
			&t.UpdatedAt,
		)

		if err != nil {
			return
		}

		t.Major = domain.Major{
			ID: majorId,
		}
		result = append(result, t)
	}

	return
}

func (m *repository) GetByID(ctx context.Context, id int64) (res domain.Classroom, err error) {
	query := `SELECT id,name,grade,major_id,created_at,updated_at
		FROM classrooms WHERE id=$1`
	
	list, err := m.fetch(ctx, query, id)
	if err != nil {
		return
	}
	if len(list) < 1 {
		res = domain.Classroom{}
		return 
	}
	res = list[0]
	return
}

func (m *repository) Fetch(ctx context.Context) (res []domain.Classroom, err error){
	query := `SELECT id,name,grade,major_id,created_at,updated_at
		FROM classrooms`
	res, err = m.fetch(ctx, query)
	if err != nil {
		return
	}
	return
}

func (m *repository) Store(ctx context.Context, c *domain.Classroom) (err error) {
	query := `INSERT INTO classrooms (name,grade,major_id,created_at,updated_at)
		VALUES($1,$2,$3,$4,$5) RETURNING id`
	err = m.Conn.QueryRowContext(ctx, query, c.Name, c.Grade, c.Major.ID, c.CreatedAt, c.UpdatedAt).Scan(&c.ID)
	if err != nil {
		return
	}
	return
}

func (m *repository) Update(ctx context.Context, c *domain.Classroom) (err error) {
	query := `UPDATE classrooms SET name=$1, grade=$2, major_id=$3, updated_at=$4
		WHERE id=$5`
	result, err := m.Conn.ExecContext(ctx, query, c.Name, c.Grade, c.Major.ID, c.UpdatedAt, c.ID)
	if err != nil {
		return
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return
	}
	if rows != 1 {
		return fmt.Errorf("expected single row affected, got %d rows affected", rows)
	}
	return
}

func (m *repository) Delete(ctx context.Context, id int64) (err error) {
	query := `DELETE FROM classrooms WHERE id=$1`
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