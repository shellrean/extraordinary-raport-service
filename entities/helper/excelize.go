package helper

import (
	"fmt"
	"time"
	"context"
	"path/filepath"
	"log"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/sync/errgroup"
	"github.com/360EntSecGroup-Skylar/excelize/v2"

	"github.com/shellrean/extraordinary-raport/domain"
)

func ReadUserFileExcel(c context.Context, file string) (res []domain.User, err error) {
	xlsx, err := excelize.OpenFile(file)
	if err != nil {
		err = fmt.Errorf("Error when open file")
		return 
	}

	users := "users"

	rows, _ := xlsx.GetRows(users)

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

func WritePlanResultFileExcel(c context.Context, plan domain.ClassroomSubjectPlan, datas []map[string]interface{}) (string, error) {
	f := excelize.NewFile()

    border, err := f.NewStyle(`{
        "border": [
        {
            "type": "left",
            "color": "#000000",
            "style": 1
        },
        {
            "type": "top", 
            "color": "#000000",
            "style": 1
        },
        {
            "type": "right", 
            "color": "#000000", 
            "style": 1
        },
        {
            "type": "bottom", 
            "color": "#000000", 
            "style": 1
        }]
    }`)

    borderWithRotate, err := f.NewStyle(`{
        "border": [
        {
            "type": "left",
            "color": "#000000",
            "style": 1
        },
        {
            "type": "top", 
            "color": "#000000",
            "style": 1
        },
        {
            "type": "right", 
            "color": "#000000", 
            "style": 1
        },
        {
            "type": "bottom", 
            "color": "#000000", 
            "style": 1
        }],
        "alignment": {
            "text_rotation": 90
        }
    }`)

	const sheet = "Sheet1"

	f.SetCellValue(sheet, "A2", "NO")
	f.SetCellValue(sheet, "B2", "NIS")
	f.SetCellValue(sheet, "C2", "NAMA")
	f.SetColWidth(sheet, "C", "C", 45)
    f.SetRowHeight(sheet, 2, 120)
	f.SetCellValue(sheet, fmt.Sprintf("%s%s", string('D'+int(plan.CountPlan)), "2"), "RATA-RATA")
	f.SetCellValue(sheet, fmt.Sprintf("%s%s", string('D'+int(plan.CountPlan)+1), "2"), "SPARK LINE")
	
	var i int
	for i = 0; i < int(plan.CountPlan); i++ {
		f.SetCellValue(sheet, string('D'+i)+"2", fmt.Sprintf("%s %d", "NILAI", i+1))
	}

	f.SetCellStyle(sheet, "A2", "C2", border)
    f.SetCellStyle(sheet, "D2", string('D'+i)+"2", borderWithRotate)
	f.SetColWidth(sheet, "D", string('D'+i), 5)
	
	var sparkLineLocations []string
	var sparkLineRange []string

	col := 3
	for i, data := range datas {
		f.SetCellValue(sheet, fmt.Sprintf("%s%d","A",col), i+1)
		f.SetCellValue(sheet, fmt.Sprintf("%s%d","B",col), data["nis"])
		f.SetCellValue(sheet, fmt.Sprintf("%s%d","C",col), data["nama"])
		i=0
		dat := data["nilai"].([]uint)
		for i = 0; i < int(plan.CountPlan); i++ {
			f.SetCellValue(sheet, fmt.Sprintf("%s%d", string('D'+i), col), dat[i])
		}
		var total float64 = 0
		for _, value := range dat {
			total += float64(value)
		}
		f.SetCellValue(sheet, fmt.Sprintf("%s%d", string('D'+i), col), fmt.Sprintf("%.2f", total/float64(len(dat))))
		f.SetCellStyle(sheet, fmt.Sprintf("%s%d","A",col), fmt.Sprintf("%s%d", string('D'+i), col), border)

		sparkLineLocations = append(sparkLineLocations, fmt.Sprintf("%s%d", string('D'+i+1), col))
		sparkLineRange = append(sparkLineRange, fmt.Sprintf("%s!%s%d:%s%d", sheet, "D", col, string('D'+i-1), col))
		col++
	}

	f.AddSparkline(sheet, &excelize.SparklineOption{
		Location: sparkLineLocations,
		Range:    sparkLineRange,
		Markers:  true,
	})

	fullPathFile := filepath.Join("storage", "app", "_tmp", uuid.NewString()+"sheet.xlsx")

	if err := f.SaveAs(fullPathFile); err != nil {
		log.Println(err)
		return "", err
    }
	
	return fullPathFile, err
}