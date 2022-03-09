package aigo

import "testing"

func TestPointString(t *testing.T) {
	p1 := Point{Row: 1, Col: 2}
	if p1.String() != "B1" {
		t.Fatalf("期望,实际%s", p1.String())
	}
}
