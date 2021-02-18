package postgres

import (
	"fmt"
	"context"
	"database/sql"

	"github.com/shellrean/extraordinary-raport/domain"
)

type repository struct {
	Conn *sql.DB
}

func New(Conn *sql.DB) domain.StudentRepository {
	return &repository{
		Conn,
	}
}

func (m *repository) fetch(ctx context.Context, query string, args ...interface{}) (result []domain.Student, err error) {
	rows, err := m.Conn.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	defer func(){
		rows.Close()
	}()

	result = make([]domain.Student, 0)
	for rows.Next() {
		t := domain.Student{}
		religionID := int64(0)
		err = rows.Scan(
			&t.ID,
			&t.SRN,
			&t.NSRN,
			&t.Name,
			&t.Gender,
			&t.BirthPlace,
			&t.BirthDate,
			&religionID,
			&t.Address,
			&t.Telp,
			&t.SchoolBefore,
			&t.AcceptedGrade,
			&t.AcceptedDate,
			&t.Familly.Status,
			&t.Familly.Order,
			&t.Father.Name,
			&t.Father.Address,
			&t.Father.Profession,
			&t.Father.Telp,
			&t.Mother.Name,
			&t.Mother.Address,
			&t.Mother.Profession,
			&t.Mother.Telp,
			&t.Guardian.Name,
			&t.Guardian.Address,
			&t.Guardian.Profession,
			&t.Guardian.Telp,
			&t.CreatedAt,
			&t.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}
		t.Religion = domain.Religion{
			ID: religionID,
		}
		result = append(result, t)
	}

	return result, nil
}

func (m *repository) Fetch(ctx context.Context, q string, cursor int64, num int64) (res []domain.Student, err error) {
	query := `SELECT id, srn, nsrn, name, gender, birth_place, birth_date, religion_id, address, telp,
			school_before, accepted_grade, accepted_date, familly_status, familly_order, father_name,
			father_address, father_profession, father_telp, mother_name, mother_address, mother_profession,
			mother_telp, guardian_name, guardian_address, guardian_profession, guardian_telp, created_at, updated_at
			FROM students WHERE (LOWER(name) LIKE '%' || $1 || '%' OR srn LIKE '%' || $2 || '%') AND id > $3 LIMIT $4`
	res, err = m.fetch(ctx, query, q, q, cursor, num)
	if err != nil {
		return nil, err
	}

	return
}

func (m *repository) GetByID(ctx context.Context, id int64) (res domain.Student, err error) {
	query := `SELECT id, srn, nsrn, name, gender, birth_place, birth_date, religion_id, address, telp,
			school_before, accepted_grade, accepted_date, familly_status, familly_order, father_name,
			father_address, father_profession, father_telp, mother_name, mother_address, mother_profession,
			mother_telp, guardian_name, guardian_address, guardian_profession, guardian_telp, created_at, updated_at
			FROM students WHERE id = $1`

    list, err := m.fetch(ctx, query, id)
    if err != nil {
        return domain.Student{}, err
    }
    if len(list) < 1 {
        return domain.Student{}, err
    }
    res = list[0]
    return
}

func (m *repository) Store(ctx context.Context, t *domain.Student) (err error) {
	query := `INSERT INTO students (srn,nsrn,name,gender,birth_place,birth_date,religion_id,address,telp,
		school_before,accepted_grade,accepted_date,familly_status,familly_order,father_name,father_address,father_profession,
		father_telp,mother_name,mother_address,mother_profession,mother_telp,guardian_name,guardian_address,guardian_profession,
		guardian_telp,created_at,updated_at) VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20,$21,$22,$23,$24,$25,$26,$27,$28)  RETURNING id`
	err = m.Conn.QueryRowContext(ctx, query, 
		t.SRN,
		t.NSRN,
		t.Name,
		t.Gender,
		t.BirthPlace,
		t.BirthDate,
		t.Religion.ID,
		t.Address,
		t.Telp,
		t.SchoolBefore,
		t.AcceptedGrade,
		t.AcceptedDate,
		t.Familly.Status,
		t.Familly.Order,
		t.Father.Name,
		t.Father.Address,
		t.Father.Profession,
		t.Father.Telp,
		t.Mother.Name,
		t.Mother.Address,
		t.Mother.Profession,
		t.Mother.Telp,
		t.Guardian.Name,
		t.Guardian.Address,
		t.Guardian.Profession,
		t.Guardian.Telp,
		t.CreatedAt,
		t.UpdatedAt).Scan(&t.ID)
	if err != nil {
		return
	}
	return
}

func (m *repository) Update(ctx context.Context, t *domain.Student) (error) {
	query := `UPDATE students SET (srn,nsrn,name,gender,birth_place,birth_date,religion_id,address,telp,
		school_before,accepted_grade,accepted_date,familly_status,familly_order,father_name,father_address,father_profession,
		father_telp,mother_name,mother_address,mother_profession,mother_telp,guardian_name,guardian_address,guardian_profession,
		guardian_telp,updated_at) = ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20,$21,$22,$23,$24,$25,$26,$27) WHERE id=$28`
	result, err := m.Conn.ExecContext(ctx, query, 
		t.SRN,
		t.NSRN,
		t.Name,
		t.Gender,
		t.BirthPlace,
		t.BirthDate,
		t.Religion.ID,
		t.Address,
		t.Telp,
		t.SchoolBefore,
		t.AcceptedGrade,
		t.AcceptedDate,
		t.Familly.Status,
		t.Familly.Order,
		t.Father.Name,
		t.Father.Address,
		t.Father.Profession,
		t.Father.Telp,
		t.Mother.Name,
		t.Mother.Address,
		t.Mother.Profession,
		t.Mother.Telp,
		t.Guardian.Name,
		t.Guardian.Address,
		t.Guardian.Profession,
		t.Guardian.Telp,
		t.UpdatedAt,
		t.ID,
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
    return nil
}

func (m *repository) Delete(ctx context.Context, id int64) (err error) {
	query := `DELETE FROM students WHERE id=$1`
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