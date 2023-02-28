package remote

import "testing"

func TestFetchOCI(t *testing.T) {
	image := "docker.io/ipedrazas/targets:helm"
	digest := "sha256:7ff97545abf105da70429b8130ed0e22685bb15cc0f071a8675734d50f30e5d9"
	type args struct {
		ociArtifactUri string
		output         string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{name: "t01", args: args{ociArtifactUri: image, output: "tmp/t2"}, want: digest, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FetchOCI(tt.args.ociArtifactUri, tt.args.output)
			if (err != nil) != tt.wantErr {
				t.Errorf("FetchOCI() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("FetchOCI() = %v, want %v", got, tt.want)
			}
		})
	}
}
