package dynamic_gui_config

import (
	"reflect"
	"testing"

	_ "github.com/andlabs/ui/winmanifest"
)

func Test_structBreakdown(t *testing.T) {
	floatValue := float64(2)

	type args struct {
		structPtr interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    []structGuiField
		wantErr bool
	}{
		{
			name: "should return error when type is not a pointer",
			args: args{
				structPtr: floatValue,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "should return error when type is not a pointer",
			args: args{
				structPtr: struct{}{},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "should return no error when type is a pointer to struct",
			args: args{
				structPtr: &struct{}{},
			},
			want:    []structGuiField{},
			wantErr: false,
		},
		{
			name: "should return error when type is a pointer to float",
			args: args{
				structPtr: &floatValue,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := structBreakdown(tt.args.structPtr)
			if (err != nil) != tt.wantErr {
				t.Errorf("structBreakdown() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("structBreakdown() got = %v, want %v", got, tt.want)
			}
		})
	}
}
