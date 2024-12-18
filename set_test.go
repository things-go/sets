package sets

import (
	"testing"
)

func TestSet(t *testing.T) {
	s := Set[int]{}
	s2 := Set[int]{}
	if len(s) != 0 {
		t.Errorf("Expected len=0: %d", len(s))
	}
	s.Insert(1, 2)
	if len(s) != 2 {
		t.Errorf("Expected len=2: %d", len(s))
	}
	s.Insert(3)
	if s.Contains(4) {
		t.Errorf("Unexpected contents: %#v", s)
	}
	if !s.Contains(1) {
		t.Errorf("Missing contents: %#v", s)
	}
	s.Delete(1)
	if s.Contains(1) {
		t.Errorf("Unexpected contents: %#v", s)
	}
	s.Insert(1)
	if s.ContainsAll(1, 2, 4) {
		t.Errorf("Unexpected contents: %#v", s)
	}
	if !s.ContainsAll(1, 2) {
		t.Errorf("Missing contents: %#v", s)
	}
	s2.Insert(1, 2, 4)
	if s.IsSuperset(s2) {
		t.Errorf("Unexpected contents: %#v", s)
	}
	if s2.IsSubset(s) {
		t.Errorf("Unexpected contents: %#v", s)
	}
	s2.Delete(4)
	if !s.IsSuperset(s2) {
		t.Errorf("Missing contents: %#v", s)
	}
	if !s2.IsSubset(s) {
		t.Errorf("Missing contents: %#v", s)
	}
	_, ok := s2.Pop()
	if !ok {
		t.Errorf("Unexpected status: %#v", ok)
	}
	s2 = New[int]()
	if s2.Len() != 0 {
		t.Errorf("Expected len=0: %d", len(s2))
	}
	v, ok := s2.Pop()
	if ok {
		t.Errorf("Unexpected status: %#v", ok)
	}
	if v != 0 {
		t.Errorf("Unexpected value: %#v", v)
	}
}

func TestSetDeleteMultiples(t *testing.T) {
	s := NewFrom(map[int]any{1: "1", 2: "2", 3: "3"})
	if len(s) != 3 {
		t.Errorf("Expected len=3: %d", len(s))
	}

	s.Delete(1, 3)
	if len(s) != 1 {
		t.Errorf("Expected len=1: %d", len(s))
	}
	if s.Contains(1) {
		t.Errorf("Unexpected contents: %#v", s)
	}
	if s.Contains(3) {
		t.Errorf("Unexpected contents: %#v", s)
	}
	if !s.Contains(2) {
		t.Errorf("Missing contents: %#v", s)
	}
}

func TestNewSet(t *testing.T) {
	s := New(1, 2, 3)
	if len(s) != 3 {
		t.Errorf("Expected len=3: %d", len(s))
	}
	if !s.Contains(1) || !s.Contains(2) || !s.Contains(3) {
		t.Errorf("Unexpected contents: %#v", s)
	}
}

func TestSetHasAny(t *testing.T) {
	a := New(1, 2, 3)

	if !a.ContainsAny(1, 4) {
		t.Errorf("expected true, got false")
	}

	if a.ContainsAny(10, 4) {
		t.Errorf("expected false, got true")
	}
}

func TestSetEquals(t *testing.T) {
	// Simple case (order doesn't matter)
	a := New[int](1, 2)
	b := New[int](2, 1)
	if !a.Equal(b) {
		t.Errorf("Expected to be equal: %v vs %v", a, b)
	}

	// It is a set; duplicates are ignored
	b = New[int](2, 2, 1)
	if !a.Equal(b) {
		t.Errorf("Expected to be equal: %v vs %v", a, b)
	}

	// Edge cases around empty sets / empty strings
	a = New[int]()
	b = New[int]()
	if !a.Equal(b) {
		t.Errorf("Expected to be equal: %v vs %v", a, b)
	}

	b = New(1, 2, 3)
	if a.Equal(b) {
		t.Errorf("Expected to be not-equal: %v vs %v", a, b)
	}

	b = New(1, 2, 0)
	if a.Equal(b) {
		t.Errorf("Expected to be not-equal: %v vs %v", a, b)
	}

	// Check for equality after mutation
	a = New[int]()
	a.Insert(1)
	if a.Equal(b) {
		t.Errorf("Expected to be not-equal: %v vs %v", a, b)
	}

	a.Insert(2)
	if a.Equal(b) {
		t.Errorf("Expected to be not-equal: %v vs %v", a, b)
	}

	a.Insert(0)
	if !a.Equal(b) {
		t.Errorf("Expected to be equal: %v vs %v", a, b)
	}

	a.Delete(0)
	if a.Equal(b) {
		t.Errorf("Expected to be not-equal: %v vs %v", a, b)
	}
}

func TestMerge(t *testing.T) {
	tests := []struct {
		s1       Set[int]
		s2       Set[int]
		expected Set[int]
	}{
		{
			New[int](1, 2, 3, 4),
			New[int](3, 4, 5, 6),
			New[int](1, 2, 3, 4, 5, 6),
		},
		{
			New[int](1, 2, 3, 4),
			New[int](1, 2, 3, 4),
			New[int](1, 2, 3, 4),
		},
		{
			New[int](1, 2, 3, 4),
			New[int](),
			New[int](1, 2, 3, 4),
		},
		{
			New[int](),
			New[int](1, 2, 3, 4),
			New[int](1, 2, 3, 4),
		},
		{
			New[int](),
			New[int](),
			New[int](),
		},
	}

	for _, test := range tests {
		intersection := test.s1.Merge(test.s2)
		if intersection.Len() != test.expected.Len() {
			t.Errorf("Expected merge.Len()=%d but got %d", test.expected.Len(), intersection.Len())
		}

		if !intersection.Equal(test.expected) {
			t.Errorf("Expected merge.Equal(expected) but not true.  merge:%v expected:%v",
				intersection.List(), test.expected.List())
		}
	}
}

func Test_Clone(t *testing.T) {
	s1 := New(1, 2, 3, 4)
	s2 := s1.Clone()
	if got := s1.Equal(s2); !got {
		t.Errorf("Expected Equal()=%v but got %v", true, got)
	}
}

func TestSetList(t *testing.T) {
	s := New(13, 12, 1, 11)
	v := s.List()

	for _, vv := range v {
		if !s.Contains(vv) {
			t.Errorf("List gave unexpected result: %#v", v)
		}
	}
}
func Test_Each(t *testing.T) {
	expect := New(1, 2, 3, 4)
	s1 := New(1, 2, 3, 4)
	s1.Each(func(item int) bool {
		if got := expect.Contains(item); !got {
			t.Errorf("Expected Equal()=%v but got %v", true, got)
		}
		return item != 3
	})
}
