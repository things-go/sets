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

// Uint8 is a set of uint8s, implemented via map[uint8]struct{} for minimal memory consumption.
type Uint8 map[uint8]struct{}

// NewUint8 creates a Uint8 from a list of values.
func NewUint8(items ...uint8) Uint8 {
	ss := Uint8{}
	ss.Insert(items...)
	return ss
}

// NewUint8From creates a Uint8 from a keys of a map[uint8](? extends interface{}).
// If the value passed in is not actually a map, this will panic.
func NewUint8From(theMap interface{}) Uint8 {
	v := reflect.ValueOf(theMap)
	ret := Uint8{}

	for _, keyValue := range v.MapKeys() {
		ret[keyValue.Interface().(uint8)] = struct{}{}
	}
	return ret
}

// Insert adds items to the set.
func (s Uint8) Insert(items ...uint8) Uint8 {
	for _, item := range items {
		s[item] = struct{}{}
	}
	return s
}

// Delete removes all items from the set.
func (s Uint8) Delete(items ...uint8) Uint8 {
	for _, item := range items {
		delete(s, item)
	}
	return s
}

// Contains returns true if and only if item is contained in the set.
func (s Uint8) Contains(item uint8) bool {
	_, contained := s[item]
	return contained
}

// ContainsAll returns true if and only if all items are contained in the set.
func (s Uint8) ContainsAll(items ...uint8) bool {
	for _, item := range items {
		if !s.Contains(item) {
			return false
		}
	}
	return true
}

// ContainsAny returns true if any items are contained in the set.
func (s Uint8) ContainsAny(items ...uint8) bool {
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
func (s Uint8) Difference(s2 Uint8) Uint8 {
	result := NewUint8()
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
func (s Uint8) Union(s2 Uint8) Uint8 {
	result := NewUint8()
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
func (s Uint8) Intersection(s2 Uint8) Uint8 {
	var walk, other Uint8
	result := NewUint8()
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
func (s Uint8) Merge(s2 Uint8) Uint8 {
	for key := range s2 {
		s[key] = struct{}{}
	}
	return s
}

// IsSuperset returns true if and only if s1 is a superset of s2.
func (s Uint8) IsSuperset(s2 Uint8) bool {
	for item := range s2 {
		if !s.Contains(item) {
			return false
		}
	}
	return true
}

// IsSubset returns true if and only if s1 is a subset of s2.
func (s Uint8) IsSubset(s2 Uint8) bool {
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
func (s Uint8) Equal(s2 Uint8) bool {
	return len(s) == len(s2) && s.IsSuperset(s2)
}

type sortableSliceOfUint8 []uint8

func (s sortableSliceOfUint8) Len() int           { return len(s) }
func (s sortableSliceOfUint8) Less(i, j int) bool { return lessUint8(s[i], s[j]) }
func (s sortableSliceOfUint8) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

// List returns the contents as a sorted uint8 slice.
func (s Uint8) List() []uint8 {
	res := make(sortableSliceOfUint8, 0, len(s))
	for key := range s {
		res = append(res, key)
	}
	sort.Sort(res)
	return res
}

// UnsortedList returns the slice with contents in random order.
func (s Uint8) UnsortedList() []uint8 {
	res := make([]uint8, 0, len(s))
	for key := range s {
		res = append(res, key)
	}
	return res
}

// Pop Returns a single element from the set.
func (s Uint8) Pop() (uint8, bool) {
	for key := range s {
		delete(s, key)
		return key, true
	}
	return 0, false
}

// Len returns the size of the set.
func (s Uint8) Len() int {
	return len(s)
}

func lessUint8(lhs, rhs uint8) bool {
	return lhs < rhs
}

// Each traverses the items in the Set, calling the provided function for each
// set member. Traversal will continue until all items in the Set have been
// visited, or if the closure returns false.
func (s Uint8) Each(f func(item uint8) bool) {
	for item := range s {
		if !f(item) {
			break
		}
	}
}

// Clone returns a new Set with a copy of s.
func (s Uint8) Clone() Uint8 {
	ns := NewUint8()
	s.Each(func(item uint8) bool {
		ns[item] = struct{}{}
		return true
	})
	return ns
}
