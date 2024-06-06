package sets

// Set sets.Set is a set of interface,
// implemented via map[T]struct{} for minimal memory consumption.
type Set[T comparable] map[T]struct{}

func newSet[T comparable](cap int) Set[T] {
	return make(Set[T], cap)
}

// New creates a T from a list of values.
func New[T comparable](items ...T) Set[T] {
	ret := newSet[T](len(items))
	return ret.Insert(items...)
}

// NewFrom creates a T from a keys of a map[T](? extends any).
// If the value passed in is not actually a map, this will panic.
func NewFrom[T comparable, V any, M ~map[T]V](m M) Set[T] {
	ret := newSet[T](len(m))
	for k := range m {
		ret[k] = struct{}{}
	}
	return ret
}

// Insert adds items to the set.
func (s Set[T]) Insert(items ...T) Set[T] {
	for _, item := range items {
		s[item] = struct{}{}
	}
	return s
}

// Delete removes all items from the set.
func (s Set[T]) Delete(items ...T) Set[T] {
	for _, item := range items {
		delete(s, item)
	}
	return s
}

// Contains returns true if and only if item is contained in the set.
func (s Set[T]) Contains(item T) bool {
	_, contained := s[item]
	return contained
}

// ContainsAll returns true if and only if all items are contained in the set.
func (s Set[T]) ContainsAll(items ...T) bool {
	for _, item := range items {
		if !s.Contains(item) {
			return false
		}
	}
	return true
}

// ContainsAny returns true if any items are contained in the set.
func (s Set[T]) ContainsAny(items ...T) bool {
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
func (s Set[T]) Difference(s2 Set[T]) Set[T] {
	result := newSet[T](len(s))
	for key := range s {
		if !s2.Contains(key) {
			result[key] = struct{}{}
		}
	}
	return result
}

// DifferenceSlice returns a slices of objects that are not in s2
// For example:
// s1 = {a1, a2, a3}
// s2 = {a1, a2, a4, a5}
// s1.Difference(s2) = {a3}
// s2.Difference(s1) = {a4, a5}.
func (s Set[T]) DifferenceSlice(s2 Set[T]) []T {
	result := make([]T, 0, len(s))
	for key := range s {
		if !s2.Contains(key) {
			result = append(result, key)
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
func (s Set[T]) Union(s2 Set[T]) Set[T] {
	result := newSet[T](len(s) + len(s2))
	for key := range s {
		result[key] = struct{}{}
	}
	for key := range s2 {
		result[key] = struct{}{}
	}
	return result
}

// UnionSlice returns a slice which includes items in either s1 or s2.
// For example:
// s1 = {a1, a2}
// s2 = {a3, a4}
// s1.Union(s2) = {a1, a2, a3, a4}
// s2.Union(s1) = {a1, a2, a3, a4}.
func (s Set[T]) UnionSlice(s2 Set[T]) []T {
	result := make([]T, 0, len(s)+len(s2))
	for key := range s {
		result = append(result, key)
	}
	for key := range s2 {
		if !s.Contains(key) {
			result = append(result, key)
		}
	}
	return result
}

// Intersection returns a new set which includes the item in BOTH s1 and s2
// For example:
// s1 = {a1, a2}
// s2 = {a2, a3}
// s1.Intersection(s2) = {a2}.
func (s Set[T]) Intersection(s2 Set[T]) Set[T] {
	var walk, other Set[T]

	if s.Len() < s2.Len() {
		walk = s
		other = s2
	} else {
		walk = s2
		other = s
	}
	ret := newSet[T](min(len(s), len(s2)))
	for key := range walk {
		if other.Contains(key) {
			ret[key] = struct{}{}
		}
	}
	return ret
}

// IntersectionSlice returns a slice which includes the item in BOTH s1 and s2
// For example:
// s1 = {a1, a2}
// s2 = {a2, a3}
// s1.Intersection(s2) = {a2}.
func (s Set[T]) IntersectionSlice(s2 Set[T]) []T {
	var walk, other Set[T]

	if s.Len() < s2.Len() {
		walk = s
		other = s2
	} else {
		walk = s2
		other = s
	}
	ret := make([]T, 0, min(len(s), len(s2)))
	for key := range walk {
		if other.Contains(key) {
			ret = append(ret, key)
		}
	}
	return ret
}

// Diff returns s diff of s2, return added, removed, remained sets
// with the given s2 set.
// For example:
// s1 = {a1, a3, a5, a7}
// s2 = {a3, a4, a5, a6}
// added = {a4, a6}
// removed = {a1, a7}
// remained = {a3, a6}
func (s Set[T]) Diff(s2 Set[T]) (added, removed, remained Set[T]) {
	removed = newSet[T](len(s))
	added = newSet[T](len(s2))
	remained = newSet[T](len(s))
	for key := range s {
		if s2.Contains(key) {
			remained[key] = struct{}{}
		} else {
			removed[key] = struct{}{}
		}
	}
	for key := range s2 {
		if !s.Contains(key) {
			added[key] = struct{}{}
		}
	}
	return added, removed, remained
}

// Diff returns s diff of s2, return added, removed, remained slices
// with the given s2 set.
// For example:
// s1 = {a1, a3, a5, a7}
// s2 = {a3, a4, a5, a6}
// added = {a4, a6}
// removed = {a1, a7}
// remained = {a3, a6}
func (s Set[T]) DiffSlice(s2 Set[T]) (added, removed, remained []T) {
	removed = make([]T, 0, len(s))
	added = make([]T, 0, len(s2))
	remained = make([]T, 0, len(s))
	for key := range s {
		if s2.Contains(key) {
			remained = append(remained, key)
		} else {
			removed = append(removed, key)
		}
	}
	for key := range s2 {
		if !s.Contains(key) {
			added = append(added, key)
		}
	}
	return added, removed, remained
}

// Merge is like Union, however it modifies the current set it's applied on
// with the given s2 set.
// For example:
// s1 = {a1, a2}
// s2 = {a3, a4}
// s1.Merge(s2), s1 = {a1, a2, a3, a4}
// s2.Merge(s1), s2 = {a1, a2, a3, a4}.
func (s Set[T]) Merge(s2 Set[T]) Set[T] {
	for item := range s2 {
		s[item] = struct{}{}
	}
	return s
}

// IsSuperset returns true if and only if s1 is a superset of s2.
func (s Set[T]) IsSuperset(s2 Set[T]) bool {
	for item := range s2 {
		if !s.Contains(item) {
			return false
		}
	}
	return true
}

// IsSubset returns true if and only if s1 is a superset of s2.
func (s Set[T]) IsSubset(s2 Set[T]) bool {
	for item := range s {
		if !s2.Contains(item) {
			return false
		}
	}
	return true
}

// List returns the contents as a sorted slice.
func (s Set[T]) List() []T {
	res := make([]T, 0, len(s))
	for key := range s {
		res = append(res, key)
	}
	return res
}

// Equal returns true if and only if s1 is equal (as a set) to s2.
// Two sets are equal if their membership is identical.
// (In practice, this means same elements, order doesn't matter).
func (s Set[T]) Equal(s2 Set[T]) bool {
	return len(s) == len(s2) && s.IsSuperset(s2)
}

// Pop Returns a single element from the set.
func (s Set[T]) Pop() (v T, ok bool) {
	for key := range s {
		delete(s, key)
		return key, true
	}
	return
}

// Len returns the size of the set.
func (s Set[T]) Len() int {
	return len(s)
}

// Each traverses the items in the Set, calling the provided function for each
// set member. Traversal will continue until all items in the Set have been
// visited, or if the closure returns false.
func (s Set[T]) Each(f func(item T) bool) {
	for item := range s {
		if !f(item) {
			break
		}
	}
}

// Clone returns a new Set with a copy of s.
func (s Set[T]) Clone() Set[T] {
	return NewFrom(s)
}
