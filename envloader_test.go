package envloader

import (
	"reflect"
	"testing"
)

type TestStruct struct {
	FirstName string `env:"first_name"`
	LastName  string `env:"last_name,optional,default=Bavarian"`
}

func TestLoad(t *testing.T) {
	var success TestStruct

	type args struct {
		vars      interface{}
		filenames []string
	}
	tests := []struct {
		name    string
		args    args
		want    TestStruct
		wantErr bool
	}{
		{
			name: "success full filled",
			args: args{vars: &success, filenames: []string{"test.env"}},
			want: TestStruct{FirstName: "Bava", LastName: "Bavarian"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Load(tt.args.vars, tt.args.filenames...); (err != nil) != tt.wantErr {
				t.Errorf("Load() error = %v, wantErr %v", err, tt.wantErr)
			}

			if reflect.DeepEqual(tt.want, tt.args.vars) {
				t.Errorf("Load() want = %v, got %v", tt.want, tt.args.vars)
			}
		})
	}
}
