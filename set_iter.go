//go:build go1.23

package sets

import (
	"iter"
)

func (s Set[T]) Values() iter.Seq[T] {
	return func(yield func(T) bool) {
		for k := range s {
			if !yield(k) {
				return
			}
		}
	}
}
