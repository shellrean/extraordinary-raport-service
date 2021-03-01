package postgres

import (
	"context"
	"database/sql"

	"github.com/shellrean/extraordinary-raport/domain"
)

type repository struct {
	Conn *sql.DB
}

func New(Conn *sql.DB) domain.StudentNoteRepository {
	return &repository {
		Conn,
	}
}

func (m *repository) Store(ctx context.Context, sn *domain.StudentNote) (err error) {
	query := `INSERT INTO student_notes (type, created_by, student_id, note, created_at, updated_at)
		VALUES($1, $2, $3, $4, $5, $6) RETURNING id`
	err = m.Conn.QueryRowContext(ctx, query,
		sn.Type,
		sn.Teacher.ID,
		sn.Student.ID,
		sn.Note,
		sn.CreatedAt,
		sn.UpdatedAt,
	).Scan(&sn.ID)
	if err != nil {
		return
	}
	return
}