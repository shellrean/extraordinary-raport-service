package postgres

import (
	"fmt"
	"context"
	"database/sql"

	"github.com/shellrean/extraordinary-raport/domain"
)

type repository struct {
	Conn	*sql.DB
}

func New(Conn *sql.DB) domain.AttendanceRepository {
	return &repository {
		Conn,
	}
}

func (m *repository) fetch(ctx context.Context, query string, args ...interface{}) (res []domain.Attendance, err error) {
	rows, err := m.Conn.QueryContext(ctx, query, args...)
	if err != nil {
		return
	}

	defer func() {
		rows.Close()
	}()

	for rows.Next() {
		t := domain.Attendance{}
		err = rows.Scan(
			&t.ID,
			&t.Student.ID,
			&t.Total,
			&t.Type,
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

func (m *repository) Fetch(ctx context.Context, cid int64) (res []domain.Attendance, err error) {
	query := `SELECT 
		a.id, 
		a.student_id, 
		a.total, 
		a.type, 
		a.created_at, 
		a.updated_at
	FROM attendances
	INNER JOIN classroom_students cs
		ON cs.id = a.student_id
	WHERE cs.classroom_academic_id=$1`

	res, err = m.fetch(ctx, query, cid)
	if err != nil {
		return
	}
	return
}

func (m *repository) GetByStudentAndType(ctx context.Context, sid int64, typ int) (res domain.Attendance, err error) {
	query := `SELECT id, student_id, total, type, created_at, updated_at
	FROM attendances WHERE student_id=$1 AND type=$2`

	list, err := m.fetch(ctx, query, sid, typ)
	if err != nil {
		return
	}
    if len(list) < 1 {
		return domain.Attendance{}, nil
    }
    res = list[0]
	return
}

func (m *repository) Store(ctx context.Context, u *domain.Attendance) (err error) {
    query := `INSERT INTO attendances (student_id, total, type, created_at, updated_at)
        VALUES ($1,$2,$3,$4,$5) RETURNING id`
        
    err = m.Conn.QueryRowContext(ctx, query, u.Student.ID, u.Total, u.Type, u.CreatedAt, u.UpdatedAt).Scan(&u.ID)
    if err != nil {
        return err
    }
    
    return
}

func (m *repository) Update(ctx context.Context, u *domain.Attendance) (err error) {
    query := `UPDATE attendances SET total=$1, updated_at=$2 WHERE id=$3`

    result, err := m.Conn.ExecContext(ctx, query, u.Total, u.UpdatedAt, u.ID)
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