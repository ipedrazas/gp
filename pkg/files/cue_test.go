package files

import (
	"testing"
)

type UnitTest struct {
	Name      string   `json:"name,omitempty" cue:"name,optional"`
	TestCases []string `json:"test_cases,omitempty" cue:"test_cases,optional"`
}

type UnitCaseSuite struct {
	Name  string     `json:"name,omitempty" cue:"name,optional"`
	Tests []UnitTest `json:"tests,omitempty" cue:"tests,optional"`
}

func TestParseCUEFile(t *testing.T) {
	utest := &UnitTest{}
	ucase := &UnitCaseSuite{}
	type args struct {
		source string
		obj    interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "t01", args: args{source: "./unittest.cue", obj: utest}},
		{name: "t02", args: args{source: "./unitcasesuite.cue", obj: ucase}},
		{name: "t03", args: args{source: "./error1.cue", obj: ucase}, wantErr: true},
		{name: "t04", args: args{source: "./unittest.yaml", obj: ucase}, wantErr: true},
		{name: "t05", args: args{source: "./unittest-error.cue", obj: nil}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ParseCUEFile(tt.args.source, tt.args.obj); (err != nil) != tt.wantErr {
				t.Errorf("ParseCUEFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
