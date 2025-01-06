package utils

func MapSlice[T any, U any](slice []T, fn func(T) U) []U {
	var result []U
	for _, v := range slice {
		result = append(result, fn(v))
	}
	return result
}
