package sets

import (
	"reflect"
	"slices"
	"sort"
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

func TestSetList(t *testing.T) {
	s := New(13, 12, 1, 11)
	v := s.List()

	for _, vv := range v {
		if !s.Contains(vv) {
			t.Errorf("List gave unexpected result: %#v", v)
		}
	}
}

func TestSetDifference(t *testing.T) {
	a := New(1, 2, 3)
	b := New(1, 2, 4, 5)

	t.Run("set", func(t *testing.T) {
		c := a.Difference(b)
		d := b.Difference(a)
		if len(c) != 1 {
			t.Errorf("Expected len=1: %d", len(c))
		}
		if !c.Contains(3) {
			t.Errorf("Unexpected contents: %#v", c.List())
		}
		if len(d) != 2 {
			t.Errorf("Expected len=2: %d", len(d))
		}
		if !d.Contains(4) || !d.Contains(5) {
			t.Errorf("Unexpected contents: %#v", d.List())
		}
	})

	t.Run("slice", func(t *testing.T) {
		c := a.DifferenceSlice(b)
		d := b.DifferenceSlice(a)
		if len(c) != 1 {
			t.Errorf("Expected len=1: %d", len(c))
		}
		if !slices.Contains(c, 3) {
			t.Errorf("Unexpected contents: %#v", c)
		}
		if len(d) != 2 {
			t.Errorf("Expected len=2: %d", len(d))
		}
		if !slices.Contains(d, 4) || !slices.Contains(d, 5) {
			t.Errorf("Unexpected contents: %#v", d)
		}
	})
}

func TestSetDiff(t *testing.T) {
	t.Run("set", func(t *testing.T) {
		a := New(1, 3, 5, 7)
		b := New(3, 4, 5, 6)

		added, removed, remained := a.Diff(b)
		wantAdded := Set[int]{
			4: {},
			6: {},
		}
		if !reflect.DeepEqual(added, wantAdded) {
			t.Errorf("Expected %#v, Got: %#v", wantAdded, added)
		}
		wantRemoved := Set[int]{
			1: {},
			7: {},
		}
		if !reflect.DeepEqual(removed, wantRemoved) {
			t.Errorf("Expected %v, Got: %v", wantRemoved, removed)
		}
		wantRemained := Set[int]{
			3: {},
			5: {},
		}
		if !reflect.DeepEqual(remained, wantRemained) {
			t.Errorf("Expected %v, Got: %v", wantRemained, remained)
		}
	})
	t.Run("slice", func(t *testing.T) {
		a := New(1, 3, 5, 7)
		b := New(3, 4, 5, 6)

		added, removed, remained := a.DiffSlice(b)
		sort.Ints(added)
		sort.Ints(removed)
		sort.Ints(remained)

		wantAdded := []int{4, 6}
		if !reflect.DeepEqual(added, wantAdded) {
			t.Errorf("Expected %#v, Got: %#v", wantAdded, added)
		}

		wantRemoved := []int{1, 7}
		if !reflect.DeepEqual(removed, wantRemoved) {
			t.Errorf("Expected %v, Got: %v", wantRemoved, removed)
		}

		wantRemained := []int{3, 5}
		if !reflect.DeepEqual(remained, wantRemained) {
			t.Errorf("Expected %v, Got: %v", wantRemained, remained)
		}
	})
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

func TestUnion(t *testing.T) {
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
		// set
		union := test.s1.Union(test.s2)
		if union.Len() != test.expected.Len() {
			t.Errorf("Expected union.Len()=%d but got %d", test.expected.Len(), union.Len())
		}
		if !union.Equal(test.expected) {
			t.Errorf("Expected union.Equal(expected) but not true.  union:%v expected:%v", union.List(), test.expected.List())
		}
		// slice
		unionSlices := test.s1.UnionSlice(test.s2)
		if len(unionSlices) != test.expected.Len() {
			t.Errorf("Expected UnionSlices length=%d but got %d", test.expected.Len(), len(unionSlices))
		}
		if !New(unionSlices...).Equal(test.expected) {
			t.Errorf("Expected union.Equal(expected) but not true.  union:%v expected:%v", unionSlices, test.expected.List())
		}
	}
}

func TestIntersection(t *testing.T) {
	tests := []struct {
		s1       Set[int]
		s2       Set[int]
		expected Set[int]
	}{
		{
			New[int](1, 2, 3, 4),
			New[int](3, 4, 5, 6),
			New[int](3, 4),
		},
		{
			New[int](1, 2, 3, 4),
			New[int](1, 2, 3, 4),
			New[int](1, 2, 3, 4),
		},
		{
			New[int](1, 2, 3, 4),
			New[int](),
			New[int](),
		},
		{
			New[int](),
			New[int](1, 2, 3, 4),
			New[int](),
		},
		{
			New[int](),
			New[int](),
			New[int](),
		},
	}

	for _, test := range tests {
		// sets
		intersection := test.s1.Intersection(test.s2)
		if intersection.Len() != test.expected.Len() {
			t.Errorf("Expected intersection.Len()=%d but got %d", test.expected.Len(), intersection.Len())
		}
		if !intersection.Equal(test.expected) {
			t.Errorf("Expected intersection.Equal(expected) but not true.  intersection:%v expected:%v",
				intersection.List(), test.expected.List())
		}

		// slice
		intersectionSlice := test.s1.IntersectionSlice(test.s2)
		if len(intersectionSlice) != test.expected.Len() {
			t.Errorf("Expected IntersectionSlice length =%d but got %d", test.expected.Len(), len(intersectionSlice))
		}
		if !New(intersectionSlice...).Equal(test.expected) {
			t.Errorf("Expected intersection.Equal(expected) but not true.  intersection:%v expected:%v",
				intersectionSlice, test.expected.List())
		}
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

func Test_Clone(t *testing.T) {
	s1 := New(1, 2, 3, 4)
	s2 := s1.Clone()
	if got := s1.Equal(s2); !got {
		t.Errorf("Expected Equal()=%v but got %v", true, got)
	}
}
