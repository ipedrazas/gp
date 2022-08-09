/*
Copyright © 2022 Ivan Pedrazas <ipedrazas@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package files

import (
	"testing"
)

func TestLoad(t *testing.T) {

	utest := &UnitTest{}
	type args struct {
		configFile string
		out        interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "t01", args: args{configFile: "./unittest.yaml", out: utest}, wantErr: false},
		{name: "t01", args: args{configFile: "./boom/unittest.yaml", out: utest}, wantErr: true},
		{name: "t01", args: args{configFile: "./unittest-error.yaml", out: &Person{}}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Load(tt.args.configFile, tt.args.out); (err != nil) != tt.wantErr {
				t.Errorf("Load() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

type Person struct {
	Name     string
	Parent   *Person
	Children complex64
}

func TestSaveAsYaml(t *testing.T) {
	utest := &UnitTest{
		Name:      "cue",
		TestCases: []string{"TestParseCUEFile", "TestParseCUEFile"},
	}

	type args struct {
		path string
		out  interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "t01", args: args{path: "./unittest.yaml", out: utest}, wantErr: false},
		{name: "t02", args: args{path: "./boom/boom.yaml", out: utest}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := SaveAsYaml(tt.args.path, tt.args.out); (err != nil) != tt.wantErr {
				t.Errorf("SaveAsYaml() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCopy(t *testing.T) {
	type args struct {
		source string
		target string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "t01", args: args{source: "./unittest.yaml", target: "./unittest.yaml"}, wantErr: false},
		{name: "t01", args: args{source: "./boom/unittest.yaml", target: "./unittest.yaml"}, wantErr: true},
		{name: "t01", args: args{source: "./unittest.yaml", target: "./boom/unittest.yaml"}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Copy(tt.args.source, tt.args.target); (err != nil) != tt.wantErr {
				t.Errorf("Copy() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
