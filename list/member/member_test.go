package member

import "testing"

func TestList(t *testing.T) {
	l := NewList()
	l.Add(`a`, `b`, nil)
	if l := len(l.Members()); l != 1 {
		t.Fatal(`want 1 got`, l)
	}
}
