package postgres

import (
    "context"
	"database/sql"

	"github.com/shellrean/extraordinary-raport/domain"
)

type csRepository struct {
	Conn 	*sql.DB
}

func NewPostgresClassroomStudentRepository(Conn *sql.DB) domain.ClassroomStudentRepository {
	return &csRepository {
		Conn,
	}
}

func (m *csRepository) fetch(ctx context.Context, query string, args ...interface{}) (result []domain.ClassroomStudent, err error) {
    rows, err := m.Conn.QueryContext(ctx, query, args...)
    if err != nil {
        return nil, err
    }

    defer func() {
        rows.Close()
    }()

    for rows.Next() {
        t := domain.ClassroomStudent{
			ClassroomAcademic: 	domain.ClassroomAcademic{},
			Student: 			domain.Student{},
		}
        err = rows.Scan(
			&t.ID,
			&t.ClassroomAcademic.ID,
			&t.Student.ID,
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

func (m *csRepository) Fetch(ctx context.Context, cursor int64, num int64) (res []domain.ClassroomStudent, err error) {
    query := `SELECT id,classroom_academic_id,student_id,created_at,updated_at
            FROM classroom_students WHERE id > $1 ORDER BY id LIMIT $2`

    res, err = m.fetch(ctx, query, cursor, num)
    if err != nil {
        return nil, err
    }

    return
}