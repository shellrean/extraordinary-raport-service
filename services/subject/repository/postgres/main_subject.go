package postgres

import (
    "fmt"
	"context"
	"database/sql"

	"github.com/shellrean/extraordinary-raport/domain"
)

type subjectRepository struct{
	Conn *sql.DB
}

func NewPostgresSubjectRepository(Conn *sql.DB) domain.SubjectRepository {
	return &subjectRepository{
		Conn,
	}
}

func (m *subjectRepository) fetch(ctx context.Context, query string, args ...interface{}) (result []domain.Subject, err error) {
    rows, err := m.Conn.QueryContext(ctx, query, args...)
    if err != nil {
        return nil, err
    }

    defer func() {
        rows.Close()
    }()

    for rows.Next() {
        t := domain.Subject{}
        err = rows.Scan(
            &t.ID,
			&t.Name,
			&t.Type,
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

func (m *subjectRepository) Fetch(ctx context.Context, cursor int64, num int64) (res []domain.Subject, err error) {
	query := `SELECT id,name,type,created_at,updated_at
		FROM subjects WHERE id > $1 ORDER BY id LIMIT $2`
	
	res, err = m.fetch(ctx, query, cursor, num)
	if err != nil {
		return nil, err
	}
	return
}

func (m *subjectRepository) GetByID(ctx context.Context, id int64) (res domain.Subject, err error) {
    query := `SELECT id,name,type,created_at,updated_at
        FROM subjects WHERE id = $1`

    list, err := m.fetch(ctx, query, id)
    if err != nil {
        return domain.Subject{}, err
    }
    if len(list) < 1 {
        return domain.Subject{}, err
    }
    res = list[0]
    return
}

func (m *subjectRepository) Store(ctx context.Context, s *domain.Subject) (err error) {
    query := `INSERT INTO subjects (name,type,created_at,updated_at)
        VALUES($1,$2,$3,$4) RETURNING id`
    err = m.Conn.QueryRowContext(ctx, query, 
        s.Name,
        s.Type,
        s.CreatedAt,
        s.UpdatedAt,
    ).Scan(&s.ID)
    if err != nil {
        return
    }
    return
}

func (m *subjectRepository) Update(ctx context.Context, s *domain.Subject) (err error) {
    query := `UPDATE subjects SET name=$1, type=$2, updated_at=$3 WHERE id=$4`
    
    result, err := m.Conn.ExecContext(ctx, query, s.Name, s.Type, s.UpdatedAt, s.ID)
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

func (m *subjectRepository) Delete(ctx context.Context, id int64) (err error) {
    query := `DELETE FROM subjects WHERE id=$1`
    
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