package helpers

import "time"

func Ptr[T any](input T) *T {
	var value T = input
	return &value
}

func ParseDate(strDate string) (time.Time, error) {
	return time.Parse("2006-01-02", strDate)
}
