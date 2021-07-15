package sets

import (
	"reflect"
)

// Set sets.Set is a set of interface,
// implemented via map[interface{}]struct{} for minimal memory consumption.
type Set struct {
	m   map[interface{}]struct{}
	cmp Comparator
}

// Option option for New.
type Option func(*Set)

// WithItems with git items.
func WithItems(items ...interface{}) Option {
	return func(s *Set) {
		s.Insert(items...)
	}
}

// WithComparator with user's Comparator only for sort
func WithComparator(cmp Comparator) Option {
	return func(s *Set) {
		s.cmp = cmp
	}
}

// New creates a interface{} from a list of values.
func New(opts ...Option) *Set {
	ss := &Set{
		m: make(map[interface{}]struct{}),
	}
	for _, opt := range opts {
		opt(ss)
	}
	return ss
}

// NewSetFrom creates a interface{} from a keys of a map[interface{}](? extends interface{}).
// If the value passed in is not actually a map, this will panic.
func NewSetFrom(theMap interface{}, opts ...Option) *Set {
	v := reflect.ValueOf(theMap)
	ret := New(opts...)

	for _, keyValue := range v.MapKeys() {
		ret.m[keyValue.Interface()] = struct{}{}
	}
	return ret
}

// Insert adds items to the set.
func (s *Set) Insert(items ...interface{}) *Set {
	for _, item := range items {
		s.m[item] = struct{}{}
	}
	return s
}

// Delete removes all items from the set.
func (s *Set) Delete(items ...interface{}) *Set {
	for _, item := range items {
		delete(s.m, item)
	}
	return s
}

// Contains returns true if and only if item is contained in the set.
func (s *Set) Contains(item interface{}) bool {
	_, contained := s.m[item]
	return contained
}

// ContainsAll returns true if and only if all items are contained in the set.
func (s *Set) ContainsAll(items ...interface{}) bool {
	for _, item := range items {
		if !s.Contains(item) {
			return false
		}
	}
	return true
}

// ContainsAny returns true if any items are contained in the set.
func (s *Set) ContainsAny(items ...interface{}) bool {
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
// s2.Difference(s1) = {a4, a5}.
func (s *Set) Difference(s2 *Set) *Set {
	result := New(WithComparator(s.cmp))
	for key := range s.m {
		if !s2.Contains(key) {
			result.m[key] = struct{}{}
		}
	}
	return result
}

// Union returns a new set which includes items in either s1 or s2.
// For example:
// s1 = {a1, a2}
// s2 = {a3, a4}
// s1.Union(s2) = {a1, a2, a3, a4}
// s2.Union(s1) = {a1, a2, a3, a4}.
func (s *Set) Union(s2 *Set) *Set {
	result := New(WithComparator(s.cmp))
	for key := range s.m {
		result.m[key] = struct{}{}
	}
	for key := range s2.m {
		result.m[key] = struct{}{}
	}
	return result
}

// Intersection returns a new set which includes the item in BOTH s1 and s2
// For example:
// s1 = {a1, a2}
// s2 = {a2, a3}
// s1.Intersection(s2) = {a2}.
func (s *Set) Intersection(s2 *Set) *Set {
	var walk, other *Set
	result := New(WithComparator(s.cmp))
	if s.Len() < s2.Len() {
		walk = s
		other = s2
	} else {
		walk = s2
		other = s
	}
	for key := range walk.m {
		if other.Contains(key) {
			result.m[key] = struct{}{}
		}
	}
	return result
}

// Merge is like Union, however it modifies the current set it's applied on
// with the given s2 set.
// For example:
// s1 = {a1, a2}
// s2 = {a3, a4}
// s1.Merge(s2), s1 = {a1, a2, a3, a4}
// s2.Merge(s1), s2 = {a1, a2, a3, a4}.
func (s *Set) Merge(s2 *Set) *Set {
	for item := range s2.m {
		s.m[item] = struct{}{}
	}
	return s
}

// IsSuperset returns true if and only if s1 is a superset of s2.
func (s *Set) IsSuperset(s2 *Set) bool {
	for item := range s2.m {
		if !s.Contains(item) {
			return false
		}
	}
	return true
}

// IsSubset returns true if and only if s1 is a superset of s2.
func (s *Set) IsSubset(s2 *Set) bool {
	for item := range s.m {
		if !s2.Contains(item) {
			return false
		}
	}
	return true
}

// List returns the contents as a sorted slice.
func (s *Set) List() []interface{} {
	res := s.UnsortedList()
	NewContainer(res, s.cmp).Sort()
	return res
}

// UnsortedList returns the slice with contents in random order.
func (s *Set) UnsortedList() []interface{} {
	res := make([]interface{}, 0, len(s.m))
	for key := range s.m {
		res = append(res, key)
	}
	return res
}

// Equal returns true if and only if s1 is equal (as a set) to s2.
// Two sets are equal if their membership is identical.
// (In practice, this means same elements, order doesn't matter).
func (s *Set) Equal(s2 *Set) bool {
	return len(s.m) == len(s2.m) && s.IsSuperset(s2)
}

// Pop Returns a single element from the set.
func (s *Set) Pop() (interface{}, bool) {
	for key := range s.m {
		delete(s.m, key)
		return key, true
	}
	return nil, false
}

// Len returns the size of the set.
func (s *Set) Len() int {
	return len(s.m)
}

// Each traverses the items in the Set, calling the provided function for each
// set member. Traversal will continue until all items in the Set have been
// visited, or if the closure returns false.
func (s *Set) Each(f func(item interface{}) bool) {
	for item := range s.m {
		if !f(item) {
			break
		}
	}
}

// Clone returns a new Set with a copy of s.
func (s *Set) Clone() *Set {
	ns := New(WithComparator(s.cmp))
	s.Each(func(item interface{}) bool {
		ns.m[item] = struct{}{}
		return true
	})
	return ns
}
