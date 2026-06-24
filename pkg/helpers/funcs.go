package helpers

func Ptr[T any](input T) *T {
	var value T = input
	return &value
}
