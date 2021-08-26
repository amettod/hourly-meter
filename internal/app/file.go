package app

import (
	"bytes"
	"embed"
	"errors"
	"fmt"
	"os"
	"path"
	"strings"
	"text/template"
	"time"
)

//go:embed template
var templateFS embed.FS

// Head is used to create a template head
type Head struct {
	Day         string
	Contract    string
	CompanyName string
	Meter       string
	Month       string
	Year        string
}

// Body is used to create a template body
type Body struct {
	StartPeriod string
	EndPeriod   string
	Value       string
}

// newHead return Head
func newHead(date time.Time, contract, companyName, meter string) *Head {
	return &Head{
		Day:         fmt.Sprintf("%02d", date.Day()),
		Contract:    contract,
		CompanyName: companyName,
		Meter:       meter,
		Month:       fmt.Sprintf("%02d", int(date.Month())),
		Year:        fmt.Sprintf("%d", date.Year()),
	}
}

// newBody return slice Body for one day template
func newBody(dailyValues []float64) ([]*Body, error) {
	if len(dailyValues) != 24 {
		return nil, errors.New("bad data: twenty-four number of values required")
	}

	var bodies []*Body

	for i, v := range dailyValues {
		bodies = append(bodies, &Body{
			StartPeriod: fmt.Sprintf("%02d", i),
			EndPeriod:   endPeriod(i),
			Value:       zeroValue(v),
		})
	}
	return bodies, nil
}

// endPeriod return 00 if i == 24
func endPeriod(i int) string {
	if i != 23 {
		return fmt.Sprintf("%02d", i+1)
	}

	return "00"
}

// zeroValue return 0 if v == 0.0
func zeroValue(v float64) string {
	if v != 0 {
		return strings.Replace(fmt.Sprintf("%.1f", v), ".", ",", 1)
	}

	return "0"
}

// toBuffer return one day template
func toBuffer(head *Head, body []*Body) (*bytes.Buffer, error) {
	tmpl, err := template.New("daily").ParseFS(templateFS, "template/daily_xml.tmpl")
	if err != nil {
		return nil, err
	}

	buff := new(bytes.Buffer)

	err = tmpl.ExecuteTemplate(buff, "head", head)
	if err != nil {
		return nil, err
	}

	err = tmpl.ExecuteTemplate(buff, "body", body)
	if err != nil {
		return nil, err
	}

	err = tmpl.ExecuteTemplate(buff, "tail", nil)
	if err != nil {
		return nil, err
	}

	return buff, nil
}

// toFile write one day buffer to file
func toFile(buff *bytes.Buffer, dirName string, h *Head) error {

	newFilename := fmt.Sprintf("80020_001_%s_%02s%02s%s.xml", h.Contract, h.Day, h.Month, h.Year)

	f, err := os.Create(path.Join(dirName, newFilename))
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write(buff.Bytes())
	if err != nil {
		return err
	}

	return nil
}

// createDir create a directory if it does not exist
func createDir(month, year int) (string, error) {
	dirName := fmt.Sprintf("80020-%02d-%d", month, year)

	if _, err := os.Stat(dirName); os.IsNotExist(err) {
		err := os.Mkdir(dirName, 0766)
		if err != nil {
			return "", err
		}
	}

	return dirName, nil
}
