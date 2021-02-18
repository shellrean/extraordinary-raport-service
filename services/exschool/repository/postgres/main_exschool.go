package postgres

import (
	"context"
	"database/sql"

	"github.com/shellrean/extraordinary-raport/domain"
)

type repository struct{
	Conn *sql.DB
}

func New(Conn *sql.DB) domain.ExschoolRepository{
	return &repository{
		Conn,
	}
}

func (m *repository) Fetch(ctx context.Context) {
	return
}

func (m *repository) Store(ctx context.Context ex *domain.Exschool) (err error) {
	return
}