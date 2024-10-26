package main

import (
	"context"
	"os"
	"path/filepath"

	"github.com/mholt/archiver/v4"
)

func zipDirectory(destName, srcName string) error {
	out, err := os.Create(destName)

	if err != nil {
		return err
	}

	defer out.Close()

	baseFolder := filepath.Base(srcName)

	files, err := archiver.FilesFromDisk(nil, map[string]string{
		srcName: baseFolder,
	})

	if err != nil {
		return err
	}

	zipArchiver := archiver.Zip{}

	return zipArchiver.Archive(context.Background(), out, files)
}

func zipRepo(src, dest string, removeOriginal bool) error {
	err := zipDirectory(dest, src)

	if err != nil {
		return err
	}

	if removeOriginal {
		return os.RemoveAll(src)
	}

	return nil
}
