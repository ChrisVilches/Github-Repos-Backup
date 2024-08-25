package main

import (
	"context"
	"github.com/mholt/archiver/v4"
	"os"
)

func zipDirectory(destName, srcName string) {
	out, _ := os.Create(destName)

	defer out.Close()
	// TODO: I think this needs more error handling, and probably the rest of the program too
	files, _ := archiver.FilesFromDisk(nil, map[string]string{
		// TODO: The nesting is extremely fucked up.
		srcName: "something-goes-here",
	})
	zipArchiver := archiver.Zip{}
	err := zipArchiver.Archive(context.Background(), out, files)
	if err != nil {
		panic(err)
	}
}

func ZipRepo(src, dest string, removeOriginal bool) {
	zipDirectory(dest, src)

	if removeOriginal {
		os.RemoveAll(src)
	}
}
