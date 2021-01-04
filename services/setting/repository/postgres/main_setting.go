package postgres

import (
	"database/sql"
	"context"

	"github.com/shellrean/extraordinary-raport/domain"
)

type settingRepository struct {
	Conn *sql.DB
}

func NewPostgresSettingRepository(Conn *sql.DB) domain.SettingRepository {
	return &settingRepository{
		Conn,
	}
}

func (m *settingRepository) GetByName(ctx context.Context, name string) (res domain.Setting, err error) {
	query := `SELECT id,name,value,created_at,updated_at 
		FROM settings WHERE name=$1`
	err = m.Conn.QueryRowContext(ctx, query, name).
		Scan(
			&res.ID,
			&res.Name,
			&res.Value,
			&res.CreatedAt,
			&res.UpdatedAt,
		)
	if err != nil {
		return
	}
	return
}