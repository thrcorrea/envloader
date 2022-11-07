package envloader

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
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

func TestLoadSlice(t *testing.T) {
	type StructTest struct {
		SliceInt    []int    `env:"SLICE_INT"`
		SliceInt16  []int16  `env:"SLICE_INT"`
		SliceInt32  []int32  `env:"SLICE_INT"`
		SliceInt64  []int64  `env:"SLICE_INT_64"`
		SliceString []string `env:"SLICE_STRING"`
	}

	type args struct {
		vars      StructTest
		filenames []string
	}
	tests := []struct {
		name    string
		args    args
		want    StructTest
		wantErr bool
	}{
		{
			name: "success full filled",
			args: args{vars: StructTest{}, filenames: []string{"test.env"}},
			want: StructTest{
				SliceInt:    []int{1, 2, 3, 4},
				SliceInt16:  []int16{1, 2, 3, 4},
				SliceInt32:  []int32{1, 2, 3, 4},
				SliceInt64:  []int64{1, 2, 3, 9223372036854775807},
				SliceString: []string{"a", "A", "ç", "ã", "1"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Load(&tt.args.vars, tt.args.filenames...); (err != nil) != tt.wantErr {
				t.Errorf("Load() error = %v, wantErr %v", err, tt.wantErr)
			}

			assert.Equal(t, tt.want, tt.args.vars)
		})
	}
}

func TestLoadDefault(t *testing.T) {
	type StructTest struct {
		DefaultString string `env:"DEFAULT_STRING,default=Bavarian"`
		DefaultInt    int    `env:"DEFAULT_INT,default=1"`
		DefaultInt16  int16  `env:"DEFAULT_INT,default=2"`
		DefaultInt32  int32  `env:"DEFAULT_INT,default=3"`
		DefaultInt64  int64  `env:"DEFAULT_INT,default=4"`
	}

	type args struct {
		vars      StructTest
		filenames []string
	}
	tests := []struct {
		name    string
		args    args
		want    StructTest
		wantErr bool
	}{
		{
			name: "success full filled",
			args: args{vars: StructTest{}, filenames: []string{"test.env"}},
			want: StructTest{
				DefaultString: "Bavarian",
				DefaultInt:    1,
				DefaultInt16:  2,
				DefaultInt32:  3,
				DefaultInt64:  4,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Load(&tt.args.vars, tt.args.filenames...); (err != nil) != tt.wantErr {
				t.Errorf("Load() error = %v, wantErr %v", err, tt.wantErr)
			}

			assert.Equal(t, tt.want, tt.args.vars)
		})
	}
}
