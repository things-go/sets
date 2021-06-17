/*
Copyright 2014 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package sets

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStringSet(t *testing.T) {
	s := String{}
	s2 := String{}
	if len(s) != 0 {
		t.Errorf("Expected len=0: %d", len(s))
	}
	s.Insert("a", "b")
	if len(s) != 2 {
		t.Errorf("Expected len=2: %d", len(s))
	}
	s.Insert("c")
	if s.Contains("d") {
		t.Errorf("Unexpected contents: %#v", s)
	}
	if !s.Contains("a") {
		t.Errorf("Missing contents: %#v", s)
	}
	s.Delete("a")
	if s.Contains("a") {
		t.Errorf("Unexpected contents: %#v", s)
	}
	s.Insert("a")
	if s.ContainsAll("a", "b", "d") {
		t.Errorf("Unexpected contents: %#v", s)
	}
	if !s.ContainsAll("a", "b") {
		t.Errorf("Missing contents: %#v", s)
	}
	s2.Insert("a", "b", "d")
	if s.IsSuperset(s2) {
		t.Errorf("Unexpected contents: %#v", s)
	}
	if s2.IsSubset(s) {
		t.Errorf("Unexpected contents: %#v", s)
	}
	s2.Delete("d")
	if !s.IsSuperset(s2) {
		t.Errorf("Missing contents: %#v", s)
	}
	if !s2.IsSubset(s) {
		t.Errorf("Unexpected contents: %#v", s)
	}
	_, ok := s2.Pop()
	if !ok {
		t.Errorf("Unexpected status: %#v", ok)
	}
	s2 = NewString()
	if s2.Len() != 0 {
		t.Errorf("Expected len=0: %d", len(s2))
	}
	v, ok := s2.Pop()
	if ok {
		t.Errorf("Unexpected status: %#v", ok)
	}
	if v != "" {
		t.Errorf("Unexpected value: %#v", v)
	}
	// improve cover
	s2 = NewStringFrom(map[string]interface{}{"1": "1", "2": "2", "3": "3"})
	s2.UnsortedList()
}

func TestStringSetDeleteMultiples(t *testing.T) {
	s := String{}
	s.Insert("a", "b", "c")
	if len(s) != 3 {
		t.Errorf("Expected len=3: %d", len(s))
	}

	s.Delete("a", "c")
	if len(s) != 1 {
		t.Errorf("Expected len=1: %d", len(s))
	}
	if s.Contains("a") {
		t.Errorf("Unexpected contents: %#v", s)
	}
	if s.Contains("c") {
		t.Errorf("Unexpected contents: %#v", s)
	}
	if !s.Contains("b") {
		t.Errorf("Missing contents: %#v", s)
	}
}

func TestNewStringSet(t *testing.T) {
	s := NewString("a", "b", "c")
	if len(s) != 3 {
		t.Errorf("Expected len=3: %d", len(s))
	}
	if !s.Contains("a") || !s.Contains("b") || !s.Contains("c") {
		t.Errorf("Unexpected contents: %#v", s)
	}
}

func TestStringSetList(t *testing.T) {
	s := NewString("z", "y", "x", "a")
	if !reflect.DeepEqual(s.List(), []string{"a", "x", "y", "z"}) {
		t.Errorf("List gave unexpected result: %#v", s.List())
	}
}

func TestStringSetDifference(t *testing.T) {
	a := NewString("1", "2", "3")
	b := NewString("1", "2", "4", "5")
	c := a.Difference(b)
	d := b.Difference(a)
	if len(c) != 1 {
		t.Errorf("Expected len=1: %d", len(c))
	}
	if !c.Contains("3") {
		t.Errorf("Unexpected contents: %#v", c.List())
	}
	if len(d) != 2 {
		t.Errorf("Expected len=2: %d", len(d))
	}
	if !d.Contains("4") || !d.Contains("5") {
		t.Errorf("Unexpected contents: %#v", d.List())
	}
}

func TestStringSetHasAny(t *testing.T) {
	a := NewString("1", "2", "3")

	if !a.ContainsAny("1", "4") {
		t.Errorf("expected true, got false")
	}

	if a.ContainsAny("0", "4") {
		t.Errorf("expected false, got true")
	}
}

func TestStringSetEquals(t *testing.T) {
	// Simple case (order doesn't matter)
	a := NewString("1", "2")
	b := NewString("2", "1")
	if !a.Equal(b) {
		t.Errorf("Expected to be equal: %v vs %v", a, b)
	}

	// It is a set; duplicates are ignored
	b = NewString("2", "2", "1")
	if !a.Equal(b) {
		t.Errorf("Expected to be equal: %v vs %v", a, b)
	}

	// Edge cases around empty sets / empty strings
	a = NewString()
	b = NewString()
	if !a.Equal(b) {
		t.Errorf("Expected to be equal: %v vs %v", a, b)
	}

	b = NewString("1", "2", "3")
	if a.Equal(b) {
		t.Errorf("Expected to be not-equal: %v vs %v", a, b)
	}

	b = NewString("1", "2", "")
	if a.Equal(b) {
		t.Errorf("Expected to be not-equal: %v vs %v", a, b)
	}

	// Check for equality after mutation
	a = NewString()
	a.Insert("1")
	if a.Equal(b) {
		t.Errorf("Expected to be not-equal: %v vs %v", a, b)
	}

	a.Insert("2")
	if a.Equal(b) {
		t.Errorf("Expected to be not-equal: %v vs %v", a, b)
	}

	a.Insert("")
	if !a.Equal(b) {
		t.Errorf("Expected to be equal: %v vs %v", a, b)
	}

	a.Delete("")
	if a.Equal(b) {
		t.Errorf("Expected to be not-equal: %v vs %v", a, b)
	}
}

func TestStringUnion(t *testing.T) {
	tests := []struct {
		s1       String
		s2       String
		expected String
	}{
		{
			NewString("1", "2", "3", "4"),
			NewString("3", "4", "5", "6"),
			NewString("1", "2", "3", "4", "5", "6"),
		},
		{
			NewString("1", "2", "3", "4"),
			NewString(),
			NewString("1", "2", "3", "4"),
		},
		{
			NewString(),
			NewString("1", "2", "3", "4"),
			NewString("1", "2", "3", "4"),
		},
		{
			NewString(),
			NewString(),
			NewString(),
		},
	}

	for _, test := range tests {
		union := test.s1.Union(test.s2)
		if union.Len() != test.expected.Len() {
			t.Errorf("Expected union.Len()=%d but got %d", test.expected.Len(), union.Len())
		}

		if !union.Equal(test.expected) {
			t.Errorf("Expected union.Equal(expected) but not true.  union:%v expected:%v", union.List(), test.expected.List())
		}
	}
}

func TestStringIntersection(t *testing.T) {
	tests := []struct {
		s1       String
		s2       String
		expected String
	}{
		{
			NewString("1", "2", "3", "4"),
			NewString("3", "4", "5", "6"),
			NewString("3", "4"),
		},
		{
			NewString("1", "2", "3", "4"),
			NewString("1", "2", "3", "4"),
			NewString("1", "2", "3", "4"),
		},
		{
			NewString("1", "2", "3", "4"),
			NewString(),
			NewString(),
		},
		{
			NewString(),
			NewString("1", "2", "3", "4"),
			NewString(),
		},
		{
			NewString(),
			NewString(),
			NewString(),
		},
	}

	for _, test := range tests {
		intersection := test.s1.Intersection(test.s2)
		if intersection.Len() != test.expected.Len() {
			t.Errorf("Expected intersection.Len()=%d but got %d", test.expected.Len(), intersection.Len())
		}

		if !intersection.Equal(test.expected) {
			t.Errorf("Expected intersection.Equal(expected) but not true.  intersection:%v expected:%v",
				intersection.List(), test.expected.List())
		}
	}
}

func TestStringMerge(t *testing.T) {
	tests := []struct {
		s1       String
		s2       String
		expected String
	}{
		{
			NewString("1", "2", "3", "4"),
			NewString("3", "4", "5", "6"),
			NewString("1", "2", "3", "4", "5", "6"),
		},
		{
			NewString("1", "2", "3", "4"),
			NewString("1", "2", "3", "4"),
			NewString("1", "2", "3", "4"),
		},
		{
			NewString("1", "2", "3", "4"),
			NewString(),
			NewString("1", "2", "3", "4"),
		},
		{
			NewString(),
			NewString("1", "2", "3", "4"),
			NewString("1", "2", "3", "4"),
		},
		{
			NewString(),
			NewString(),
			NewString(),
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

func TestString_Each(t *testing.T) {
	expect := NewString("1", "2", "3", "4")
	s1 := NewString("1", "2", "3", "4")
	s1.Each(func(item string) bool {
		require.True(t, expect.Contains(item))
		return item != "3"
	})
}

func TestString_Clone(t *testing.T) {
	s1 := NewString("1", "2", "3", "4")
	s2 := s1.Clone()

	require.True(t, s1.Equal(s2))
}
