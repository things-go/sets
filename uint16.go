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

// Uint16 is a set of uint16s, implemented via map[uint16]struct{} for minimal memory consumption.
type Uint16 map[uint16]struct{}

// NewUint16 creates a Uint16 from a list of values.
func NewUint16(items ...uint16) Uint16 {
	ss := Uint16{}
	ss.Insert(items...)
	return ss
}

// NewUint16From creates a Uint16 from a keys of a map[uint16](? extends interface{}).
// If the value passed in is not actually a map, this will panic.
func NewUint16From(theMap interface{}) Uint16 {
	v := reflect.ValueOf(theMap)
	ret := Uint16{}

	for _, keyValue := range v.MapKeys() {
		ret[keyValue.Interface().(uint16)] = struct{}{}
	}
	return ret
}

// Insert adds items to the set.
func (s Uint16) Insert(items ...uint16) Uint16 {
	for _, item := range items {
		s[item] = struct{}{}
	}
	return s
}

// Delete removes all items from the set.
func (s Uint16) Delete(items ...uint16) Uint16 {
	for _, item := range items {
		delete(s, item)
	}
	return s
}

// Contains returns true if and only if item is contained in the set.
func (s Uint16) Contains(item uint16) bool {
	_, contained := s[item]
	return contained
}

// ContainsAll returns true if and only if all items are contained in the set.
func (s Uint16) ContainsAll(items ...uint16) bool {
	for _, item := range items {
		if !s.Contains(item) {
			return false
		}
	}
	return true
}

// ContainsAny returns true if any items are contained in the set.
func (s Uint16) ContainsAny(items ...uint16) bool {
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
func (s Uint16) Difference(s2 Uint16) Uint16 {
	result := NewUint16()
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
func (s Uint16) Union(s2 Uint16) Uint16 {
	result := NewUint16()
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
func (s Uint16) Intersection(s2 Uint16) Uint16 {
	var walk, other Uint16
	result := NewUint16()
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
func (s Uint16) Merge(s2 Uint16) Uint16 {
	for key := range s2 {
		s[key] = struct{}{}
	}
	return s
}

// IsSuperset returns true if and only if s1 is a superset of s2.
func (s Uint16) IsSuperset(s2 Uint16) bool {
	for item := range s2 {
		if !s.Contains(item) {
			return false
		}
	}
	return true
}

// IsSubset returns true if and only if s1 is a subset of s2.
func (s Uint16) IsSubset(s2 Uint16) bool {
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
func (s Uint16) Equal(s2 Uint16) bool {
	return len(s) == len(s2) && s.IsSuperset(s2)
}

type sortableSliceOfUint16 []uint16

func (s sortableSliceOfUint16) Len() int           { return len(s) }
func (s sortableSliceOfUint16) Less(i, j int) bool { return lessUint16(s[i], s[j]) }
func (s sortableSliceOfUint16) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

// List returns the contents as a sorted uint16 slice.
func (s Uint16) List() []uint16 {
	res := make(sortableSliceOfUint16, 0, len(s))
	for key := range s {
		res = append(res, key)
	}
	sort.Sort(res)
	return res
}

// UnsortedList returns the slice with contents in random order.
func (s Uint16) UnsortedList() []uint16 {
	res := make([]uint16, 0, len(s))
	for key := range s {
		res = append(res, key)
	}
	return res
}

// Pop Returns a single element from the set.
func (s Uint16) Pop() (uint16, bool) {
	for key := range s {
		delete(s, key)
		return key, true
	}
	return 0, false
}

// Len returns the size of the set.
func (s Uint16) Len() int {
	return len(s)
}

func lessUint16(lhs, rhs uint16) bool {
	return lhs < rhs
}

// Each traverses the items in the Set, calling the provided function for each
// set member. Traversal will continue until all items in the Set have been
// visited, or if the closure returns false.
func (s Uint16) Each(f func(item uint16) bool) {
	for item := range s {
		if !f(item) {
			break
		}
	}
}

// Clone returns a new Set with a copy of s.
func (s Uint16) Clone() Uint16 {
	ns := NewUint16()
	s.Each(func(item uint16) bool {
		ns[item] = struct{}{}
		return true
	})
	return ns
}
