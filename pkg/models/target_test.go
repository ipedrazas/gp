package models

import (
	"testing"
)

func Test_getArch(t *testing.T) {
	type args struct {
		platform string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "t01", args: args{platform: "linux/arm64"}, want: "arm64"},
		{name: "t02", args: args{platform: "linux"}, want: ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getArch(tt.args.platform); got != tt.want {
				t.Errorf("getArch() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTarget_GetDockerImage(t *testing.T) {
	type fields struct {
		Name           string
		Registry       string
		RegistryUserId string
	}
	type args struct {
		appName string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{name: "t01",
			fields: fields{
				Name:           "test",
				Registry:       "harbor.alacasa.uk",
				RegistryUserId: "ivan",
			},
			args: args{appName: "unittest"},
			want: "harbor.alacasa.uk/ivan/unittest:tag"},
		{name: "t02",
			fields: fields{
				Name:           "test",
				Registry:       "",
				RegistryUserId: "ivan",
			},
			args: args{appName: "unittest"},
			want: "ivan/unittest:tag"},
		{name: "t03",
			fields: fields{
				Name:           "test",
				Registry:       "docker.io",
				RegistryUserId: "ivan",
			},
			args: args{appName: "unittest"},
			want: "ivan/unittest:tag"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			target := &Target{
				Name:           tt.fields.Name,
				Registry:       tt.fields.Registry,
				RegistryUserId: tt.fields.RegistryUserId,
			}
			target.SetDockerImage(tt.args.appName, "tag")
			if target.Image != tt.want {
				t.Errorf("target.Image == %v, want %v", target.Image, tt.want)
			}
		})
	}
}

func TestTarget_IsAvailable(t *testing.T) {
	type fields struct {
		Name string
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name:   "t01",
			fields: fields{Name: "available"},
			want:   true,
		},
		// {
		// 	name:   "t02",
		// 	fields: fields{Name: "local"},
		// 	want:   false,
		// },
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := &Target{
				Name: tt.fields.Name,
			}
			if got := tr.IsAvailable(); got != tt.want {
				t.Errorf("Target.IsAvailable() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_validateActions(t *testing.T) {
	type args struct {
		actions  []string
		services []string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{name: "t01",
			args: args{
				actions:  []string{"build", "hydrate"},
				services: []string{"build", "hydrate"},
			},
			want: true},
		{name: "t02",
			args: args{
				actions:  []string{"hydrate"},
				services: []string{"build", "hydrate"},
			},
			want: true},
		{name: "t03",
			args: args{
				actions:  []string{"hydrate"},
				services: []string{"build"},
			},
			want: false},
		{name: "t04",
			args: args{
				actions:  []string{},
				services: []string{"build"},
			},
			want: true},
		{name: "t05",
			args: args{
				actions:  []string{"build", "hydrate"},
				services: []string{},
			},
			want: false},
		{name: "t06",
			args: args{
				actions:  []string{"hydrate", "build"},
				services: []string{"build", "hydrate"},
			},
			want: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := validateActions(tt.args.actions, tt.args.services); got != tt.want {
				t.Errorf("validateActions() = %v, want %v", got, tt.want)
			}
		})
	}
}
