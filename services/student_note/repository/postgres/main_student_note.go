package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/shellrean/extraordinary-raport/domain"
)

type repository struct {
	Conn *sql.DB
}

func New(Conn *sql.DB) domain.StudentNoteRepository {
	return &repository {
		Conn,
	}
}

func (m *repository) fetch(ctx context.Context, query string, args ...interface{}) (result []domain.StudentNote, err error) {
	rows, err := m.Conn.QueryContext(ctx, query, args...)
	if err != nil {
		return
	}

	defer func() {
		rows.Close()
	}()

	for rows.Next() {
		t := domain.StudentNote{}
		err = rows.Scan(
			&t.ID,
			&t.Type,
			&t.Student.ID,
			&t.Teacher.ID,
			&t.Note,
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

func (m *repository) FetchByClassroom(ctx context.Context, id int64) (res []domain.StudentNote, err error) {
	query := `SELECT 
		sn.id, 
		sn.type, 
		sn.student_id, 
		sn.created_by, 
		sn.note, 
		sn.created_at, 
		sn.updated_at 
	FROM student_notes sn
	INNER JOIN classroom_students cs
		ON sn.student_id = cs.id
	WHERE cs.classroom_academic_id=$1`

	res, err  = m.fetch(ctx, query, id)
	if err != nil {
		return nil, err
	}
	return
}

func (m *repository) FetchByTypeAndClassroom(ctx context.Context, id int64, typ int64) (res []domain.StudentNote, err error) {
	query := `SELECT 
		sn.id, 
		sn.type, 
		sn.student_id, 
		sn.created_by, 
		sn.note, 
		sn.created_at, 
		sn.updated_at 
	FROM student_notes sn
	INNER JOIN classroom_students cs
		ON sn.student_id = cs.id
	WHERE cs.classroom_academic_id=$1 AND sn.type=$2`

	res, err  = m.fetch(ctx, query, id, typ)
	if err != nil {
		return nil, err
	}
	return
}

func (m *repository) GetByStudentAndType(ctx context.Context, id int64, typ int) (res domain.StudentNote, err error) {
	query := `SELECT 
		sn.id, 
		sn.type, 
		sn.student_id, 
		sn.created_by, 
		sn.note, 
		sn.created_at, 
		sn.updated_at 
	FROM student_notes sn
	WHERE sn.type=$1 AND sn.student_id=$2`

	list, err := m.fetch(ctx, query, typ, id)
    if err != nil {
        return domain.StudentNote{}, err
    }
    if len(list) < 1 {
        return domain.StudentNote{}, err
    }
    res = list[0]
    return
}

func (m *repository) Store(ctx context.Context, sn *domain.StudentNote) (err error) {
	query := `INSERT INTO student_notes (type, created_by, student_id, note, created_at, updated_at)
		VALUES($1, $2, $3, $4, $5, $6) RETURNING id`
	err = m.Conn.QueryRowContext(ctx, query,
		sn.Type,
		sn.Teacher.ID,
		sn.Student.ID,
		sn.Note,
		sn.CreatedAt,
		sn.UpdatedAt,
	).Scan(&sn.ID)
	if err != nil {
		return
	}
	return
}

func (m *repository) Update(ctx context.Context, sn *domain.StudentNote) (err error) {
	query := `UPDATE student_notes SET type=$1, note=$2, updated_at=$3 WHERE id=$4`
	result, err := m.Conn.ExecContext(ctx, query, sn.Type, sn.Note, sn.UpdatedAt, sn.ID)
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