package logic

import (
	"reflect"
	"testing"
)

func TestDifference(t *testing.T) {
	tests := []struct {
		name string
		a    []string
		b    []string
		want []string
	}{
		{
			name: "basic difference",
			a:    []string{"a", "b", "c"},
			b:    []string{"b"},
			want: []string{"a", "c"},
		},
		{
			name: "no difference",
			a:    []string{"a", "b"},
			b:    []string{"a", "b"},
			want: nil,
		},
		{
			name: "empty a",
			a:    []string{},
			b:    []string{"a"},
			want: nil,
		},
		{
			name: "empty b",
			a:    []string{"a", "b"},
			b:    []string{},
			want: []string{"a", "b"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Difference(tt.a, tt.b); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Difference() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIntersection(t *testing.T) {
	tests := []struct {
		name string
		a    []string
		b    []string
		want []string
	}{
		{
			name: "basic intersection",
			a:    []string{"a", "b", "c"},
			b:    []string{"b", "c", "d"},
			want: []string{"b", "c"},
		},
		{
			name: "no intersection",
			a:    []string{"a", "b"},
			b:    []string{"c", "d"},
			want: nil,
		},
		{
			name: "empty a",
			a:    []string{},
			b:    []string{"a"},
			want: nil,
		},
		{
			name: "empty b",
			a:    []string{"a"},
			b:    []string{},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Intersection(tt.a, tt.b); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Intersection() = %v, want %v", got, tt.want)
			}
		})
	}
}
