package app

import (
	"errors"
	"os"
	"time"
)

// App all values
type App struct {
	Contract    string
	CompanyName string
	Meter       string
	Coefficient float64
	Rows        []*Row
	FirstDay    int
	FirstHour   int
	DaysInMonth int
	Month       int
	Year        int
	Total       float64
}

// New return App
func New(filename, contract, companyName, meter string, coefficient float64) (*App, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	rows, dataMeter, err := parse(data)
	if err != nil {
		return nil, err
	}

	if meter == "" && dataMeter != "" {
		meter = dataMeter
	}

	firstRow := rows[0]

	return &App{
		Contract:    contract,
		CompanyName: companyName,
		Meter:       meter,
		Coefficient: coefficient,
		Rows:        rows,
		FirstDay:    firstRow.Date.Day(),
		FirstHour:   firstRow.Date.Hour(),
		Month:       int(firstRow.Date.Month()),
		Year:        firstRow.Date.Year(),
		DaysInMonth: daysInMonth(firstRow.Date)}, nil
}

// daysInMonth return days in month at a selected period
func daysInMonth(t time.Time) int {
	return t.AddDate(0, 1, -t.Day()).Day()
}

// Run application start
func (a *App) Run() error {
	if a.DaysInMonth == 0 {
		return errors.New("there should be no 0 days in a month")
	}

	dirName, err := createDir(a.Month, a.Year)
	if err != nil {
		return err
	}

	index := 0

	for i := 1; i <= a.DaysInMonth; i++ {
		daily := make([]float64, 24)

		if i >= a.FirstDay {
			for j := 0; j < len(daily); j++ {
				if j+1 >= a.FirstHour {
					if index < len(a.Rows) {
						p := a.Rows[index].PPlus * a.Coefficient
						a.Total += p

						index++

						daily[j] = p
					}
				}
			}
		}

		date := time.Date(a.Year, time.Month(a.Month), i, 0, 0, 0, 0, time.UTC)
		head := newHead(date, a.Contract, a.CompanyName, a.Meter)

		body, err := newBody(daily)
		if err != nil {
			return err
		}

		buff, err := toBuffer(head, body)
		if err != nil {
			return err
		}

		err = toFile(buff, dirName, head)
		if err != nil {
			return err
		}
	}

	return nil
}
