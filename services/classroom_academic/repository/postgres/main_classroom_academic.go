package postgres

import (
	"database/sql"
	"context"

	"github.com/shellrean/extraordinary-raport/domain"
)

type classroomAcademicRepository struct {
	Conn *sql.DB
}

func NewPostgresClassroomAcademicRepository(Conn *sql.DB) domain.ClassroomAcademicRepository {
	return &classroomAcademicRepository {
		Conn,
	}
}

func (m *classroomAcademicRepository) fetch(ctx context.Context, query string, args ...interface{}) (result []domain.ClassroomAcademic, err error) {
    rows, err := m.Conn.QueryContext(ctx, query, args...)
    if err != nil {
        return nil, err
    }

    defer func() {
        rows.Close()
    }()

    result = []domain.ClassroomAcademic{}
    for rows.Next() {
		t := domain.ClassroomAcademic{}
		var academicID int64
		var teacherId int64
        err = rows.Scan(
            &t.ID,
			&academicID,
			&teacherId,
            &t.CreatedAt,
            &t.UpdatedAt,
        )

        if err != nil {
            return nil, err
		}
		
		academic := domain.Academic {
			ID:	academicID,
		}
		teacher := domain.User {
			ID: teacherId,
		}
		t.Academic = academic
		t.Teacher = teacher

        result = append(result, t)
    }

    return
}

func (m *classroomAcademicRepository) Fetch(ctx context.Context) (res []domain.ClassroomAcademic, err error) {
	query := `SELECT id, academic_id, teacher_id, created_at, updated_at
		FROM classroom_academics`
	
	res, err = m.fetch(ctx, query)
	if err != nil {
		return
	}
	return
}