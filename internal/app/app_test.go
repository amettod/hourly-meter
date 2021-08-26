package app

import (
	"fmt"
	"os"
	"reflect"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	type args struct {
		filename    string
		contract    string
		companyName string
		counter     string
		coefficient float64
	}
	tests := []struct {
		name    string
		args    args
		want    *App
		wantErr bool
	}{
		{
			name:    "empty args",
			args:    args{},
			want:    nil,
			wantErr: true,
		},
		{
			name: "empty file",
			args: args{
				filename:    "testdata/empty.html",
				contract:    "98765432",
				companyName: "OOO STAR",
				counter:     "23456789",
				coefficient: 4000,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "first row",
			args: args{
				filename:    "testdata/first_row.html",
				contract:    "98765432",
				companyName: "OOO STAR",
				counter:     "",
				coefficient: 4000,
			},
			want: &App{
				Contract:    "98765432",
				CompanyName: "OOO STAR",
				Meter:       "23456789",
				Coefficient: 4000,
				Rows: []*Row{
					{
						ID:     2,
						PPlus:  0.0705,
						PMinus: 0.0000,
						QPlus:  0.0063,
						QMinus: 0.0019,
						Period: "30+30",
						Note:   "-",
						Date:   time.Date(2020, 2, 1, 1, 0, 0, 0, time.UTC),
					},
				},
				FirstDay:    1,
				FirstHour:   1,
				DaysInMonth: 29,
				Month:       2,
				Year:        2020,
				Total:       0,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.args.filename, tt.args.contract, tt.args.companyName, tt.args.counter, tt.args.coefficient)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestApp_Run(t *testing.T) {
	type fields struct {
		Contract    string
		CompanyName string
		Counter     string
		Coefficient float64
		Rows        []*Row
		FirstDay    int
		FirstHour   int
		DaysInMonth int
		Month       int
		Year        int
		Total       float64
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name:    "empty fields",
			fields:  fields{},
			wantErr: true,
		},
		{
			name: "first row",
			fields: fields{
				Contract:    "98765432",
				CompanyName: "OOO STAR",
				Counter:     "23456789",
				Coefficient: 4000,
				Rows: []*Row{
					{
						ID:     2,
						PPlus:  0.0705,
						PMinus: 0.0000,
						QPlus:  0.0063,
						QMinus: 0.0019,
						Period: "30+30",
						Note:   "-",
						Date:   time.Date(2020, 2, 1, 1, 0, 0, 0, time.UTC),
					},
				},
				FirstDay:    1,
				FirstHour:   1,
				DaysInMonth: 29,
				Month:       2,
				Year:        2020,
				Total:       0,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &App{
				Contract:    tt.fields.Contract,
				CompanyName: tt.fields.CompanyName,
				Meter:       tt.fields.Counter,
				Coefficient: tt.fields.Coefficient,
				Rows:        tt.fields.Rows,
				FirstDay:    tt.fields.FirstDay,
				FirstHour:   tt.fields.FirstHour,
				DaysInMonth: tt.fields.DaysInMonth,
				Month:       tt.fields.Month,
				Year:        tt.fields.Year,
				Total:       tt.fields.Total,
			}
			if err := a.Run(); (err != nil) != tt.wantErr {
				t.Errorf("App.Run() error = %v, wantErr %v", err, tt.wantErr)
			}

			dirName := fmt.Sprintf("80020-%02d-%d", a.Month, a.Year)

			err := os.RemoveAll(dirName)
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}
