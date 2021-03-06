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

func New(Conn *sql.DB) domain.ClassroomAcademicRepository {
	return &repository {
		Conn,
	}
}

func (m *repository) fetch(ctx context.Context, query string, args ...interface{}) (res []domain.ClassroomAcademic, err error) {
    rows, err := m.Conn.QueryContext(ctx, query, args...)
    if err != nil {
        return nil, err
    }

    defer func() {
        rows.Close()
    }()

    for rows.Next() {
		t := domain.ClassroomAcademic{
			Academic: domain.Academic{},
			Classroom: domain.Classroom{Major: domain.Major{}},
			Teacher: domain.User{},
		}
        err = rows.Scan(
            &t.ID,
			&t.Academic.ID,
			&t.Classroom.ID,
			&t.Teacher.ID,
			&t.Classroom.Name,
			&t.Classroom.Major.Name,
			&t.Teacher.Name,
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

func (m *repository) Fetch(ctx context.Context, academicID int64) (res []domain.ClassroomAcademic, err error) {
	query := `SELECT
		ca.id,
		ca.academic_id,
		ca.classroom_id,
		ca.teacher_id,
		c.name,
		m.name,
		u.name,
		ca.created_at,
		ca.updated_at
	FROM 
		classroom_academics ca
	INNER JOIN classrooms c
		ON c.id = ca.classroom_id
	INNER JOIN users u
		ON u.id = ca.teacher_id
	INNER JOIN majors m
		on m.id = c.major_id
	WHERE ca.academic_id=$1`
	
	res, err = m.fetch(ctx, query, academicID)
    if err != nil {
        return nil, err
    }

    return
}

func (m *repository) FetchByAcademicAndTeacher(ctx context.Context, academicID int64, userID int64) (res []domain.ClassroomAcademic, err error) {
	query := `SELECT
		ca.id,
		ca.academic_id,
		ca.classroom_id,
		ca.teacher_id,
		c.name,
		m.name,
		u.name,
		ca.created_at,
		ca.updated_at
	FROM 
		classroom_academics ca
	INNER JOIN classrooms c
		ON c.id = ca.classroom_id
	INNER JOIN users u
		ON u.id = ca.teacher_id
	INNER JOIN majors m
		on m.id = c.major_id
	WHERE ca.academic_id=$1 AND ca.teacher_id=$2`
	
	res, err = m.fetch(ctx, query, academicID, userID)
    if err != nil {
        return nil, err
    }

    return
}

func (m *repository) Store(ctx context.Context, ca *domain.ClassroomAcademic) (err error) {
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

func (m *repository) GetByAcademicAndClass(ctx context.Context, a int64, c int64) (res domain.ClassroomAcademic, err error) {
	query := `SELECT
		ca.id,
		ca.academic_id,
		ca.classroom_id,
		ca.teacher_id,
		c.name,
		m.name,
		u.name,
		ca.created_at,
		ca.updated_at
	FROM 
		classroom_academics ca
	INNER JOIN classrooms c
		ON c.id = ca.classroom_id
	INNER JOIN users u
		ON u.id = ca.teacher_id
	INNER JOIN majors m
		on m.id = c.major_id
	WHERE ca.academic_id=$1 AND ca.classroom_id=$2`

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

func (m *repository) GetByID(ctx context.Context, id int64) (res domain.ClassroomAcademic, err error) {
	query := `SELECT
		ca.id,
		ca.academic_id,
		ca.classroom_id,
		ca.teacher_id,
		c.name,
		m.name,
		u.name,
		ca.created_at,
		ca.updated_at
	FROM 
		classroom_academics ca
	INNER JOIN classrooms c
		ON c.id = ca.classroom_id
	INNER JOIN users u
		ON u.id = ca.teacher_id
	INNER JOIN majors m
		on m.id = c.major_id
	WHERE ca.id=$1`

	list, err := m.fetch(ctx, query, id)
    if err != nil {
        return domain.ClassroomAcademic{}, err
    }
    if len(list) < 1 {
        return domain.ClassroomAcademic{}, err
    }
    res = list[0]
    return
}

func (m *repository) Update(ctx context.Context, ca *domain.ClassroomAcademic) (err error) {
	query := `UPDATE classroom_academics SET classroom_id=$1, teacher_id=$2, updated_at=$3
		WHERE id=$4`
	result, err := m.Conn.ExecContext(ctx, query,
		ca.Classroom.ID,
		ca.Teacher.ID,
		ca.UpdatedAt,
		ca.ID,
	)
	if err != nil {
		return
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return
	}

	if rows != 1 {
		return fmt.Errorf("expected single row affected, got %d rows affected", rows)
	}
	return
}

func (m *repository) Delete(ctx context.Context, id int64) (err error) {
	query := `DELETE FROM classroom_academics WHERE id=$1`
	result, err := m.Conn.ExecContext(ctx, query, id)
	if err != nil {
		return
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return
	}

	if rows != 1 {
		return fmt.Errorf("expected single row affected, got %d rows affected", rows) 
	}
	return
}