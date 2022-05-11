package models

import "testing"

func Test_parseCMD(t *testing.T) {
	type args struct {
		cmd string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "t01", args: args{cmd: "ls -latrh"}, want: "[\"ls\", \"-latrh\"]"},
		{name: "t02", args: args{cmd: "myapp"}, want: "[\"/myapp\"]"},
		{name: "t03", args: args{cmd: "/myapp"}, want: "[\"/myapp\"]"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseCMD(tt.args.cmd); got != tt.want {
				t.Errorf("parseCMD() = %v, want %v", got, tt.want)
			}
		})
	}
}
