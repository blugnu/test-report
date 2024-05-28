package internal

func coalesce[T comparable](values ...T) T {
	z := *new(T)
	for _, v := range values {
		if v != z {
			return v
		}
	}
	return z
}
