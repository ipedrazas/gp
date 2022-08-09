package models

import (
	"testing"
)

func TestCompose_GetServiceNames(t *testing.T) {
	type fields struct {
		Services map[string]Service
	}

	svcs := make(map[string]Service)
	svc := &Service{
		Image: "test",
	}

	svcs["uno"] = *svc
	svcs["dos"] = *svc
	svcs["tres"] = *svc

	tests := []struct {
		name   string
		fields fields
		want   []string
	}{
		{name: "t01", fields: fields{
			Services: svcs,
		}, want: []string{"uno", "dos", "tres"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Compose{
				Services: tt.fields.Services,
			}
			got := c.GetServiceNames()
			for _, entry := range got {
				if !contains(entry, tt.want) {
					t.Errorf("Compose.GetServiceNames() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}
