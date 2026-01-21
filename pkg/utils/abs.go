package utils

func Abs[T ~int | ~int64 | ~int32 | ~float32 | ~float64](v T) T {
	if v < 0 {
		return -v
	}
	return v
}
