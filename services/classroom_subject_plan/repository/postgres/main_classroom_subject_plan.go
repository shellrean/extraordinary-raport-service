package postgres

import (
	"database/sql"
	"context"
	"fmt"

	"github.com/shellrean/extraordinary-raport/domain"
)

type csPlanRepo struct {
	Conn *sql.DB
}

func NewPostgresClassroomSubjectPlanRepository(Conn *sql.DB) domain.ClassroomSubjectPlanRepository {
	return &csPlanRepo {
		Conn,
	}
}

func (m csPlanRepo) fetch(ctx context.Context, query string, args ...interface{}) (res []domain.ClassroomSubjectPlan, err error) {
	rows, err := m.Conn.QueryContext(ctx, query, args...)
	if err != nil {
		return
	}

	defer func() {
		rows.Close()
	}()

	for rows.Next() {
		csp := domain.ClassroomSubjectPlan{
			Teacher:	domain.User{},
			Subject:	domain.ClassroomSubject{},
			Classroom:	domain.ClassroomAcademic{},
		}
		err = rows.Scan(
			&csp.ID,
			&csp.Type,
			&csp.Name,
			&csp.Desc,
			&csp.Teacher.ID,
			&csp.Subject.ID,
			&csp.Subject.Subject.Name,
			&csp.Classroom.ID,
			&csp.CountPlan,
			&csp.MaxPoint,
			&csp.CreatedAt,
			&csp.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		res = append(res, csp)
	}

	return
}

func (m csPlanRepo) Store(ctx context.Context, csp *domain.ClassroomSubjectPlan) (err error){
	query := `INSERT INTO classroom_subject_plans (type,name,description,teacher_id,classroom_subject_id,classroom_academic_id,count_plan,max_point,created_at,updated_at) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10) RETURNING id`

	err = m.Conn.QueryRowContext(ctx, query, 
		csp.Type,
		csp.Name,
		csp.Desc,
		csp.Teacher.ID,
		csp.Subject.ID,
		csp.Classroom.ID,
		csp.CountPlan,
		csp.MaxPoint,
		csp.CreatedAt,
		csp.UpdatedAt,
	).Scan(&csp.ID)

	return
}

func (m csPlanRepo) GetByID(ctx context.Context, id int64) (res domain.ClassroomSubjectPlan, err error) {
	query := `
	SELECT 
		csp.id,
		csp.type, 
		csp.name, 
		csp.description, 
		csp.teacher_id,
		csp.classroom_subject_id,
		s.name,
		csp.classroom_academic_id,
		csp.count_plan,
		csp.max_point,
		csp.created_at,
		csp.updated_at
	FROM 
		classroom_subject_plans csp
	INNER JOIN classroom_subjects cs
		ON cs.id = csp.classroom_subject_id
	INNER JOIN subjects s
		ON s.id = cs.subject_id
	WHERE csp.id=$1`

	list, err := m.fetch(ctx, query, id)
	if err != nil {
		return domain.ClassroomSubjectPlan{}, err
	} 
	if len(list) < 1 {
		return domain.ClassroomSubjectPlan{}, err
	}
	res = list[0]
	return
}

func (m csPlanRepo) FetchByClassroom(ctx context.Context, id int64) (res []domain.ClassroomSubjectPlan, err error) {
	query := `
	SELECT 
		csp.id,
		csp.type, 
		csp.name, 
		csp.description, 
		csp.teacher_id,
		csp.classroom_subject_id,
		s.name,
		csp.classroom_academic_id,
		csp.count_plan,
		csp.max_point,
		csp.created_at,
		csp.updated_at
	FROM 
		classroom_subject_plans csp
	INNER JOIN classroom_subjects cs
		ON cs.id = csp.classroom_subject_id
	INNER JOIN subjects s
		ON s.id = cs.subject_id
	WHERE csp.classroom_academic_id=$1`

	res, err = m.fetch(ctx, query, id)
	if err != nil {
		return
	} 

	return
}

func (m csPlanRepo) FetchByTeacher(ctx context.Context, id int64) (res []domain.ClassroomSubjectPlan, err error) {
	query := `
	SELECT 
		csp.id,
		csp.type, 
		csp.name, 
		csp.description, 
		csp.teacher_id,
		csp.classroom_subject_id,
		s.name,
		csp.classroom_academic_id,
		csp.count_plan,
		csp.max_point,
		csp.created_at,
		csp.updated_at
	FROM 
		classroom_subject_plans csp
	INNER JOIN classroom_subjects cs
		ON cs.id = csp.classroom_subject_id
	INNER JOIN subjects s
		ON s.id = cs.subject_id
	WHERE csp.teacher_id=$1`

	res, err = m.fetch(ctx, query, id)
	if err != nil {
		return
	} 

	return
}

func (m csPlanRepo) FetchByTeacherAndClassroom(ctx context.Context, tid int64, cid int64) (res []domain.ClassroomSubjectPlan, err error) {
	query := `
	SELECT 
		csp.id,
		csp.type, 
		csp.name, 
		csp.description, 
		csp.teacher_id,
		csp.classroom_subject_id,
		s.name,
		csp.classroom_academic_id,
		csp.count_plan,
		csp.max_point,
		csp.created_at,
		csp.updated_at
	FROM 
		classroom_subject_plans csp
	INNER JOIN classroom_subjects cs
		ON cs.id = csp.classroom_subject_id
	INNER JOIN subjects s
		ON s.id = cs.subject_id
	WHERE csp.teacher_id=$1 AND csp.classroom_academic_id=$2`

	res, err = m.fetch(ctx, query, tid, cid)
	if err != nil {
		return
	} 

	return
}

func (m csPlanRepo) Update(ctx context.Context, csp *domain.ClassroomSubjectPlan) (err error) {
	query := `UPDATE classroom_subject_plans SET type=$1, name=$2, description=$3, teacher_id=$4, classroom_subject_id=$5, classroom_academic_id=$6, count_plan=$7, max_point=$8, updated_at=$9 WHERE id=$10`

	result, err := m.Conn.ExecContext(ctx, query, 
		csp.Type,
		csp.Name,
		csp.Desc,
		csp.Teacher.ID,
		csp.Subject.ID,
		csp.Classroom.ID,
		csp.CountPlan,
		csp.MaxPoint,
		csp.UpdatedAt,
		csp.ID,
	)

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

func (m csPlanRepo) Delete(ctx context.Context, id int64) (err error) {
	query := `DELETE FROM classroom_subject_plans WHERE id=$1`
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