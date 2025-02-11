package utils

func MapSlice[T any, U any](slice []T, fn func(T) U) []U {
	var result []U
	for _, v := range slice {
		result = append(result, fn(v))
	}
	return result
}

func FilterSlice[T any](slice []T, fn func(T) bool) []T {
    var result []T
    for _, v := range slice {
        if fn(v) {
            result = append(result, v)
        }
    }

    return result
}

func SliceIsAll[T any](slice []T, fn func(T) bool) bool {
    for _, v := range slice {
        if fn(v) {
            continue
        } else {
            return false
        }
    }

    return true
}

func SliceCount[T any](slice []T, fn func(T) bool) int {
    result := 0
    for _, v := range slice {
        if fn(v) {
            result += 1
        }
    }

    return result
}
