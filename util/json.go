package util

import (
	"encoding/json"
	"io"
	"os"
)

func WriteJSON[T any](filename string, items []T) error {
	bytes, err := json.MarshalIndent(items, "", "  ")
	if err != nil {
		return err
	}

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(bytes)
	if err != nil {
		return err
	}

	return nil
}

func ReadJSON[T any](filename string) ([]T, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	bytes, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var items []T

	if err := json.Unmarshal(bytes, &items); err != nil {
		return nil, err
	}

	return items, nil
}
