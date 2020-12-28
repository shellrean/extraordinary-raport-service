package postgres

import (
	"context"
	"database/sql"

	"github.com/shellrean/extraordinary-raport/domain"
)

type postgresStudentRepository struct {
	Conn *sql.DB
}

func NewPostgresStudentRepository(Conn *sql.DB) domain.StudentRepository {
	return &postgresStudentRepository{
		Conn,
	}
}

func (m *postgresStudentRepository) fetch(ctx context.Context, query string, args ...interface{}) (result []domain.Student, err error) {
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

func (m *postgresStudentRepository) Fetch(ctx context.Context, cursor int64, num int64) (res []domain.Student, err error) {
	query := `SELECT id, srn, nsrn, name, gender, birth_place, birth_date, religion_id, address, telp,
			school_before, accepted_grade, accepted_date, familly_status, familly_order, father_name,
			father_address, father_profession, father_telp, mother_name, mother_address, mother_profession,
			mother_telp, guardian_name, guardian_address, guardian_profession, guardian_telp, created_at, updated_at
			FROM students WHERE id > $1 LIMIT $2`
	res, err = m.fetch(ctx, query, cursor, num)
	if err != nil {
		return nil, err
	}

	return
}