package app

import (
	"bytes"
	"fmt"
	"os"
	"path"
	"reflect"
	"testing"
	"time"
)

func Test_newHead(t *testing.T) {
	type args struct {
		date        time.Time
		contract    string
		companyName string
		counter     string
	}
	tests := []struct {
		name string
		args args
		want *Head
	}{
		{
			name: "empty head",
			args: args{},
			want: &Head{
				Day:         "01",
				Contract:    "",
				CompanyName: "",
				Meter:       "",
				Month:       "01",
				Year:        "1",
			},
		},
		{
			name: "head witch",
			args: args{
				date:        time.Date(2020, 2, 2, 10, 0, 0, 0, time.UTC),
				contract:    "98765432",
				companyName: "OOO STAR",
				counter:     "23456789",
			},
			want: &Head{
				Day:         "02",
				Contract:    "98765432",
				CompanyName: "OOO STAR",
				Meter:       "23456789",
				Month:       "02",
				Year:        "2020",
			},
		},
		{
			name: "time only",
			args: args{
				date:        time.Date(2020, 11, 17, 2, 0, 0, 0, time.UTC),
				contract:    "",
				companyName: "",
				counter:     "",
			},
			want: &Head{
				Day:         "17",
				Contract:    "",
				CompanyName: "",
				Meter:       "",
				Month:       "11",
				Year:        "2020",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newHead(tt.args.date, tt.args.contract, tt.args.companyName, tt.args.counter); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newHead() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_newBody(t *testing.T) {
	type args struct {
		dailyValues []float64
	}
	tests := []struct {
		name    string
		args    args
		want    []*Body
		wantErr bool
	}{
		{
			name: "dab data",
			args: args{
				dailyValues: []float64{0, 0},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "empty data",
			args: args{
				dailyValues: make([]float64, 24),
			},
			want: []*Body{
				{
					StartPeriod: "00",
					EndPeriod:   "01",
					Value:       "0",
				},
				{
					StartPeriod: "01",
					EndPeriod:   "02",
					Value:       "0",
				},
				{
					StartPeriod: "02",
					EndPeriod:   "03",
					Value:       "0",
				},
				{
					StartPeriod: "03",
					EndPeriod:   "04",
					Value:       "0",
				},
				{
					StartPeriod: "04",
					EndPeriod:   "05",
					Value:       "0",
				},
				{
					StartPeriod: "05",
					EndPeriod:   "06",
					Value:       "0",
				},
				{
					StartPeriod: "06",
					EndPeriod:   "07",
					Value:       "0",
				},
				{
					StartPeriod: "07",
					EndPeriod:   "08",
					Value:       "0",
				},
				{
					StartPeriod: "08",
					EndPeriod:   "09",
					Value:       "0",
				},
				{
					StartPeriod: "09",
					EndPeriod:   "10",
					Value:       "0",
				},
				{
					StartPeriod: "10",
					EndPeriod:   "11",
					Value:       "0",
				},
				{
					StartPeriod: "11",
					EndPeriod:   "12",
					Value:       "0",
				},
				{
					StartPeriod: "12",
					EndPeriod:   "13",
					Value:       "0",
				},
				{
					StartPeriod: "13",
					EndPeriod:   "14",
					Value:       "0",
				},
				{
					StartPeriod: "14",
					EndPeriod:   "15",
					Value:       "0",
				},
				{
					StartPeriod: "15",
					EndPeriod:   "16",
					Value:       "0",
				},
				{
					StartPeriod: "16",
					EndPeriod:   "17",
					Value:       "0",
				},
				{
					StartPeriod: "17",
					EndPeriod:   "18",
					Value:       "0",
				},
				{
					StartPeriod: "18",
					EndPeriod:   "19",
					Value:       "0",
				},
				{
					StartPeriod: "19",
					EndPeriod:   "20",
					Value:       "0",
				},
				{
					StartPeriod: "20",
					EndPeriod:   "21",
					Value:       "0",
				},
				{
					StartPeriod: "21",
					EndPeriod:   "22",
					Value:       "0",
				},
				{
					StartPeriod: "22",
					EndPeriod:   "23",
					Value:       "0",
				},
				{
					StartPeriod: "23",
					EndPeriod:   "00",
					Value:       "0",
				}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := newBody(tt.args.dailyValues)
			if (err != nil) != tt.wantErr {
				t.Errorf("newBody() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newBody() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_createDir(t *testing.T) {
	type args struct {
		month int
		year  int
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "create dir 80020-02-2020",
			args: args{
				month: 2,
				year:  2020,
			},
			want:    "80020-02-2020",
			wantErr: false,
		},
		{
			name: "does not create a directory if it exists",
			args: args{
				month: 2,
				year:  2020,
			},
			want:    "80020-02-2020",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := createDir(tt.args.month, tt.args.year)
			if (err != nil) != tt.wantErr {
				t.Errorf("createDir() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("createDir() = %v, want %v", got, tt.want)
			}
		})
	}

	err := os.Remove(tests[0].want)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_toBuffer(t *testing.T) {
	head := &Head{
		Day:         "02",
		Contract:    "98765432",
		CompanyName: "OOO STAR",
		Meter:       "23456789",
		Month:       "02",
		Year:        "2020",
	}
	body := []*Body{
		{
			StartPeriod: "00",
			EndPeriod:   "01",
			Value:       "45.0",
		},
		{
			StartPeriod: "01",
			EndPeriod:   "02",
			Value:       "38.1",
		},
		{
			StartPeriod: "02",
			EndPeriod:   "03",
			Value:       "49.3",
		},
		{
			StartPeriod: "03",
			EndPeriod:   "04",
			Value:       "56.7",
		},
		{
			StartPeriod: "04",
			EndPeriod:   "05",
			Value:       "34.2",
		},
		{
			StartPeriod: "05",
			EndPeriod:   "06",
			Value:       "43.12",
		},
		{
			StartPeriod: "06",
			EndPeriod:   "07",
			Value:       "24.90",
		},
		{
			StartPeriod: "07",
			EndPeriod:   "08",
			Value:       "23.4",
		},
		{
			StartPeriod: "08",
			EndPeriod:   "09",
			Value:       "24.23",
		},
		{
			StartPeriod: "09",
			EndPeriod:   "10",
			Value:       "65.24",
		},
		{
			StartPeriod: "10",
			EndPeriod:   "11",
			Value:       "65.13",
		},
		{
			StartPeriod: "11",
			EndPeriod:   "12",
			Value:       "84.23",
		},
		{
			StartPeriod: "12",
			EndPeriod:   "13",
			Value:       "32.9",
		},
		{
			StartPeriod: "13",
			EndPeriod:   "14",
			Value:       "23.5",
		},
		{
			StartPeriod: "14",
			EndPeriod:   "15",
			Value:       "45.2",
		},
		{
			StartPeriod: "15",
			EndPeriod:   "16",
			Value:       "34.56",
		},
		{
			StartPeriod: "16",
			EndPeriod:   "17",
			Value:       "43.85",
		},
		{
			StartPeriod: "17",
			EndPeriod:   "18",
			Value:       "30.43",
		},
		{
			StartPeriod: "18",
			EndPeriod:   "19",
			Value:       "94.32",
		},
		{
			StartPeriod: "19",
			EndPeriod:   "20",
			Value:       "34.23",
		},
		{
			StartPeriod: "20",
			EndPeriod:   "21",
			Value:       "45.12",
		},
		{
			StartPeriod: "21",
			EndPeriod:   "22",
			Value:       "23.5",
		},
		{
			StartPeriod: "22",
			EndPeriod:   "23",
			Value:       "93.3",
		},
		{
			StartPeriod: "23",
			EndPeriod:   "00",
			Value:       "3.2",
		},
	}

	got, err := toBuffer(head, body)
	if err != nil {
		t.Errorf("CreateOneDayBuffer() error = %v", err)
	}

	tests := [][]byte{
		[]byte(`<timestamp>20200202235555</timestamp>`),
		[]byte(`<day>20200202</day>`),
		[]byte(`<name>OOO STAR</name>`),
		[]byte(`<accountpoint code="9876543223456789" name="">`),
		[]byte(`<value status="0">24.23</value>`),
		[]byte(`<period start="2100" end="2200">`),
		[]byte(`</message>`),
	}

	for _, tt := range tests {
		t.Run(string(tt), func(t *testing.T) {
			if !bytes.Contains(got.Bytes(), tt) {
				t.Errorf("got don't contain %v", tt)
			}
		})
	}
}

func Test_toFile(t *testing.T) {
	type args struct {
		buff    *bytes.Buffer
		dirName string
		h       *Head
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "create file empty args",
			args: args{
				buff:    &bytes.Buffer{},
				dirName: "",
				h:       &Head{},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := toFile(tt.args.buff, tt.args.dirName, tt.args.h); (err != nil) != tt.wantErr {
				t.Errorf("toFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			filename := fmt.Sprintf("80020_001_%s_%02s%02s%s.xml", tt.args.h.Contract, tt.args.h.Day, tt.args.h.Month, tt.args.h.Year)
			if _, err := os.Stat(path.Join(tt.args.dirName, filename)); os.IsNotExist(err) {
				t.Errorf("file %s not created", filename)
			}

			err := os.Remove(path.Join(tt.args.dirName, filename))
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}
