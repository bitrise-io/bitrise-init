package analytics

import (
	"reflect"
	"testing"
)

func Test_initData(t *testing.T) {
	tests := []struct {
		name       string
		data, want map[string]interface{}
	}{
		{
			name: "Empty data",
			data: map[string]interface{}{},
			want: map[string]interface{}{
				"source": "scanner",
			},
		},
		{
			name: "nil data",
			data: nil,
			want: map[string]interface{}{
				"source": "scanner",
			},
		},
		{
			name: "source is overwritten",
			data: map[string]interface{}{
				"source": "A",
			},
			want: map[string]interface{}{
				"source": "scanner",
			},
		},
		{
			name: "Existing data",
			data: map[string]interface{}{
				"A": "B",
			},
			want: map[string]interface{}{
				"source": "scanner",
				"A":      "B",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := initData(tt.data); !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("initData() = %v, want %s", got, tt.want)
			}
		})
	}
}
