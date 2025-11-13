package inscriptos

import (
	"testing"

	"github.com/google/uuid"
)

func TestHelpers_ParseUUID(t *testing.T) {
	id := uuid.New()
	got := parseUUID(id.String())
	if got != id {
		t.Fatalf("expected %v, got %v", id, got)
	}
	if parseUUID(nil) != uuid.Nil {
		t.Fatalf("expected zero UUID for nil")
	}
}

func TestHelpers_ToBool(t *testing.T) {
	cases := []struct {
		in   interface{}
		want bool
	}{
		{true, true}, {false, false},
		{int64(1), true}, {int64(0), false},
		{int(1), true}, {int(0), false},
		{float64(1), true}, {float64(0), false},
		{[]byte("1"), true}, {[]byte("0"), false},
		{"true", true}, {"false", false},
		{"1", true}, {"0", false},
		{struct{}{}, false},
	}
	for i, c := range cases {
		if got := toBool(c.in); got != c.want {
			t.Fatalf("case %d: toBool(%v) = %v, want %v", i, c.in, got, c.want)
		}
	}
}

func TestHelpers_ToInt(t *testing.T) {
	cases := []struct {
		in   interface{}
		want int
	}{
		{int(5), 5}, {int32(6), 6}, {int64(7), 7},
		{float32(2.9), 2}, {float64(3.1), 3},
		{[]byte("42"), 42}, {struct{}{}, 0},
	}
	for i, c := range cases {
		if got := toInt(c.in); got != c.want {
			t.Fatalf("case %d: toInt(%v) = %v, want %v", i, c.in, got, c.want)
		}
	}
}

func TestHelpers_ToFloat64(t *testing.T) {
	cases := []struct {
		in   interface{}
		want float64
	}{
		{float64(1.5), 1.5}, {float32(1.25), 1.25},
		{int(2), 2.0}, {int64(3), 3.0},
		{[]byte("4.75"), 4.75}, {struct{}{}, 0.0},
	}
	for i, c := range cases {
		if got := toFloat64(c.in); (got-c.want) > 0.0001 || (c.want-got) > 0.0001 {
			t.Fatalf("case %d: toFloat64(%v) = %v, want %v", i, c.in, got, c.want)
		}
	}
}
