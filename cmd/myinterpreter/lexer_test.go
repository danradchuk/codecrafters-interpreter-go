package main

import (
	"reflect"
	"testing"
)

func Test_lexify(t *testing.T) {
	type args struct {
		input []byte
	}
	tests := []struct {
		name  string
		args  args
		want  []Token
		want1 []error
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := lexify(tt.args.input)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("lexify() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("lexify() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
