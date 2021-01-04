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
		t := domain.ClassroomAcademic{
			Academic: domain.Academic{},
			Classroom: domain.Classroom{},
			Teacher: domain.User{},
		}
        err = rows.Scan(
            &t.ID,
			&t.Academic.ID,
			&t.Classroom.ID,
			&t.Teacher.ID,
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

func (m *classroomAcademicRepository) Fetch(ctx context.Context) (res []domain.ClassroomAcademic, err error) {
	query := `SELECT id, academic_id, classroom_id, teacher_id, created_at, updated_at
		FROM classroom_academics`
	
	res, err = m.fetch(ctx, query)
	if err != nil {
		return
	}
	return
}

func (m *classroomAcademicRepository) Store(ctx context.Context, ca *domain.ClassroomAcademic) (err error) {
	query := `INSERT INTO classroom_academics (academic_id, classroom_id, teacher_id, created_at, updated_at)
		VALUES($1,$2,$3,$4,$5) RETURNING id`
	err = m.Conn.QueryRowContext(ctx, query, 
		ca.Academic.ID,
		ca.Classroom.ID,
		ca.Teacher.ID,
		ca.CreatedAt,
		ca.UpdatedAt,
	).Scan(&ca.ID)
	if err != nil {
		return
	}
	return 
}

func (m *classroomAcademicRepository) GetByAcademicAndClass(ctx context.Context, a int64, c int64) (res domain.ClassroomAcademic, err error) {
	query := `SELECT id,academic_id,classroom_id, teacher_id,createdAt,updated_at
		FROM classroom_academics WHERE academic_id=$1 AND classroom_id=$2`

	list, err := m.fetch(ctx, query, a, c)
	if err != nil {
		return
	}
	if len(list) < 1 {
		return domain.ClassroomAcademic{}, err
	}
	res = list[0]
	return
}