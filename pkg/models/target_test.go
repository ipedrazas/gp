package models

import (
	"os"
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

func TestTarget_Run(t *testing.T) {
	type fields struct {
		Name           string
		Platform       []string
		Domain         string
		Overwrite      bool
		Registry       string
		RegistryUserId string
		Image          string
		Compose        string
		Actions        []string
		DockerBuild    bool
		Paused         bool
	}
	type args struct {
		comp   *Component
		gitSha string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// test paused
		{name: "t01", fields: fields{Name: "test", Paused: true}, args: args{}, wantErr: false},
		// test no docker build no compose
		{name: "t02", fields: fields{Name: "test", DockerBuild: false}, args: args{}, wantErr: false},
		// test no docker build no compose
		{name: "t03", fields: fields{Name: "test", DockerBuild: false, Compose: "nofile.yaml"}, args: args{}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := &Target{
				Name:           tt.fields.Name,
				Platform:       tt.fields.Platform,
				Domain:         tt.fields.Domain,
				Overwrite:      tt.fields.Overwrite,
				Registry:       tt.fields.Registry,
				RegistryUserId: tt.fields.RegistryUserId,
				Image:          tt.fields.Image,
				Compose:        tt.fields.Compose,
				Actions:        tt.fields.Actions,
				DockerBuild:    tt.fields.DockerBuild,
				Paused:         tt.fields.Paused,
			}
			if err := tr.Run(tt.args.comp, tt.args.gitSha); (err != nil) != tt.wantErr {
				t.Errorf("Target.Run() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_setEnvVars(t *testing.T) {
	version := "123"
	git := "abc"
	type args struct {
		version string
		gitSha  string
	}
	tests := []struct {
		name string
		args args
	}{
		{name: "t01", args: args{version: version, gitSha: git}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setEnvVars(tt.args.version, tt.args.gitSha)
			vversion := os.Getenv("TAG")
			if version != vversion {
				t.Errorf("Target.setEnvVars() expected %v got %v", version, vversion)
			}
		})
	}
}
