package postgres

import (
	"database/sql"
	"context"

	"github.com/shellrean/extraordinary-raport/domain"
)

type sprRepo struct {
	Conn *sql.DB
}

func NewPostgresClassroomSubjectPlanResult(Conn *sql.DB) domain.ClassroomSubjectPlanResultRepository{
	return &sprRepo{
		Conn,
	}
}

func (m sprRepo) Store(ctx context.Context, spr *domain.ClassroomSubjectPlanResult) (err error) {
	query := `INSERT INTO classroom_subject_plan_results (student_id, subject_id, plan_id, number, updated_by)
	VALUES ($1, $2, $3, $4, $5) RETURNING id`

	err = m.Conn.QueryRowContext(ctx, query, spr.Student.ID, spr.Subject.ID, spr.Plan.ID, spr.Number, spr.UpdatedBy.ID).Scan(&spr.ID)

	return
}

