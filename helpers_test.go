package main

import "testing"

func TestPortLabel(t *testing.T) {
	cases := []struct {
		i    int
		want string
	}{
		{0, "1"},
		{1, "2"},
		{9, "10"},
	}
	for _, c := range cases {
		if got := portLabel(c.i); got != c.want {
			t.Errorf("portLabel(%d) = %q, want %q", c.i, got, c.want)
		}
	}
}

func TestBitGet64(t *testing.T) {
	if v, ok := bitGet64(nil, 0); ok || v {
		t.Errorf("bitGet64(nil, 0) = (%v, %v), want (false, false)", v, ok)
	}

	mask := int64(0b0000_0101) // bits 0 and 2 set
	if v, ok := bitGet64(&mask, 0); !ok || !v {
		t.Errorf("bitGet64(mask, 0) = (%v, %v), want (true, true)", v, ok)
	}
	if v, ok := bitGet64(&mask, 1); !ok || v {
		t.Errorf("bitGet64(mask, 1) = (%v, %v), want (false, true)", v, ok)
	}
	if v, ok := bitGet64(&mask, 2); !ok || !v {
		t.Errorf("bitGet64(mask, 2) = (%v, %v), want (true, true)", v, ok)
	}
}

func TestBitGet32(t *testing.T) {
	if v, ok := bitGet32(nil, 0); ok || v {
		t.Errorf("bitGet32(nil, 0) = (%v, %v), want (false, false)", v, ok)
	}

	mask := int32(0b10)
	if v, ok := bitGet32(&mask, 1); !ok || !v {
		t.Errorf("bitGet32(mask, 1) = (%v, %v), want (true, true)", v, ok)
	}
	if v, ok := bitGet32(&mask, 0); !ok || v {
		t.Errorf("bitGet32(mask, 0) = (%v, %v), want (false, true)", v, ok)
	}
}

func TestArrayLenInt32(t *testing.T) {
	a := []int32{1, 2, 3}
	b := []int32{1, 2}
	if n := arrayLenInt32(nil); n != 0 {
		t.Errorf("arrayLenInt32(nil) = %d, want 0", n)
	}
	if n := arrayLenInt32(&a, &b); n != 3 {
		t.Errorf("arrayLenInt32(&a, &b) = %d, want 3", n)
	}
	if n := arrayLenInt32(&b, nil); n != 2 {
		t.Errorf("arrayLenInt32(&b, nil) = %d, want 2", n)
	}
}

func TestArrayLenString(t *testing.T) {
	a := []string{"x", "y"}
	if n := arrayLenString(nil); n != 0 {
		t.Errorf("arrayLenString(nil) = %d, want 0", n)
	}
	if n := arrayLenString(&a); n != 2 {
		t.Errorf("arrayLenString(&a) = %d, want 2", n)
	}
}

func TestAtInt32(t *testing.T) {
	a := []int32{10, 20, 30}
	if v := atInt32(nil, 0); v != nil {
		t.Errorf("atInt32(nil, 0) = %v, want nil", v)
	}
	if v := atInt32(&a, -1); v != nil {
		t.Errorf("atInt32(&a, -1) = %v, want nil", v)
	}
	if v := atInt32(&a, 3); v != nil {
		t.Errorf("atInt32(&a, 3) = %v, want nil (out of range)", v)
	}
	if v := atInt32(&a, 1); v == nil || *v != 20 {
		t.Errorf("atInt32(&a, 1) = %v, want 20", v)
	}
}

func TestAtInt64(t *testing.T) {
	a := []int64{100, 200}
	if v := atInt64(&a, 1); v == nil || *v != 200 {
		t.Errorf("atInt64(&a, 1) = %v, want 200", v)
	}
	if v := atInt64(&a, 5); v != nil {
		t.Errorf("atInt64(&a, 5) = %v, want nil", v)
	}
}

func TestAtString(t *testing.T) {
	a := []string{"port1", "port2"}
	if v := atString(&a, 0); v == nil || *v != "port1" {
		t.Errorf("atString(&a, 0) = %v, want \"port1\"", v)
	}
	if v := atString(&a, 2); v != nil {
		t.Errorf("atString(&a, 2) = %v, want nil", v)
	}
}

func TestDerefString(t *testing.T) {
	if got := derefString(nil); got != "" {
		t.Errorf("derefString(nil) = %q, want \"\"", got)
	}
	s := "hello"
	if got := derefString(&s); got != "hello" {
		t.Errorf("derefString(&s) = %q, want \"hello\"", got)
	}
}
