package rel

import "github.com/marcelocantos/frozen"

// CombineOp specifies which pairings to include in Combine().
type CombineOp int

// The following masks control which elements to include in Combine().
const (
	OnlyOnLHS CombineOp = 1 << iota // Include elements only found on lhs.
	InBoth                          // Include elements found on both sides.
	OnlyOnRHS                       // Include elements only found on rhs.

	AllPairs = OnlyOnLHS | InBoth | OnlyOnRHS
)

// Pair represents a pair of values.
type Pair struct {
	a, b Value
}

// Combine returns a map of names to pairs of corresponding Values from a and b.
// Which names appear in the output is determined by the masks provided in op.
func Combine(a, b Tuple, op CombineOp) map[string]Pair {
	names := make(map[string]Pair, a.Count()+b.Count())
	for e := a.Enumerator(); e.MoveNext(); {
		aName, aValue := e.Current()
		bValue, found := b.Get(aName)
		if !found && (op&OnlyOnLHS != 0) || found && (op&InBoth != 0) {
			names[aName] = Pair{aValue, bValue}
		}
	}
	for e := b.Enumerator(); e.MoveNext(); {
		bName, bValue := e.Current()
		_, found := a.Get(bName)
		if !found && (op&OnlyOnRHS != 0) {
			names[bName] = Pair{nil, bValue}
		}
	}
	return names
}

// CombineNames returns names from a and b according to the given mask.
func CombineNames(a, b Tuple, op CombineOp) Names {
	var sb frozen.SetBuilder
	for name := range Combine(a, b, op) {
		sb.Add(name)
	}
	return Names(sb.Finish())
}

// Merge returns the merger of a and b, if possible or nil otherwise.
// Success requires that common names map to equal values.
func Merge(a, b Tuple) Tuple {
	t := NewTuple()
	for name, pair := range Combine(a, b, AllPairs) {
		if pair.a == nil {
			t = t.With(name, pair.b)
		} else if pair.b == nil || pair.a.Equal(pair.b) {
			t = t.With(name, pair.a)
		} else {
			return nil
		}
	}
	return t
}

// MergeLeftToRight returns the merger of a and b. Key from tuples to the right
// override tuples to the left.
func MergeLeftToRight(t Tuple, ts ...Tuple) Tuple {
	for _, u := range ts {
		for e := u.Enumerator(); e.MoveNext(); {
			name, value := e.Current()
			t = t.With(name, value)
		}
	}
	return t
}
