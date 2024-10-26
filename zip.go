package main

import (
	"context"
	"github.com/mholt/archiver/v4"
	"os"
)

func zipDirectory(destName, srcName string) error {
	out, _ := os.Create(destName)

	defer out.Close()

	files, err := archiver.FilesFromDisk(nil, map[string]string{
		// TODO: The nesting is extremely fucked up.
		srcName: "something-goes-here",
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
