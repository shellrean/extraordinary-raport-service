package postgres

import (
	"database/sql"
	"context"

	"github.com/shellrean/extraordinary-raport/domain"
)

type classroomRepository struct {
	Conn *sql.DB
}

func NewPostgresClassroomRepository(Conn *sql.DB) domain.ClassroomRepository {
	return &classroomRepository{
		Conn,
	}
}

func (m *classroomRepository) fetch(ctx context.Context, query string, args ...interface{})(result []domain.Classroom, err error) {
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

func (m *classroomRepository) Fetch(ctx context.Context) (res []domain.Classroom, err error){
	query := `SELECT id,name,grade,major_id,created_at,updated_at
		FROM classrooms`
	res, err = m.fetch(ctx, query)
	if err != nil {
		return
	}
	return
}