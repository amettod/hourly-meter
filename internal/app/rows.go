package app

import (
	"bytes"
	"errors"
	"strconv"
	"time"
)

// Row all values on a single row
type Row struct {
	ID     int
	PPlus  float64
	PMinus float64
	QPlus  float64
	QMinus float64
	Period string
	Note   string
	Date   time.Time
}

// parse analyzes row by row and finds values
func parse(data []byte) ([]*Row, string, error) {
	if len(data) == 0 {
		return nil, "", errors.New("bad data: shouldn't be empty")
	}

	var rows []*Row
	var meter string

	// ten values per row
	var tenRows []string

	for _, row := range bytes.Split(data, []byte("\n")) {
		if bytes.Contains(row, []byte("TD")) { // checks the content of the tag
			tenRows = append(tenRows, value(row))
			if len(tenRows) == 10 {
				r, err := toRows(tenRows)
				if err != nil {
					return nil, "", err
				}
				rows = append(rows, r)
				tenRows = nil
			}
		} else if bytes.Contains(row, []byte("H2")) && meter == "" {
			meter = serial(row)
		}
	}

	return rows, meter, nil
}

// toRows cast type
func toRows(tenRows []string) (*Row, error) {
	id, err := strconv.Atoi(tenRows[0])
	if err != nil {
		return nil, err
	}

	pPlus, err := strconv.ParseFloat(tenRows[1], 64)
	if err != nil {
		return nil, err
	}

	pMinus, err := strconv.ParseFloat(tenRows[2], 64)
	if err != nil {
		return nil, err
	}

	qPlus, err := strconv.ParseFloat(tenRows[3], 64)
	if err != nil {
		return nil, err
	}

	qMinus, err := strconv.ParseFloat(tenRows[4], 64)
	if err != nil {
		return nil, err
	}

	if len(tenRows[9]) < 10 {
		return nil, errors.New("bad value: date length must be more than 10")
	}
	date, err := strconv.ParseInt(tenRows[9][:10], 10, 64)
	if err != nil {
		return nil, err
	}

	return &Row{
		ID:     id,
		PPlus:  pPlus,
		PMinus: pMinus,
		QPlus:  qPlus,
		QMinus: qMinus,
		Period: tenRows[7],
		Note:   tenRows[8],
		Date:   time.Unix(date, 0).Local().UTC()}, nil
}

// serial return the value of the meter serial number
func serial(row []byte) string {
	if startIndex := bytes.LastIndexByte(row, ' '); startIndex != -1 {
		if index := bytes.IndexByte(row[startIndex:], ')'); index != -1 {
			endIndex := startIndex + index
			return string(row[startIndex+1 : endIndex])
		}
	}

	return ""
}

// value return the value from the row
func value(row []byte) string {
	var val string

	if startIndex := bytes.IndexByte(row, '>'); startIndex != -1 {
		if index := bytes.IndexByte(row[startIndex:], '<'); index != -1 {
			endIndex := startIndex + index
			val = string(row[startIndex+1 : endIndex])
		}
	}

	return val
}
