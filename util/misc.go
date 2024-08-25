package util

type WithID interface {
	GetID() string
}

func PatchList[T WithID](list, replacements []T) []T {
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
