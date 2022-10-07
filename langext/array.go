package langext

func BoolCount(arr ...bool) int {
	c := 0
	for _, v := range arr {
		if v {
			c++
		}
	}
	return c
}

func IntRange(start int, end int) []int {
	r := make([]int, 0, end-start)
	for i := start; i < end; i++ {
		r = append(r, i)
	}
	return r
}

func ForceArray[T any](v []T) []T {
	if v == nil {
		return make([]T, 0)
	} else {
		return v
	}
}

func ReverseArray[T any](v []T) {
	for i, j := 0, len(v)-1; i < j; i, j = i+1, j-1 {
		v[i], v[j] = v[j], v[i]
	}
}

func InArray[T comparable](needle T, haystack []T) bool {
	for _, v := range haystack {
		if v == needle {
			return true
		}
	}
	return false
}

func ArrUnique[T comparable](array []T) []T {
	m := make(map[T]bool, len(array))
	for _, v := range array {
		m[v] = true
	}
	result := make([]T, 0, len(m))
	for v := range m {
		result = append(result, v)
	}
	return result
}

func ArrEqualsExact[T comparable](arr1 []T, arr2 []T) bool {
	if len(arr1) != len(arr2) {
		return false
	}

	for i := range arr1 {
		if arr1[i] != arr2[i] {
			return false
		}
	}
	return true
}

func ArrFirst[T comparable](arr []T, comp func(v T) bool) (T, bool) {
	for _, v := range arr {
		if comp(v) {
			return v, true
		}
	}
	return *new(T), false
}

func AddToSet[T comparable](set []T, add T) []T {
	for _, v := range set {
		if v == add {
			return set
		}
	}
	return append(set, add)
}
