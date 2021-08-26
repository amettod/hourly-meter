package app

import (
	"reflect"
	"testing"
	"time"
)

func Test_value(t *testing.T) {
	type args struct {
		row []byte
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "value empty",
			args: args{
				row: []byte(`<TD class=style21></TD>`),
			},
			want: "",
		},
		{
			name: "value PPlus",
			args: args{
				row: []byte(`<TD class=style21>2</TD>`),
			},
			want: "2",
		},
		{
			name: "value UTC",
			args: args{
				row: []byte(`<TD class=style21>1580522400000</TD></TR>`),
			},
			want: "1580522400000",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := value(tt.args.row); got != tt.want {
				t.Errorf("value() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_serial(t *testing.T) {
	type args struct {
		row []byte
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "row empty",
			args: args{
				row: []byte(`"<H2></H2>"`),
			},
			want: "",
		},
		{
			name: "row witch Count",
			args: args{
				row: []byte(`<H2>М234 (Сетевой адрес - 17, Серийный номер - 23456789)</H2>`),
			},
			want: "23456789",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := serial(tt.args.row); got != tt.want {
				t.Errorf("serial() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_toRows(t *testing.T) {
	type args struct {
		tenRows []string
	}
	tests := []struct {
		name    string
		args    args
		want    *Row
		wantErr bool
	}{
		{
			name: "fail id",
			args: args{
				tenRows: []string{"fail"},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "fail pPlus",
			args: args{
				tenRows: []string{"1", "fail"},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "fail pMinus",
			args: args{
				tenRows: []string{"1", "1", "fail"},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "fail qPlus",
			args: args{
				tenRows: []string{"1", "1", "1", "fail"},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "fail qMinus",
			args: args{
				tenRows: []string{"1", "1", "1", "1", "fail"},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "fail date len 4",
			args: args{
				tenRows: []string{"1", "1", "1", "1", "1", "1", "1", "1", "1", "fail"},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "fail date len 10",
			args: args{
				tenRows: []string{"1", "1", "1", "1", "1", "1", "1", "1", "1", "failfailfa"},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "valid tenRows",
			args: args{
				tenRows: []string{
					"4",
					"0.0686",
					"0.0000",
					"0.0018",
					"0.0031",
					"02:00",
					"01.02.20",
					"30+30",
					"-",
					"1580522400000",
				},
			},
			want: &Row{
				ID:     4,
				PPlus:  0.0686,
				PMinus: 0.0000,
				QPlus:  0.0018,
				QMinus: 0.0031,
				Period: "30+30",
				Note:   "-",
				Date:   time.Date(2020, 2, 1, 2, 0, 0, 0, time.UTC),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := toRows(tt.args.tenRows)
			if (err != nil) != tt.wantErr {
				t.Errorf("toRows() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("toRows() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parse(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		args    args
		want    []*Row
		want1   string
		wantErr bool
	}{
		{
			name:    "data empty",
			args:    args{},
			want:    nil,
			want1:   "",
			wantErr: true,
		},
		{
			name: "valid data",
			args: args{
				data: []byte(`<HTML>
									   <HEAD>
									   <META http-equiv=Content-Type content="text/html; charset=windows-1251">
									   <link rel='stylesheet' href='_img/my.css' type='text/css' />
									   <!--[if IE]><SCRIPT language=javascript src="../res/flot/excanvas.min.js" type=text/javascript></SCRIPT><![endif]-->
									   <SCRIPT language=javascript src="../res/flot/jquery.js" type=text/javascript></SCRIPT>
									   <SCRIPT language=javascript src="../res/flot/jquery.flot.js" type=text/javascript></SCRIPT>
									   <SCRIPT language=javascript src="../res/flot/jquery.flot.selection.js" type=text/javascript></SCRIPT>
									   <BODY bgcolor='#F4FFE4'>

									   <H1>Профили мощности за период <br>с  00:00 01.02.2020 по  24:00 29.02.2020</H1>
									   <DIV id='placeholder' style="width:640px;height:300px;"></DIV>
									   <DIV id='overview'></DIV>
									   <div id="choices"></div><br>

									   <H2>М234 (Сетевой адрес - 17, Серийный номер - 23456789)</H2>
									   <H2></H2>
									   <table id="tableGraph2" style="width:1100;border-collapse:collapse;margin:1em 0">
									   <THEAD>
					  <TR>
					  <TH class=style20>№</TH>
					  <TH class=style20>P+, кВт</TH>
					  <TH class=style20>P-, кВт</TH>
					  <TH class=style20>Q+, квар</TH>
					  <TH class=style20>Q-, квар</TH>
					  <TH class=style20>Время</TH>
					  <TH class=style20>Дата</TH>
					  <TH class=style20>Период, мин.</TH>
					  <TH class=style20>Примечание</TH>
					  <TH class=style20>UTC(мс)</TH></TR></THEAD>
					  <TBODY>
					  <TR>
					  <TD class=style21>2</TD>
					  <TD class=style21>0.0705</TD>
					  <TD class=style21>0.0000</TD>
					  <TD class=style21>0.0063</TD>
					  <TD class=style21>0.0019</TD>
					  <TD class=style21>01:00</TD>
					  <TD class=style21>01.02.20</TD>
					  <TD class=style21>30+30</TD>
					  <TD class=style21>-</TD>
					  <TD class=style21>1580518800000</TD></TR>
					  <TR>`),
			},
			want: []*Row{
				{
					ID:     2,
					PPlus:  0.0705,
					PMinus: 0,
					QPlus:  0.0063,
					QMinus: 0.0019,
					Period: "30+30",
					Note:   "-",
					Date:   time.Date(2020, 2, 1, 1, 0, 0, 0, time.UTC),
				},
			},
			want1:   "23456789",
			wantErr: false,
		},
		{
			name: "fail date",
			args: args{
				data: []byte(`
					  <TD class=style21>2</TD>
					  <TD class=style21>0.0705</TD>
					  <TD class=style21>0.0000</TD>
					  <TD class=style21>0.0063</TD>
					  <TD class=style21>0.0019</TD>
					  <TD class=style21>01:00</TD>
					  <TD class=style21>01.02.20</TD>
					  <TD class=style21>30+30</TD>
					  <TD class=style21>-</TD>
					  <TD class=style21>fail</TD></TR>`),
			},
			want:    nil,
			want1:   "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := parse(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parse() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("parse() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
