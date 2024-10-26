package util

type withID interface {
	GetID() string
}

func ListToReadonlyChannel[T any](list []T, size int) <-chan T {
	res := make(chan T, size)

	go func() {
		for _, elem := range list {
			res <- elem
		}

		close(res)
	}()

	return res
}

func PatchList[T withID](list, replacements []T) []T {
	data := make(map[string]T)

	for _, item := range list {
		data[item.GetID()] = item
	}

	for _, item := range replacements {
		data[item.GetID()] = item
	}

	var result []T

	for _, item := range data {
		result = append(result, item)
	}

	return result
}
