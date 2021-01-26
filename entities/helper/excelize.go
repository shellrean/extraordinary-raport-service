package helper

import (
	"fmt"
	"time"
	"context"

	"golang.org/x/crypto/bcrypt"
	"golang.org/x/sync/errgroup"
	"github.com/360EntSecGroup-Skylar/excelize"

	"github.com/shellrean/extraordinary-raport/domain"
)

func ReadUserFileExcel(c context.Context, file string) (res []domain.User, err error) {
	xlsx, err := excelize.OpenFile(file)
	if err != nil {
		err = fmt.Errorf("Error when open file")
		return 
	}

	users := "users"

	rows := xlsx.GetRows(users)

	g, _ := errgroup.WithContext(c)
	chanUser := make(chan domain.User)
	for i, row := range rows {
		row := row
		i := i
		g.Go(func() error {
			password, err := bcrypt.GenerateFromPassword([]byte(row[2]), 10)
			if err != nil {
				return fmt.Errorf("Error when generate password student at row %d",i)
			}
			user := domain.User{
				Name: row[0],
				Email: row[1],
				Password: string(password),
				Role: domain.RoleTeacher,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}
			chanUser <- user
			return nil
		})
	}

	go func() {
		err := g.Wait() 
		if err != nil {
			return
		}
		close(chanUser)
	}()

	for user := range chanUser {
		res = append(res, user)
	}

	if err := g.Wait(); err != nil {
		return nil, err
	}

	return
}