package postgres

import (
    "fmt"
    "context"
    "database/sql"
    "strings"

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
            &t.Student.SRN,
            &t.Student.NSRN,
            &t.Student.Name,
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
    query := `SELECT 
        cs.id,
        cs.classroom_academic_id,
        cs.student_id,
        s.srn,
        s.nsrn,
        s.name,
        cs.created_at,
        cs.updated_at
    FROM 
        classroom_students cs
    INNER JOIN students s
        ON s.id = cs.student_id
    WHERE cs.id > $1 ORDER BY cs.id LIMIT $2`

    res, err = m.fetch(ctx, query, cursor, num)
    if err != nil {
        return nil, err
    }

    return
}

func (m *csRepository) GetByID(ctx context.Context, id int64) (res domain.ClassroomStudent, err error) {
    query := `SELECT 
        cs.id,
        cs.classroom_academic_id,
        cs.student_id,
        s.srn,
        s.nsrn,
        s.name,
        cs.created_at,
        cs.updated_at
    FROM 
        classroom_students cs
    INNER JOIN students s
        ON s.id = cs.student_id
    WHERE cs.id=$1`

    list, err := m.fetch(ctx, query, id)
    if err != nil {
        return domain.ClassroomStudent{}, err
    }
    if len(list) < 1 {
        return domain.ClassroomStudent{}, err
    }
    res = list[0]
    return
}

func (m *csRepository) GetByAcademicAndStudent(ctx context.Context, academic int64, student int64) (res domain.ClassroomStudent, err error) {
    query := `SELECT 
        cs.id,
        cs.classroom_academic_id,
        cs.student_id,
        s.srn,
        s.nsrn,
        s.name,
        cs.created_at,
        cs.updated_at 
    FROM 
        classroom_students cs
    INNER JOIN students s
        ON s.id = cs.student_id
    WHERE classroom_academic_id 
        IN (SELECT id FROM classroom_academics WHERE academic_id=$1)
    AND student_id=$2`

    list, err := m.fetch(ctx, query, academic, student)
    if err != nil {
        return domain.ClassroomStudent{}, err
    }
    if len(list) < 1 {
        return domain.ClassroomStudent{}, err
    }
    res = list[0]
    return
}

func (m *csRepository) GetByClassroomAcademic(ctx context.Context, academicID int64) (res []domain.ClassroomStudent, err error) {
    query := `SELECT 
        cs.id,
        cs.classroom_academic_id,
        cs.student_id,
        s.srn,
        s.nsrn,
        s.name,
        cs.created_at,
        cs.updated_at 
    FROM 
        classroom_students cs
    INNER JOIN students s
        ON s.id = cs.student_id
    WHERE classroom_academic_id=$1`

    res, err = m.fetch(ctx, query, academicID)
    if err != nil {
        return nil, err
    }

    return
}

func (m *csRepository) Store(ctx context.Context, cs *domain.ClassroomStudent) (err error) {
    query := `INSERT INTO classroom_students (classroom_academic_id, student_id,created_at, updated_at)
            VALUES ($1,$2,$3,$4) RETURNING id`
        
    err = m.Conn.QueryRowContext(ctx, query, cs.ClassroomAcademic.ID, cs.Student.ID,cs.CreatedAt, cs.UpdatedAt).Scan(&cs.ID)
    if err != nil {
        return err
    }
    
    return
}

func (m *csRepository) Update(ctx context.Context, cs *domain.ClassroomStudent) (err error) {
    query := `UPDATE classroom_students SET classroom_academic_id=$1, student_id=$2, updated_at=$3
        WHERE id=$4`
    
    result, err := m.Conn.ExecContext(ctx, query, cs.ClassroomAcademic.ID, cs.Student.ID, cs.UpdatedAt, cs.ID)
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

func (m *csRepository) Delete(ctx context.Context, id int64) (err error) {
    query := `DELETE FROM classroom_students WHERE id=$1`
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

func (m *csRepository) StoreMultiple(ctx context.Context, cas []domain.ClassroomStudent) (err error) {
    var valueStrings []string
    var valueArgs []interface{}

    for i, item := range cas {
        valueStrings = append(valueStrings, fmt.Sprintf("($%d,$%d,$%d,$%d)", (4*i+1), (4*i+2), (4*i+3), (4*i+4)))

        valueArgs = append(valueArgs, item.ClassroomAcademic.ID)
        valueArgs = append(valueArgs, item.Student.ID)
        valueArgs = append(valueArgs, item.CreatedAt)
        valueArgs = append(valueArgs, item.UpdatedAt)
    }
    query := `INSERT INTO classroom_students (classroom_academic_id, student_id,created_at, updated_at)
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