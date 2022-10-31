package util

func Pointer[E any](e E) *E {
	return &e
}
