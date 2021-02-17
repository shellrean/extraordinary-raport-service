package postgres

import (
	"context"
	"database/sql"

	"github.com/shellrean/extraordinary-raport/domain"
)

type exschoolRepo struct{
	Conn *sql.DB
}

func NewPostgresExschoolRepository(Conn *sql.DB) domain.ExschoolRepository{
	return &ExschoolRepo{
		Conn,
	}
}

func (m exschoolRepo) Fetch(ctx context.Context) {
	return
}

func (m exschoolRepo) Store(ctx context.Context ex *domain.Exschool) (err error) {
	return
}