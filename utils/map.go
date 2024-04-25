package utils

type Pair[K, V any] struct {
	Key   K
	Value V
}

func MapToPairs[K comparable, V any](m map[K]V) []Pair[K, V] {
	pairs := make([]Pair[K, V], 0, len(m))
	for k, v := range m {
		pairs = append(pairs, Pair[K, V]{k, v})
	}
	return pairs
}

func PairValues[K, V any](pairs []Pair[K, V]) []V {
	values := make([]V, 0, len(pairs))
	for _, pair := range pairs {
		values = append(values, pair.Value)
	}
	return values
}

func Merge[K comparable, V any](a, b map[K]V) map[K]V {
	m := make(map[K]V, len(a)+len(b))
	for k, v := range a {
		m[k] = v
	}
	for k, v := range b {
		m[k] = v
	}
	return m
}
