package extloadtest

import (
	"reflect"
	"testing"
)

func Test_getRandomIndices(t *testing.T) {
	type args struct {
		count2Change int
		maxLen       int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "getRandomIndices",
			args: args{
				count2Change: 2,
				maxLen:       10,
			},
			want: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getRandomIndices(tt.args.count2Change, tt.args.maxLen); !reflect.DeepEqual(len(got), tt.want) {
				t.Errorf("getRandomIndices() = %v, want %v", got, tt.want)
			}
		})
	}
}
