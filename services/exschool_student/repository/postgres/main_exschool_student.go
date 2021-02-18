package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/shellrean/extraordinary-raport/domain"
)

type repository struct {
	Conn 	*sql.DB
}

func New(Conn *sql.DB) domain.ExschoolStudentRepository {
	return &repository {
		Conn,
	}
}

func (m *repository) fetch(ctx context.Context, query string, args ...interface{}) (res []domain.ExschoolStudent, err error) {
	rows, err := m.Conn.QueryContext(ctx, query, args...)
	if err != nil {
		return
	}

	defer func() {
		rows.Close()
	}()

	for rows.Next() {
		t := domain.ExschoolStudent{}
		err = rows.Scan(
			&t.ID,
			&t.Exschool.ID,
			&t.Student.ID,
			&t.CreatedAt,
			&t.UpdatedAt,
		)

		if err != nil {
			return
		}

		res = append(res, t)
	}
	return
}

func (m *repository) FetchByClassroom(ctx context.Context, cID int64) (res []domain.ExschoolStudent, err error) {
	query := `
	SELECT
		exs.id,
		exs.exschool_id,
		exs.student_id,
		exs.created_at,
		exs.updated_at
	FROM exschool_students exs
	INNER JOIN classroom_students cs
		ON cs.id=exs.student_id
	WHERE cs.classroom_academic_id=$1`

	res, err = m.fetch(ctx, query, cID)
	if err != nil {
		return
	}
	return
}

func (m *repository) GetByID(ctx context.Context, id int64) (res domain.ExschoolStudent, err error) {
	query := `
	SELECT
		exs.id,
		exs.exschool_id,
		exs.student_id,
		exs.created_at,
		exs.updated_at
	FROM exschool_students exs 
	WHERE exs.id=$1`

	list, err := m.fetch(ctx, query, id)
    if err != nil {
        return domain.ExschoolStudent{}, err
    }
    if len(list) < 1 {
        return domain.ExschoolStudent{}, err
    }
    res = list[0]
    return
}

func (m *repository) GetByExschoolAndStudent(ctx context.Context, exID, sID int64) (res domain.ExschoolStudent, err error) {
	query := `
	SELECT
		exs.id,
		exs.exschool_id,
		exs.student_id,
		exs.created_at,
		exs.updated_at
	FROM exschool_students exs 
	WHERE exs.exschool_id=$1 AND exs.student_id=$2`

	list, err := m.fetch(ctx, query, exID, sID)
    if err != nil {
        return domain.ExschoolStudent{}, err
    }
    if len(list) < 1 {
        return domain.ExschoolStudent{}, err
    }
    res = list[0]
    return
}

func (m *repository) Store(ctx context.Context, exs *domain.ExschoolStudent) (err error) {
	query := `INSERT INTO exschool_students (exschool_id, student_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4) RETURNING id`
	err = m.Conn.QueryRowContext(ctx, query, exs.Exschool.ID, exs.Student.ID, exs.CreatedAt, exs.UpdatedAt).Scan(&exs.ID)
	if err != nil {
		return err
	}
	return
}

func (m *repository) Delete(ctx context.Context, id int64) (err error) {
	query := `DELETE FROM exschool_students WHERE id=$1`
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