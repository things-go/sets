/*
Copyright The Kubernetes Authors.

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
	"sort"
)

// Int is a set of ints, implemented via map[int]struct{} for minimal memory consumption.
type Int map[int]struct{}

// NewInt creates a Int from a list of values.
func NewInt(items ...int) Int {
	ss := Int{}
	ss.Insert(items...)
	return ss
}

// NewIntFrom creates a Int from a keys of a map[int](? extends interface{}).
// If the value passed in is not actually a map, this will panic.
func NewIntFrom(theMap interface{}) Int {
	v := reflect.ValueOf(theMap)
	ret := Int{}

	for _, keyValue := range v.MapKeys() {
		ret[keyValue.Interface().(int)] = struct{}{}
	}
	return ret
}

// Insert adds items to the set.
func (s Int) Insert(items ...int) Int {
	for _, item := range items {
		s[item] = struct{}{}
	}
	return s
}

// Delete removes all items from the set.
func (s Int) Delete(items ...int) Int {
	for _, item := range items {
		delete(s, item)
	}
	return s
}

// Contains returns true if and only if item is contained in the set.
func (s Int) Contains(item int) bool {
	_, contained := s[item]
	return contained
}

// ContainsAll returns true if and only if all items are contained in the set.
func (s Int) ContainsAll(items ...int) bool {
	for _, item := range items {
		if !s.Contains(item) {
			return false
		}
	}
	return true
}

// ContainsAny returns true if any items are contained in the set.
func (s Int) ContainsAny(items ...int) bool {
	for _, item := range items {
		if s.Contains(item) {
			return true
		}
	}
	return false
}

// Difference returns a set of objects that are not in s2
// For example:
// s1 = {a1, a2, a3}
// s2 = {a1, a2, a4, a5}
// s1.Difference(s2) = {a3}
// s2.Difference(s1) = {a4, a5}
func (s Int) Difference(s2 Int) Int {
	result := NewInt()
	for key := range s {
		if !s2.Contains(key) {
			result[key] = struct{}{}
		}
	}
	return result
}

// Union returns a new set which includes items in either s1 or s2.
// For example:
// s1 = {a1, a2}
// s2 = {a3, a4}
// s1.Union(s2) = {a1, a2, a3, a4}
// s2.Union(s1) = {a1, a2, a3, a4}
func (s Int) Union(s2 Int) Int {
	result := NewInt()
	for key := range s {
		result[key] = struct{}{}
	}
	for key := range s2 {
		result[key] = struct{}{}
	}
	return result
}

// Intersection returns a new set which includes the item in BOTH s1 and s2
// For example:
// s1 = {a1, a2}
// s2 = {a2, a3}
// s1.Intersection(s2) = {a2}
func (s Int) Intersection(s2 Int) Int {
	var walk, other Int
	result := NewInt()
	if s.Len() < s2.Len() {
		walk = s
		other = s2
	} else {
		walk = s2
		other = s
	}
	for key := range walk {
		if other.Contains(key) {
			result[key] = struct{}{}
		}
	}
	return result
}

// Merge adds item from s2 set into s1
// For example:
// s1 = {a1, a2}
// s2 = {a3, a4}
// s1.Merge(s2) = {a1, a2, a3, a4}
func (s Int) Merge(s2 Int) Int {
	for key := range s2 {
		s[key] = struct{}{}
	}
	return s
}

// IsSuperset returns true if and only if s1 is a superset of s2.
func (s Int) IsSuperset(s2 Int) bool {
	for item := range s2 {
		if !s.Contains(item) {
			return false
		}
	}
	return true
}

// IsSubset returns true if and only if s1 is a subset of s2.
func (s Int) IsSubset(s2 Int) bool {
	for item := range s {
		if !s2.Contains(item) {
			return false
		}
	}
	return true
}

// Equal returns true if and only if s1 is equal (as a set) to s2.
// Two sets are equal if their membership is identical.
// (In practice, this means same elements, order doesn't matter)
func (s Int) Equal(s2 Int) bool {
	return len(s) == len(s2) && s.IsSuperset(s2)
}

type sortableSliceOfInt []int

func (s sortableSliceOfInt) Len() int           { return len(s) }
func (s sortableSliceOfInt) Less(i, j int) bool { return lessInt(s[i], s[j]) }
func (s sortableSliceOfInt) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

// List returns the contents as a sorted int slice.
func (s Int) List() []int {
	res := make(sortableSliceOfInt, 0, len(s))
	for key := range s {
		res = append(res, key)
	}
	sort.Sort(res)
	return res
}

// UnsortedList returns the slice with contents in random order.
func (s Int) UnsortedList() []int {
	res := make([]int, 0, len(s))
	for key := range s {
		res = append(res, key)
	}
	return res
}

// Pop Returns a single element from the set.
func (s Int) Pop() (int, bool) {
	for key := range s {
		delete(s, key)
		return key, true
	}
	return 0, false
}

// Len returns the size of the set.
func (s Int) Len() int {
	return len(s)
}

func lessInt(lhs, rhs int) bool {
	return lhs < rhs
}

// Each traverses the items in the Set, calling the provided function for each
// set member. Traversal will continue until all items in the Set have been
// visited, or if the closure returns false.
func (s Int) Each(f func(item int) bool) {
	for item := range s {
		if !f(item) {
			break
		}
	}
}

// Clone returns a new Set with a copy of s.
func (s Int) Clone() Int {
	ns := NewInt()
	s.Each(func(item int) bool {
		ns[item] = struct{}{}
		return true
	})
	return ns
}
