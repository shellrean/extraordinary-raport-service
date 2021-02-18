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

func New(Conn *sql.DB) domain.ClassroomSubjectPlanResultRepository{
	return &repository{
		Conn,
	}
}

func (m *repository) fetch(ctx context.Context, query string, args ...interface{}) (result []domain.ClassroomSubjectPlanResult, err error) {
    rows, err := m.Conn.QueryContext(ctx, query, args...)
    if err != nil {
        return nil, err
    }

    defer func() {
        rows.Close()
    }()

    for rows.Next() {
        t := domain.ClassroomSubjectPlanResult{
		}
        err = rows.Scan(
			&t.ID,
			&t.Index,
			&t.Student.ID,
			&t.Subject.ID,
			&t.Plan.ID,
			&t.Number,
			&t.UpdatedBy.ID,
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

func (m *repository) Store(ctx context.Context, spr *domain.ClassroomSubjectPlanResult) (err error) {
	query := `INSERT INTO classroom_subject_plan_results (index, student_id, subject_id, plan_id, number, updated_by, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id`

	err = m.Conn.QueryRowContext(ctx, query, spr.Index, spr.Student.ID, spr.Subject.ID, spr.Plan.ID, spr.Number, spr.UpdatedBy.ID, spr.CreatedAt, spr.UpdatedAt).Scan(&spr.ID)

	return
}

func (m *repository) FetchByPlan(ctx context.Context, planID int64) (res []domain.ClassroomSubjectPlanResult, err error) {
	query := `SELECT id, index, student_id, subject_id, plan_id, number, updated_by, created_at, updated_at FROM classroom_subject_plan_results WHERE plan_id=$1`

	res, err = m.fetch(ctx, query, planID)
	if err != nil {
        return nil, err
    }

    return
}

func (m *repository) GetByPlanIndexStudent(ctx context.Context, p int64, i int, s int64) (res domain.ClassroomSubjectPlanResult, err error) {
	query := `SELECT id, index, student_id, subject_id, plan_id, number, updated_by, created_at, updated_at FROM classroom_subject_plan_results WHERE plan_id=$1 AND index=$2 AND student_id=$3`

    list, err := m.fetch(ctx, query, p, i, s)
    if err != nil {
        return domain.ClassroomSubjectPlanResult{}, err
    }
    if len(list) < 1 {
        return domain.ClassroomSubjectPlanResult{}, err
    }
    res = list[0]
    return
}

func (m *repository) Update(ctx context.Context, spr *domain.ClassroomSubjectPlanResult) (err error) {
    query := `UPDATE classroom_subject_plan_results SET index=$1, student_id=$2, subject_id=$3, plan_id=$4, number=$5, updated_by=$6, updated_at=$7 WHERE id=$8`

	result, err := m.Conn.ExecContext(ctx, query, spr.Index, spr.Student.ID, spr.Subject.ID, spr.Plan.ID, spr.Number, spr.UpdatedBy.ID, spr.UpdatedAt, spr.ID)
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