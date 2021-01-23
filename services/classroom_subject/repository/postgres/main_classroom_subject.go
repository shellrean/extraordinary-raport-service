package postgres

import (
	"database/sql"
	"context"
	"fmt"
	"strings"

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

func (m *csuRepo) fetch(ctx context.Context, query string, args ...interface{}) (res []domain.ClassroomSubject, err error) {
    rows, err := m.Conn.QueryContext(ctx, query, args...)
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

	res, err = m.fetch(ctx, query, academicClassroomID)
    if err != nil {
        return nil, err
    }

    return
}

func (m *csuRepo) Store(ctx context.Context, cs *domain.ClassroomSubject) (err error) {
	query := `INSERT INTO classroom_subjects (classroom_academic_id, subject_id, teacher_id, mgn, created_at, updated_at)
		VALUES($1,$2,$3,$4,$5,$6) RETURNING id`

	err = m.Conn.QueryRowContext(ctx, query,
		cs.ClassroomAcademic.ID,
		cs.Subject.ID,
		cs.Teacher.ID,
		cs.MGN,
		cs.CreatedAt,
		cs.UpdatedAt,
	).Scan(&cs.ID)

	if err != nil {
		return
	}

	return
}

func (m *csuRepo) GetByID(ctx context.Context, id int64) (res domain.ClassroomSubject, err error) {
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
	WHERE cs.id=$1`

	list, err := m.fetch(ctx, query, id)
    if err != nil {
        return domain.ClassroomSubject{}, err
    }
    if len(list) < 1 {
        return domain.ClassroomSubject{}, err
    }
    res = list[0]
    return
}

func (m *csuRepo) GetByClassroomAndSubject(ctx context.Context, academicClassroomID int64, subjectID int64) (res domain.ClassroomSubject, err error) {
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
	WHERE cs.classroom_academic_id=$1
		AND cs.subject_id=$2 
	LIMIT 1`

	list, err := m.fetch(ctx, query, academicClassroomID, subjectID)
    if err != nil {
        return domain.ClassroomSubject{}, err
    }
    if len(list) < 1 {
        return domain.ClassroomSubject{}, err
    }
    res = list[0]
    return
}

func (m *csuRepo) Update(ctx context.Context, cs *domain.ClassroomSubject) (err error) {
	query := `UPDATE classroom_subjects SET subject_id=$1, teacher_id=$2, mgn=$3, updated_at=$4 
		WHERE id=$5`
	
	result, err := m.Conn.ExecContext(ctx, query, cs.Subject.ID, cs.Teacher.ID, cs.MGN, cs.UpdatedAt, cs.ID)
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

func (m *csuRepo) Delete(ctx context.Context, id int64) (err error) {
    query := `DELETE FROM classroom_subjects WHERE id=$1`
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

func (m *csuRepo) StoreMultiple(ctx context.Context, cs []domain.ClassroomSubject) (err error) {
	var valueStrings []string
    var valueArgs []interface{}

    for i, item := range cs {
        valueStrings = append(valueStrings, fmt.Sprintf("($%d,$%d,$%d,$%d,$%d, $%d)", (6*i+1), (6*i+2), (6*i+3), (6*i+4), (6*i+5), (6*i+6)))

		valueArgs = append(valueArgs, item.ClassroomAcademic.ID)
		valueArgs = append(valueArgs, item.Subject.ID)
		valueArgs = append(valueArgs, item.Teacher.ID)
		valueArgs = append(valueArgs, item.MGN)
        valueArgs = append(valueArgs, item.CreatedAt)
        valueArgs = append(valueArgs, item.UpdatedAt)
	}
	
	query := `INSERT INTO classroom_subjects (classroom_academic_id, subject_id, teacher_id, mgn, created_at, updated_at)
		VALUES %s`
	query = fmt.Sprintf(query, strings.Join(valueStrings, ","))
	tx, err := m.Conn.Begin()
    if err != nil {
        return
    }

    _, err = tx.ExecContext(ctx, query, valueArgs...)
    if err != nil {
        _ = tx.Rollback()
        return
    }
    
    err = tx.Commit()
    
    return
}