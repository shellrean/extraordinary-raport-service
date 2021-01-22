package postgres

import (
	"database/sql"
	"context"

	"github.com/shellrean/extraordinary-raport/domain"
)

type csuRepo struct {
	Conn 	*sql.DB
}

func NewPostgresClassroomSubjectRepository(Conn *sql.DB) domain.ClassroomSubjectRepository {
	return &csuRepo{
		Conn,
	}
}

func (m *csuRepo) FetchByClassroom(ctx context.Context, academicClassroomID int64) (res []domain.ClassroomSubject, err error) {
	query := `SELECT
		cs.id,
		cs.classroom_academic_id,
		cs.subject_id,
		cs.teacher_id,
		cs.mgn,
		u.name,
		s.name,
		cs.created_at,
		cs.updated_at
	FROM 
		classroom_subjects cs
	INNER JOIN users u
		ON u.id = cs.teacher_id
	INNER JOIN subjects s
		ON s.id = cs.subject_id
	WHERE cs.classroom_academic_id=$1`

	rows, err := m.Conn.QueryContext(ctx, query, academicClassroomID)
	if err != nil {
		return nil, err
	}

	defer func() {
		rows.Close()
	}()

	for rows.Next() {
		t := domain.ClassroomSubject{
			ClassroomAcademic: domain.ClassroomAcademic{},
			Subject: domain.Subject{},
			Teacher: domain.User{},
		}
		err = rows.Scan(
			&t.ID,
			&t.ClassroomAcademic.ID,
			&t.Subject.ID,
			&t.Teacher.ID,
			&t.MGN,
			&t.Teacher.Name,
			&t.Subject.Name,
			&t.CreatedAt,
			&t.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		res = append(res, t)
	}

	return
}